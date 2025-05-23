package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	netscanner "github.com/deluxesande/network-scanner/netscanner/utils"
	"github.com/deluxesande/network-scanner/subnet"
	"github.com/deluxesande/network-scanner/utils"

	"github.com/fatih/color"
)

func subnet_scan(chosen []string) []utils.Device {
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

	return devices
}

func main() {
	help := flag.Bool("h", false, "Show help message")
	subnetFlag := flag.String("subnet", "", "Comma-separated list of subnets to scan (e.g., 192.168.1.0/24,10.0.0.0/24)")
	credits := flag.Bool("credit", false, "Show program credits")
	output := flag.String("output", "", "Output file for JSON results")
	openTcp := flag.Bool("tcp", false, "Scan for open TCP ports (default: false)")
	openUdp := flag.Bool("udp", false, "Scan for open UDP ports (default: false)")
	version := flag.Bool("version", false, "Show version information")

	flag.Parse()

	var chosen []string

	if *help {
		netscanner.DisplayHelp()
		return
	}

	if *credits {
		netscanner.DisplayCredits()
		return
	}

	if *version {
		netscanner.DisplayVersion()
		return
	}

	if *openTcp && *openUdp {
		// Run both TCP and UDP scans if both flags are provided
		netscanner.ScanTcp()
		netscanner.ScanUdp()
		return // Terminate the program after running both scans
	}

	if *openTcp {
		netscanner.ScanTcp()
		if flag.NFlag() == 1 { // Check if --tcp is the only flag provided
			return
		}
	}

	if *openUdp {
		netscanner.ScanUdp()
		if flag.NFlag() == 1 { // Check if --udp is the only flag provided
			return
		}
	}

	// Determine subnets to scan
	if *subnetFlag != "" {
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

	devices := subnet_scan(chosen)

	// Print results
	netscanner.PrintResults(devices)

	if *output != "" {
		subnet.ExportToJSON(devices, *output)
	}

	// Prompt to exit
	fmt.Println("\nPress Enter to exit the program.")
	bufio.NewReader(os.Stdin).ReadString('\n')
}
