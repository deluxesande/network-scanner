package ping

import (
	"fmt"
	"net"
)

func GetIpAddress(website string) (string, error) {
	if website == "" {
		return "", fmt.Errorf("please provide a website to resolve its IP address")
	}

	ips, err := net.LookupIP(website)
	if err != nil {
		return "", fmt.Errorf("failed to resolve %s: %w", website, err)
	}

	// Prefer IPv4; the rest of the scanner is IPv4 only
	for _, ip := range ips {
		if v4 := ip.To4(); v4 != nil {
			return v4.String(), nil
		}
	}

	if len(ips) > 0 {
		return ips[0].String(), nil
	}

	return "", fmt.Errorf("no IP address found for %s", website)
}
