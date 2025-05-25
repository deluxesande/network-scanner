package ping

import (
	"flag"
	"fmt"
	"net"

	"github.com/fatih/color"
)

func GetIpAddress() (string, error) {
	args := flag.Args()

	if len(args) < 1 {
		return "", fmt.Errorf("%s", color.RedString("❌ Please provide a website to resolve its IP address"))
	}

	website := args[0]

	ips, err := net.LookupIP(website)

	if err != nil {
		return "", fmt.Errorf("failed to resolve IP for website %s: %v", website, err)
	}

	if len(ips) > 0 {
		return ips[0].String(), nil
	}

	return "", fmt.Errorf("no IP address found for %s", website)
}
