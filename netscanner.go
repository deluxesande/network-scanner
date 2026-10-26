package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/deluxesande/network-scanner/netscanner/ping"
	netscanner "github.com/deluxesande/network-scanner/netscanner/utils"
	"github.com/deluxesande/network-scanner/subnet"
	"github.com/deluxesande/network-scanner/utils"

	"github.com/fatih/color"
)

func main() {
	if err := run(); err != nil {
		color.Red("❌ %s", err)
		os.Exit(1)
	}
}

func run() error {
	subnetFlag := flag.String("subnet", "", "Comma-separated list of subnets to scan (e.g., 192.168.1.0/24,10.0.0.0/24)")
	credits := flag.Bool("credits", false, "Show program credits")
	output := flag.String("output", "", "Output file for JSON results")
	openTcp := flag.Bool("tcp", false, "Scan for open TCP ports (default: false)")
	openUdp := flag.Bool("udp", false, "Scan for open UDP ports (default: false)")
	version := flag.Bool("version", false, "Show version information")
	ipaddr := flag.String("ip", "", "Resolve a website to its IP address (e.g., --ip google.com)")
	concurrency := flag.Int("concurrency", 100, "Hosts probed in parallel")
	timeout := flag.Duration("timeout", time.Second, "Per-host probe timeout")
	noDNS := flag.Bool("no-dns", false, "Skip reverse DNS lookups")
	force := flag.Bool("force", false, "Allow scanning non-private ranges")

	flag.Usage = netscanner.DisplayHelp
	flag.Parse()

	// 1. Informational flags
	if *credits {
		netscanner.DisplayCredits()
		return nil
	}
	if *version {
		netscanner.DisplayVersion()
		return nil
	}

	// 2. Port and resolve modes read the positional host/start/end args
	if *openTcp {
		if err := netscanner.ScanTcp(); err != nil {
			return err
		}
	}
	if *openUdp {
		if err := netscanner.ScanUdp(); err != nil {
			return err
		}
	}
	if *ipaddr != "" {
		ip, err := ping.GetIpAddress(*ipaddr)
		if err != nil {
			return err
		}
		color.Green("✅ Found IP: %s", ip)
	}

	// 3. Sweep subnets only when no mode flag was used, or --subnet was given
	if *subnetFlag == "" && (*openTcp || *openUdp || *ipaddr != "") {
		return nil
	}

	var chosen []string
	if *subnetFlag != "" {
		chosen = strings.Split(*subnetFlag, ",")
	} else {
		color.Cyan("🔍 Detecting local subnets...")
		subnets, err := subnet.GetLocalSubnets()
		if err != nil {
			return err
		}
		if len(subnets) == 0 {
			return fmt.Errorf("no subnets found")
		}

		chosen, err = subnet.AskSubnetChoice(subnets)
		if err != nil {
			return err
		}
	}

	color.Green("\n🌐 Scanning selected subnet(s)...\n")

	// Intercepted here only: port scans and the picker keep default Ctrl-C behaviour
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	devices, err := subnet.ScanSubnets(ctx, chosen, utils.ScanOptions{
		Concurrency: *concurrency,
		Timeout:     *timeout,
		ResolveDNS:  !*noDNS,
		AllowPublic: *force,
		OnProgress: func(done, total int) {
			fmt.Printf("\r🔍 %d/%d hosts probed...", done, total)
		},
	})
	fmt.Println()
	if err != nil {
		return err
	}

	netscanner.PrintResults(devices)

	if *output != "" {
		if err := subnet.ExportToJSON(devices, *output); err != nil {
			return err
		}
		color.Cyan("📁 Results saved to %s\n", *output)
	}
	return nil
}
