package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/deluxesande/network-scanner/subnet"
	"github.com/deluxesande/network-scanner/tcp"
	"github.com/deluxesande/network-scanner/utils"

	"github.com/fatih/color"
)

func displayCredits() {
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

func displayHelp() {
	color.Green(`Usage: netscanner [options]

Options:
  -h,          Show this help message and exit
  --version       Show version information and exit
  --output FILE   Specify output file for JSON results (default: results.json)
  --subnet SUBNET Specify a specific subnet to scan (e.g., 192.168.1.0/24)
  --credits       Display program credits and exit

Examples:
  netscanner --output output.json
  netscanner --subnet 192.168.1.0/24
  `)
}

func printResults(devices []utils.Device) {
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

func scanTcp() {
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

func main() {
	help := flag.Bool("h", false, "Show help message")
	subnetFlag := flag.String("subnet", "", "Comma-separated list of subnets to scan (e.g., 192.168.1.0/24,10.0.0.0/24)")
	credits := flag.Bool("credit", false, "Show program credits")
	output := flag.String("output", "", "Output file for JSON results")
	openTcp := flag.Bool("tcp", false, "Scan for open TCP ports (default: false)")

	flag.Parse()

	var chosen []string

	if *help {
		displayHelp()
		return
	}

	if *credits {
		displayCredits()
		return
	}

	if *openTcp {
		scanTcp()
		return
	}

	// Determine subnets to scan
	if *subnetFlag != "" {
		color.Cyan("🔍 Using provided subnets...")
		chosen = strings.Split(*subnetFlag, ",")
	} else {
		color.Cyan("🔍 Detecting local subnets...")
		subnets := subnet.GetLocalSubnets()
		if len(subnets) == 0 {
			color.Red("No subnets found.")
			return
		}
		chosen = subnet.AskSubnetChoice(subnets)
	}

	color.Green("\n🌐 Scanning selected subnet(s)...\n")
	var wg sync.WaitGroup
	results := make(chan utils.Device, 1024)

	for _, subnetChoice := range chosen {
		wg.Add(1)
		go subnet.ScanSubnet(subnetChoice, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var devices []utils.Device
	for device := range results {
		devices = append(devices, device)
	}

	// Attach MACs
	color.Yellow("\n🔄 Retrieving MAC addresses...")
	macTable := subnet.GetMacTable()
	for i := range devices {
		if mac, ok := macTable[devices[i].IP]; ok {
			devices[i].MAC = mac
		}
	}

	// Print results
	printResults(devices)

	if *output != "" {
		subnet.ExportToJSON(devices, *output)
	} else {
		subnet.ExportToJSON(devices, "results.json")
	}

	// Prompt to exit
	fmt.Println("\nPress Enter to exit the program.")
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println("👋 Exiting the program. Goodbye!")
}
