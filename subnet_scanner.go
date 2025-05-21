package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/fatih/color"
)

const (
	startIP     = 1
	endIP       = 254
	concurrency = 100
)

type Device struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	MAC      string `json:"mac"`
}

func getLocalSubnets() []string {
	var subnets []string
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting interfaces:", err)
		return subnets
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ip := ipnet.IP.Mask(ipnet.Mask)
				subnet := ip.String()
				subnets = append(subnets, subnet)
			}
		}
	}
	return subnets
}

func pingWithTTL(ip string) (bool, string) {
	cmd := exec.Command("ping", "-n", "1", "-w", "1000", ip) // Windows ping
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, ""
	}
	outStr := string(output)
	if !strings.Contains(outStr, "Reply from") {
		return false, ""
	}

	// Extract TTL
	ttl := ""
	lines := strings.Split(outStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "TTL=") {
			parts := strings.Split(line, "TTL=")
			if len(parts) > 1 {
				ttl = strings.TrimSpace(parts[1])
			}
			break
		}
	}

	return true, estimateOS(ttl)
}

func estimateOS(ttl string) string {
	switch ttl {
	case "32":
		return "Windows 95/98/ME"
	case "64":
		return "Linux/macOS"
	case "128":
		return "Windows"
	case "255":
		return "Cisco/Unix"
	case "60":
		return "FreeBSD"
	case "200":
		return "Solaris"
	default:
		return "Unknown"
	}
}

func scanSubnet(subnet string, wg *sync.WaitGroup, results chan<- Device) {
	defer wg.Done()
	var innerWg sync.WaitGroup
	sem := make(chan struct{}, concurrency)

	base := strings.Join(strings.Split(subnet, ".")[:3], ".") + "."

	for i := startIP; i <= endIP; i++ {
		ip := fmt.Sprintf("%s%d", base, i)
		sem <- struct{}{}
		innerWg.Add(1)

		go func(ip string) {
			defer innerWg.Done()
			if alive, os := pingWithTTL(ip); alive {
				names, _ := net.LookupAddr(ip)
				hostname := ""
				if len(names) > 0 {
					hostname = names[0]
				}
				results <- Device{IP: ip, Hostname: hostname, OS: os}
			}
			<-sem
		}(ip)
	}

	innerWg.Wait()
}

func askSubnetChoice(subnets []string) []string {
	fmt.Println("Choose a subnet to scan:")
	for i, s := range subnets {
		fmt.Printf("  [%d] %s.0/24\n", i+1, s)
	}
	fmt.Print("  [0] All subnets\n")

	fmt.Print("Enter choice (e.g. 0 or 1,2): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "0" {
		return subnets
	}

	var selected []string
	for _, part := range strings.Split(input, ",") {
		idx := -1
		fmt.Sscanf(part, "%d", &idx)
		if idx > 0 && idx <= len(subnets) {
			selected = append(selected, subnets[idx-1])
		}
	}
	return selected
}

func getMacTable() map[string]string {
	table := make(map[string]string)
	cmd := exec.Command("arp", "-a")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error fetching ARP table:", err)
		return table
	}

	lines := strings.Split(string(output), "\n")
	re := regexp.MustCompile(`(?m)(\d+\.\d+\.\d+\.\d+)\s+([a-fA-F0-9\-]{17})`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			ip := strings.TrimSpace(matches[1])
			mac := strings.ToUpper(strings.ReplaceAll(matches[2], "-", ":"))
			table[ip] = mac
		}
	}
	return table
}

func exportToJSON(devices []Device, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		color.Red("❌ Failed to create JSON file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(devices)
	if err != nil {
		color.Red("❌ Failed to encode JSON: %v", err)
		return
	}
	color.Cyan("📁 Results saved to %s\n", filename)
}

func main() {
	color.Cyan("🔍 Detecting local subnets...")
	subnets := getLocalSubnets()
	if len(subnets) == 0 {
		color.Red("No subnets found.")
		return
	}

	chosen := askSubnetChoice(subnets)

	color.Green("\n🌐 Scanning selected subnet(s)...\n")
	var wg sync.WaitGroup
	results := make(chan Device, 1024)

	for _, subnet := range chosen {
		wg.Add(1)
		go scanSubnet(subnet, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var devices []Device
	for device := range results {
		devices = append(devices, device)
	}

	// Attach MACs
	color.Yellow("\n🔄 Retrieving MAC addresses...")
	macTable := getMacTable()
	for i := range devices {
		if mac, ok := macTable[devices[i].IP]; ok {
			devices[i].MAC = mac
		}
	}

	// Print results
	fmt.Println("\n📋 Active Devices Found:")
	fmt.Println("-------------------------------------------------------------------")
	color.Set(color.FgHiYellow)
	fmt.Printf("%-16s %-30s %-15s %-17s\n", "IP Address", "Hostname", "OS", "MAC Address")
	color.Unset()
	fmt.Println("-------------------------------------------------------------------")
	for _, d := range devices {
		fmt.Printf("%-16s %-30s %-15s %-17s\n", d.IP, d.Hostname, d.OS, d.MAC)
	}
	fmt.Println("-------------------------------------------------------------------")
	color.Green("✅ Done. %d device(s) detected.\n", len(devices))

	// Export to JSON
	exportToJSON(devices, "results.json")
}
