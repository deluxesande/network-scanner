package netscanner

import (
	"flag"
	"fmt"
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

func DisplayHelp() {
	color.Green(`Usage: netscanner [options]

Options:
  -h,          Show this help message and exit
  --version       Show version information and exit
  --output FILE   Specify output file for JSON results
  --subnet SUBNET Specify a specific subnet to scan (e.g., 192.168.1.0/24)
  --tcp <host> <port> <port> HOST STARTPORT ENDPORT Scan for open TCP ports on a specific host
  --credits       Display program credits and exit

Examples:
  netscanner --tcp 192.168.1.10 80 100
  netscanner --subnet 192.168.1.0/24 --output output.json
  `)
}

func DisplayVersion() {
	color.Green(`Netscanner
Version: 0.2.0
Build Date: 2025-05-23
Author: Deluxe Sande
Repository: https://github.com/deluxesande/network-scanner
	`)
}

func PrintResults(devices []utils.Device) {
	fmt.Println("\n📋 Active Devices Found:")
	fmt.Println("----------------------------------------------------------------------------------")
	color.Set(color.FgHiYellow)
	fmt.Printf("%-16s %-30s %-15s %-17s\n", "IP Address", "Hostname", "OS", "MAC Address")
	color.Unset()
	fmt.Println("----------------------------------------------------------------------------------")
	for _, d := range devices {
		fmt.Printf("%-16s %-30s %-15s %-17s\n", d.IP, d.Hostname, d.OS, d.MAC)
	}
	if len(devices) == 0 {
		color.Red("❌ No devices found.")
	}
	fmt.Println("----------------------------------------------------------------------------------")
	color.Green("✅ Done. %d device(s) detected.\n", len(devices))
}

func ScanTcp() {
	// Ensure the user has provided the required arguments
	args := flag.Args()

	if len(args) < 3 {
		color.Red("❌ Please provide the host, start port, and end port for TCP scanning.")
		fmt.Println("Usage: netscanner --tcp <host> <startPort> <endPort>")
		return
	}
	host := args[0]
	startPort, err1 := strconv.Atoi(args[1])
	endPort, err2 := strconv.Atoi(args[2])

	// Validate the ports
	if err1 != nil || err2 != nil || startPort < 1 || endPort > 65535 || startPort > endPort {
		color.Red("❌ Error: Invalid port range. Ports must be integers between 1 and 65535, and startPort must be <= endPort.")
		return
	}

	// Perform the TCP scan
	color.Cyan("🔍 Scanning for open TCP ports on %s from port %d to %d...", host, startPort, endPort)
	openPorts := tcp.ScanOpenTcpPorts(host, startPort, endPort)

	// Display the results
	if len(openPorts) > 0 {
		color.Green("✅ Open ports found:")
		for port, service := range openPorts {
			fmt.Printf(" - Port %d: %s (Version: %s)\n", port, service.Service, service.Version)
		}
	} else {
		color.Yellow("⚠️ No open ports found in the specified range.")
	}
}

func ScanUdp() {
	args := flag.Args()

	if len(args) < 3 {
		color.Red("❌ Please provide the host, start port, and end port for UDP scanning.")
		fmt.Println("Usage: netscanner --udp <host> <startPort> <endPort>")
		return
	}

	host := args[0]
	startPort, err1 := strconv.Atoi(args[1])
	endPort, err2 := strconv.Atoi(args[2])

	// Validate the ports
	if err1 != nil || err2 != nil || startPort < 1 || endPort > 65535 || startPort > endPort {
		color.Red("❌ Error: Invalid port range. Ports must be integers between 1 and 65535, and startPort must be <= endPort.")
		return
	}

	color.Cyan("🔍 Scanning for open UDP ports on %s from port %d to %d...", host, startPort, endPort)
	openPorts := udp.ScanOpenUdpPorts(host, startPort, endPort)

	if len(openPorts) > 0 {
		color.Green("✅ Open UDP ports found:")
		for port, service := range openPorts {
			fmt.Printf(" - Port %d: %s\n", port, service)
		}
	} else {
		color.Yellow("⚠️ No open UDP ports found in the specified range.")
	}
}
