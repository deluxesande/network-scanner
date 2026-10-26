package subnet

import (
	"net/netip"
	"os/exec"
	"strings"
	"testing"

	"github.com/deluxesande/network-scanner/utils"
)

func TestParseTarget(t *testing.T) {
	cases := map[string]string{
		"192.168.1.0/24": "192.168.1.0/24",
		"10.0.0.5/16":    "10.0.0.0/16",
		"192.168.1.10":   "192.168.1.10/32",
	}
	for raw, want := range cases {
		prefix, err := parseTarget(raw)
		if err != nil {
			t.Fatalf("parse %s: %v", raw, err)
		}
		if got := prefix.String(); got != want {
			t.Errorf("%s: got %s, want %s", raw, got, want)
		}
	}

	// Garbage used to panic in the old first-three-octets split
	for _, raw := range []string{"foo", "", "192.168.1.0/33", "::1/128"} {
		if _, err := parseTarget(raw); err == nil {
			t.Errorf("%q should error", raw)
		}
	}
}

func TestEnumerateHosts(t *testing.T) {
	cases := map[string]int{
		"192.168.1.0/24": 254,
		"10.0.0.0/30":    2,
		"10.0.0.0/31":    2,
		"10.0.0.5/32":    1,
	}
	for cidr, want := range cases {
		prefix, err := parseTarget(cidr)
		if err != nil {
			t.Fatalf("parse %s: %v", cidr, err)
		}
		hosts, err := enumerateHosts(prefix)
		if err != nil {
			t.Fatalf("enumerate %s: %v", cidr, err)
		}
		if len(hosts) != want {
			t.Errorf("%s: got %d hosts, want %d", cidr, len(hosts), want)
		}
	}

	if _, err := enumerateHosts(netip.MustParsePrefix("10.0.0.0/8")); err == nil {
		t.Error("/8 should be refused, not swept")
	}
}

func TestIsPrivate(t *testing.T) {
	cases := map[string]bool{
		"192.168.1.0/24": true,
		"10.0.0.0/8":     true,
		"172.16.0.0/12":  true,
		"127.0.0.0/8":    true,
		"8.8.8.0/24":     false,
	}
	for cidr, want := range cases {
		if got := isPrivate(netip.MustParsePrefix(cidr)); got != want {
			t.Errorf("%s: got %v, want %v", cidr, got, want)
		}
	}
}

func TestEstimateOS(t *testing.T) {
	cases := map[string]string{
		"":    "Unknown",
		"0":   "Unknown",
		"abc": "Unknown",
		"300": "Unknown",
		"32":  "Windows 95/98/ME",
		"63":  "Linux", // 64 minus one hop
		"64":  "Linux",
		"127": "Windows NT", // 128 minus one hop
		"128": "Windows NT",
		"255": "Cisco IOS",
	}
	for ttl, want := range cases {
		if got := EstimateOS(ttl); !strings.Contains(got, want) {
			t.Errorf("ttl %q: got %q, want it to contain %q", ttl, got, want)
		}
	}

	if got := EstimateOS("63"); !strings.Contains(got, "1 hop") {
		t.Errorf("ttl 63 should report the hop count: %q", got)
	}
}

func TestParseArpLines(t *testing.T) {
	output := "? (192.168.1.1) at 0:1c:42:0:0:8 on en0 ifscope [ethernet]\n" +
		"192.168.1.5      00-11-22-33-44-55     dynamic\n" +
		"gateway (192.168.1.254) at 00:11:22:33:44:55 [ether] on wlan0\n"
	table := parseArpLines(output)

	if table["192.168.1.1"] != "00:1C:42:00:00:08" {
		t.Errorf("darwin single-digit MAC: got %q", table["192.168.1.1"])
	}
	if table["192.168.1.5"] != "00:11:22:33:44:55" {
		t.Errorf("windows dashed MAC: got %q", table["192.168.1.5"])
	}
	if table["192.168.1.254"] != "00:11:22:33:44:55" {
		t.Errorf("linux MAC: got %q", table["192.168.1.254"])
	}
}

func TestNormalizeMac(t *testing.T) {
	cases := map[string]string{
		"0:1c:42:0:0:8":     "00:1C:42:00:00:08",
		"00-11-22-33-44-55": "00:11:22:33:44:55",
		"aa:bb:cc:dd:ee:ff": "AA:BB:CC:DD:EE:FF",
	}
	for raw, want := range cases {
		if got := normalizeMac(raw); got != want {
			t.Errorf("%q: got %q, want %q", raw, got, want)
		}
	}
}

func TestMergeDevicesFiltersAndSorts(t *testing.T) {
	pinged := []utils.Device{
		{IP: "192.168.1.20", OS: "Linux/macOS"},
		{IP: "192.168.1.3", OS: "Windows"},
	}
	macs := map[string]string{
		"192.168.1.20": "00:11:22:33:44:55",
		"192.168.1.7":  "AA:BB:CC:DD:EE:FF", // ICMP filtered, still on the wire
		"192.168.1.1":  "00:11:22:33:44:66", // gateway, inside the prefix
		"192.168.2.10": "00:11:22:33:44:77", // neighbour, outside the target
	}
	prefixes := []netip.Prefix{netip.MustParsePrefix("192.168.1.0/24")}

	devices := mergeDevices(pinged, macs, prefixes)
	if len(devices) != 4 {
		t.Fatalf("got %d devices, want 4: %+v", len(devices), devices)
	}

	wantOrder := []string{"192.168.1.1", "192.168.1.3", "192.168.1.7", "192.168.1.20"}
	for i, want := range wantOrder {
		if devices[i].IP != want {
			t.Errorf("position %d: got %s, want %s", i, devices[i].IP, want)
		}
	}
	if devices[2].OS != "Unknown" || devices[2].MAC != "AA:BB:CC:DD:EE:FF" {
		t.Errorf("ARP-only host not filled in: %+v", devices[2])
	}
}

// Live check: needs a working ping binary and loopback
func TestPingWithTTL(t *testing.T) {
	if _, err := exec.LookPath("ping"); err != nil {
		t.Skip("ping not installed")
	}
	alive, _ := PingWithTTL("127.0.0.1")
	if !alive {
		t.Fatal("loopback did not answer ping")
	}
}
