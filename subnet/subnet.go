package subnet

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/deluxesande/network-scanner/utils"
)

const (
	defaultConcurrency = 100
	defaultTimeout     = 1 * time.Second
	maxHosts           = 65536 // Refuse anything wider than a /16
	dnsTimeout         = 1 * time.Second
)

var (
	ttlPattern = regexp.MustCompile(`(?i)ttl[=:]\s*(\d+)`)
	arpPattern = regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)\D+(([0-9a-fA-F]{1,2}[:-]){5}[0-9a-fA-F]{1,2})`)
)

// Clamp caller values into ranges we know are safe to run
func withDefaults(opts utils.ScanOptions) utils.ScanOptions {
	if opts.Concurrency <= 0 || opts.Concurrency > 1000 {
		opts.Concurrency = defaultConcurrency
	}
	if opts.Timeout <= 0 || opts.Timeout > 10*time.Second {
		opts.Timeout = defaultTimeout
	}
	return opts
}

// Sweep every host in each subnet, then merge with the kernel's ARP cache
func ScanSubnets(ctx context.Context, subnets []string, opts utils.ScanOptions) ([]utils.Device, error) {
	opts = withDefaults(opts)

	if len(subnets) == 0 {
		return nil, fmt.Errorf("no subnet selected")
	}

	if _, err := exec.LookPath("ping"); err != nil {
		return nil, fmt.Errorf("ping not found on PATH; install iputils-ping (Debian) or iputils (Fedora)")
	}

	var prefixes []netip.Prefix
	var targets []netip.Addr
	for _, raw := range subnets {
		prefix, err := parseTarget(raw)
		if err != nil {
			return nil, err
		}

		// A typo'd public range is a real incident, not a scan
		if !opts.AllowPublic && !isPrivate(prefix) {
			return nil, fmt.Errorf("%s is not a private range; pass --force to scan it anyway", prefix)
		}

		hosts, err := enumerateHosts(prefix)
		if err != nil {
			return nil, err
		}
		prefixes = append(prefixes, prefix)
		targets = append(targets, hosts...)
	}

	devices := sweep(ctx, targets, opts)

	// After the sweep, the cache holds every host that answered on the wire
	return mergeDevices(devices, GetMacTable(), prefixes), nil
}

func sweep(ctx context.Context, targets []netip.Addr, opts utils.ScanOptions) []utils.Device {
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		devices []utils.Device
		sem     = make(chan struct{}, opts.Concurrency)
		done    int
	)

	// One shared semaphore: total concurrency stays flat across subnets
	for _, target := range targets {
		select {
		case <-ctx.Done():
			wg.Wait()
			return devices
		case sem <- struct{}{}:
		}

		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			defer func() { <-sem }()

			hostCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
			alive, osName := pingHost(hostCtx, ip)
			cancel()

			var device utils.Device
			if alive {
				device = utils.Device{IP: ip, OS: osName}
				if opts.ResolveDNS {
					device.Hostname = lookupHostname(ctx, ip)
				}
			}

			mu.Lock()
			if alive {
				devices = append(devices, device)
			}
			done++
			if opts.OnProgress != nil {
				opts.OnProgress(done, len(targets))
			}
			mu.Unlock()
		}(target.String())
	}

	wg.Wait()
	return devices
}

// Union of ping replies and ARP entries; ICMP-filtered hosts still show up
func mergeDevices(pinged []utils.Device, macs map[string]string, prefixes []netip.Prefix) []utils.Device {
	byIP := make(map[string]utils.Device, len(pinged)+len(macs))
	for _, device := range pinged {
		byIP[device.IP] = device
	}

	for ip, mac := range macs {
		addr, err := netip.ParseAddr(ip)
		if err != nil || !inPrefixes(addr, prefixes) {
			continue // ARP also holds the gateway and neighbours outside the target
		}
		if device, ok := byIP[ip]; ok {
			device.MAC = mac
			byIP[ip] = device
			continue
		}
		byIP[ip] = utils.Device{IP: ip, MAC: mac, OS: "Unknown"}
	}

	devices := make([]utils.Device, 0, len(byIP))
	for _, device := range byIP {
		devices = append(devices, device)
	}
	sort.Slice(devices, func(i, j int) bool { return lessIP(devices[i].IP, devices[j].IP) })
	return devices
}

func inPrefixes(addr netip.Addr, prefixes []netip.Prefix) bool {
	for _, prefix := range prefixes {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}

// Numeric order, so identical scans produce identical files
func lessIP(a, b string) bool {
	addrA, errA := netip.ParseAddr(a)
	addrB, errB := netip.ParseAddr(b)
	if errA != nil || errB != nil {
		return a < b
	}
	return addrA.Less(addrB)
}

// Accept 192.168.1.0/24 or a bare 192.168.1.10
func parseTarget(raw string) (netip.Prefix, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return netip.Prefix{}, fmt.Errorf("empty subnet")
	}

	if strings.Contains(raw, "/") {
		prefix, err := netip.ParsePrefix(raw)
		if err != nil {
			return netip.Prefix{}, fmt.Errorf("invalid subnet %q: %w", raw, err)
		}
		if !prefix.Addr().Is4() {
			return netip.Prefix{}, fmt.Errorf("subnet %q is IPv6; only IPv4 is supported", raw)
		}
		return prefix.Masked(), nil
	}

	addr, err := netip.ParseAddr(raw)
	if err != nil || !addr.Is4() {
		return netip.Prefix{}, fmt.Errorf("invalid subnet %q: expected 192.168.1.0/24 or 192.168.1.10", raw)
	}
	return netip.PrefixFrom(addr, 32), nil
}

// Every usable host in the prefix, network and broadcast excluded
func enumerateHosts(prefix netip.Prefix) ([]netip.Addr, error) {
	prefix = prefix.Masked()
	size := int64(1) << (32 - prefix.Bits())

	firstUsable, lastUsable := int64(1), size-2
	if prefix.Bits() >= 31 { // /31 and /32 have no network or broadcast address
		firstUsable, lastUsable = 0, size-1
	}
	if total := lastUsable - firstUsable + 1; total > maxHosts {
		return nil, fmt.Errorf("subnet %s has %d hosts, limit is %d", prefix, total, maxHosts)
	}

	hosts := make([]netip.Addr, 0, lastUsable-firstUsable+1)
	addr := prefix.Addr()
	for i := int64(0); i <= lastUsable; i++ {
		if i >= firstUsable {
			hosts = append(hosts, addr)
		}
		addr = addr.Next()
	}
	return hosts, nil
}

// RFC1918, loopback and link-local; anything else needs --force
func isPrivate(prefix netip.Prefix) bool {
	addr := prefix.Addr()
	return addr.IsPrivate() || addr.IsLoopback() || addr.IsLinkLocalUnicast()
}

// Alive means the ping binary exited 0; the OS is best-effort from the TTL line
func pingHost(ctx context.Context, ip string) (bool, string) {
	countFlag := "-c"
	if runtime.GOOS == "windows" {
		countFlag = "-n"
	}

	// The context deadline is the timeout, so there is no per-OS -W/-w difference
	output, err := exec.CommandContext(ctx, "ping", countFlag, "1", ip).Output()
	if err != nil {
		return false, ""
	}
	return true, EstimateOS(parseTTL(string(output)))
}

// Backward-compatible wrapper; ScanSubnets calls pingHost directly
func PingWithTTL(ip string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	return pingHost(ctx, ip)
}

func parseTTL(output string) string {
	match := ttlPattern.FindStringSubmatch(output)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

// Default TTLs as shipped, ascending; each router decrements by one per hop
var initialTTLs = []struct {
	ttl  int
	name string
}{
	{32, "Windows 95/98/ME, NT 3.51/4.0"},
	{60, "macOS (legacy)/FreeBSD"},
	{62, "Windows CE/Mobile"},
	{64, "Linux/macOS/iOS/Android/BSD"},
	{100, "Google Chrome OS"},
	{128, "Windows NT and later"},
	{200, "Solaris (legacy)"},
	{252, "IBM AIX"},
	{255, "Cisco IOS/Solaris/HP-UX/AIX"},
}

// Exact TTLs are the OS default; lower values are that default minus one per hop
func EstimateOS(ttl string) string {
	value, err := strconv.Atoi(ttl)
	if err != nil || value <= 0 || value > 255 {
		return "Unknown"
	}

	for _, entry := range initialTTLs {
		if value == entry.ttl {
			return entry.name
		}
		if value < entry.ttl {
			hops := entry.ttl - value
			if hops == 1 {
				return fmt.Sprintf("%s (~1 hop)", entry.name)
			}
			return fmt.Sprintf("%s (~%d hops)", entry.name, hops)
		}
	}

	return "Unknown"
}

// Best-effort PTR with its own deadline; a dead resolver must not stall the sweep
func lookupHostname(ctx context.Context, ip string) string {
	lookupCtx, cancel := context.WithTimeout(ctx, dnsTimeout)
	defer cancel()

	names, err := net.DefaultResolver.LookupAddr(lookupCtx, ip)
	if err != nil || len(names) == 0 {
		return ""
	}
	return strings.TrimSuffix(names[0], ".")
}

func GetLocalSubnets() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("list interfaces: %w", err)
	}

	var subnets []string
	seen := make(map[string]bool)
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.To4() == nil || ipnet.IP.IsLoopback() {
				continue
			}
			prefix, err := netip.ParsePrefix(ipnet.String())
			if err != nil || prefix.Addr().IsLinkLocalUnicast() {
				continue
			}
			cidr := prefix.Masked().String()
			if !seen[cidr] {
				seen[cidr] = true
				subnets = append(subnets, cidr)
			}
		}
	}

	sort.Strings(subnets)
	return subnets, nil
}

func AskSubnetChoice(subnets []string) ([]string, error) {
	fmt.Println("Choose a subnet to scan:")
	for i, raw := range subnets {
		fmt.Printf("  [%d] %s\n", i+1, raw)
	}
	fmt.Print("  [0] All subnets\nEnter choice (e.g. 0 or 1,2): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("read choice: %w", err)
	}

	input = strings.TrimSpace(input)
	if input == "0" {
		return subnets, nil
	}

	var selected []string
	for _, part := range strings.Split(input, ",") {
		index, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || index < 1 || index > len(subnets) {
			return nil, fmt.Errorf("invalid choice %q: pick 0 or a number from 1 to %d", part, len(subnets))
		}
		selected = append(selected, subnets[index-1])
	}
	return selected, nil
}

// MAC addresses the kernel has resolved, keyed by IP
func GetMacTable() map[string]string {
	if runtime.GOOS == "linux" {
		if table, err := readProcArp(); err == nil {
			return table
		}
	}

	output, err := exec.Command("arp", "-a").Output()
	if err != nil {
		return map[string]string{} // No ARP tool; hosts still report, just without MACs
	}
	return parseArpLines(string(output))
}

// /proc/net/arp needs no external tool; incomplete entries carry 00:00:00:00:00:00
func readProcArp() (map[string]string, error) {
	data, err := os.ReadFile("/proc/net/arp")
	if err != nil {
		return nil, err
	}

	table := make(map[string]string)
	for i, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if i == 0 || len(fields) < 4 || fields[3] == "00:00:00:00:00:00" {
			continue // Header, malformed, or unresolved
		}
		table[fields[0]] = normalizeMac(fields[3])
	}
	return table, nil
}

func parseArpLines(output string) map[string]string {
	table := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		match := arpPattern.FindStringSubmatch(line)
		if len(match) < 3 {
			continue
		}
		table[match[1]] = normalizeMac(match[2])
	}
	return table
}

// Uppercase and zero-padded: 0:1c:42:0:0:8 -> 00:1C:42:00:00:08
func normalizeMac(raw string) string {
	parts := strings.FieldsFunc(strings.ToUpper(raw), func(r rune) bool { return r == ':' || r == '-' })
	for i, part := range parts {
		if len(part) == 1 {
			parts[i] = "0" + part
		}
	}
	return strings.Join(parts, ":")
}

func ExportToJSON(devices []utils.Device, filename string) error {
	data, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		return fmt.Errorf("encode results: %w", err)
	}
	data = append(data, '\n')

	// Write beside the target then rename, so an interrupt never leaves half a report
	temp, err := os.CreateTemp(filepath.Dir(filename), ".netscanner-*.json")
	if err != nil {
		return fmt.Errorf("create %s: %w", filename, err)
	}
	defer os.Remove(temp.Name())

	// Reports are shared; CreateTemp's 0600 would hide them from the team
	if err := temp.Chmod(0o644); err != nil {
		temp.Close()
		return fmt.Errorf("chmod %s: %w", filename, err)
	}

	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return fmt.Errorf("write %s: %w", filename, err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close %s: %w", filename, err)
	}
	return os.Rename(temp.Name(), filename)
}
