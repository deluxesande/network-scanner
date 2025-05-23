package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/deluxesande/network-scanner/subnet"
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
  -h, --help          Show this help message and exit
  -v, --version       Show version information and exit
  -o, --output FILE   Specify output file for JSON results (default: results.json)
  -s, --subnet SUBNET Specify a specific subnet to scan (e.g., 192.168.1.0/24)
  -a, --all           Scan all detected subnets
  -m, --mac           Include MAC address resolution in the scan
  -t, --ttl           Display TTL values and estimated OS for each device
  -c, --credits       Display program credits and exit

Examples:
  netscanner -o output.json
  netscanner --subnet 192.168.1.0/24
  netscanner -a -m -o devices.json
  netscanner -t
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

func main() {
	help := flag.Bool("h", false, "Show help message")
	subnetFlag := flag.String("subnet", "", "Comma-separated list of subnets to scan (e.g., 192.168.1.0/24,10.0.0.0/24)")
	credits := flag.Bool("c", false, "Show program credits")
	output := flag.String("output", "", "Output file for JSON results")

	// Define a map for aliases
	aliases := map[string]*string{
		"o": output,
		"s": subnetFlag,
	}

	// Define short flags
	var shortFlags = make(map[string]string)
	for alias := range aliases {
		shortFlags[alias] = ""
		value := shortFlags[alias]
		flag.StringVar(&value, alias, "", fmt.Sprintf("Alias for --%s", alias))
		shortFlags[alias] = value
	}

	flag.Parse()

	// Resolve aliases
	for alias, target := range aliases {
		if shortFlags[alias] != "" {
			*target = shortFlags[alias]
		}
	}

	var chosen []string

	if *help {
		displayHelp()
		return
	}

	if *credits {
		displayCredits()
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
