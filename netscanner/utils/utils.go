package netscanner

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/deluxesande/network-scanner/tcp"
	"github.com/deluxesande/network-scanner/udp"
	"github.com/deluxesande/network-scanner/utils"
	"github.com/fatih/color"
)

func DisplayCredits() {
	color.Green(`                                                      
	_                                           
	____   ____| |_   ___  ____ ____ ____  ____   ____  ____ 
   |  _ \ / _  )  _) /___)/ ___) _  |  _ \|  _ \ / _  )/ ___)
   | | | ( (/ /| |__|___ ( (__( ( | | | | | | | ( (/ /| |    
   |_| |_|\____)\___|___/ \____)_||_|_| |_|_| |_|\____)_|    														 
   

	Netscanner: Lightweight CLI tool to scan your local network for active devices,
	detect operating systems via TTL, and export results to JSON.

	Created by deluxesande
	GitHub: https://github.com/deluxesande/net-scanner
	`)
}

// Rendered from the registered flags, so it cannot drift from reality
func DisplayHelp() {
	fmt.Fprint(os.Stdout, `Usage: netscanner [options]

Options:
`)
	flag.CommandLine.SetOutput(os.Stdout)
	flag.PrintDefaults()
	fmt.Fprint(os.Stdout, `
Examples:
  netscanner
  netscanner --subnet 192.168.1.0/24 --output output.json
  netscanner --subnet 10.0.0.0/16 --concurrency 200 --timeout 2s --no-dns
  netscanner --tcp 192.168.1.10 80 100
  netscanner --udp 192.168.1.10 50 100
  netscanner --ip example.com
`)
}

// Set at build time via -ldflags; see .goreleaser.yaml.
var (
	Version   = "dev"
	BuildDate = "unknown"
)

func DisplayVersion() {
	color.Green(`Netscanner
Version: %s
Build Date: %s
Author: Deluxe Sande
Repository: https://github.com/deluxesande/network-scanner
	`, Version, BuildDate)
}

func PrintResults(devices []utils.Device) {
	fmt.Println("\n📋 Active Devices Found:")
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	color.Set(color.FgHiYellow)
	fmt.Printf("%-16s %-26s %-42s %-17s\n", "IP Address", "Hostname", "OS", "MAC Address")
	color.Unset()
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	for _, d := range devices {
		fmt.Printf("%-16s %-26s %-42s %-17s\n", d.IP, d.Hostname, d.OS, d.MAC)
	}
	if len(devices) == 0 {
		color.Red("❌ No devices found.")
	}
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	color.Green("✅ Done. %d device(s) detected.\n", len(devices))
}

// Both port scans read the host and range from the positional args
func parsePortArgs(proto string) (string, int, int, error) {
	args := flag.Args()

	if len(args) < 3 {
		return "", 0, 0, fmt.Errorf("provide a host, start port and end port: netscanner --%s <host> <startPort> <endPort>", proto)
	}

	startPort, errStart := strconv.Atoi(args[1])
	endPort, errEnd := strconv.Atoi(args[2])
	if errStart != nil || errEnd != nil || startPort < 1 || endPort > 65535 || startPort > endPort {
		return "", 0, 0, fmt.Errorf("invalid port range: ports are 1-65535 and startPort must be <= endPort")
	}

	return args[0], startPort, endPort, nil
}

func ScanTcp() error {
	host, startPort, endPort, err := parsePortArgs("tcp")
	if err != nil {
		return err
	}

	color.Cyan("🔍 Scanning for open TCP ports on %s from port %d to %d...", host, startPort, endPort)
	openPorts := tcp.ScanOpenTcpPorts(host, startPort, endPort)

	if len(openPorts) == 0 {
		color.Yellow("⚠️ No open ports found in the specified range.")
		return nil
	}

	color.Green("✅ Open ports found:")
	for port, service := range openPorts {
		fmt.Printf(" - Port %d: %s (Version: %s)\n", port, service.Service, service.Version)
	}
	return nil
}

func ScanUdp() error {
	host, startPort, endPort, err := parsePortArgs("udp")
	if err != nil {
		return err
	}

	color.Cyan("🔍 Scanning for open UDP ports on %s from port %d to %d...", host, startPort, endPort)
	openPorts := udp.ScanOpenUdpPorts(host, startPort, endPort)

	if len(openPorts) == 0 {
		color.Yellow("⚠️ No open UDP ports found in the specified range.")
		return nil
	}

	color.Green("✅ Open UDP ports found:")
	for port, service := range openPorts {
		fmt.Printf(" - Port %d: %s\n", port, service)
	}
	return nil
}
