package tcp

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/deluxesande/network-scanner/utils"
)

func scanPort(host string, port int, results chan<- map[int]utils.ServiceInfo, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the WaitGroup counter when the goroutine completes

	address := host + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", address, time.Second*1)

	if err == nil {
		defer conn.Close()

		// Simulate banner grabbing (you can replace this with actual banner grabbing logic)
		banner := grabBanner(host, port)

		// Identify the service
		service, version := identifyService(port, banner)

		// Send the open port to the channel
		results <- map[int]utils.ServiceInfo{port: {Service: service, Version: version}}
	}
}

func grabBanner(host string, port int) string {
	address := net.JoinHostPort(host, strconv.Itoa(port))

	conn, err := net.DialTimeout("tcp", address, time.Second*2)
	if err != nil {
		return ""
	}

	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	// Read the first 1024 bytes of the banner
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		return ""
	}

	// Convert the bytes to a string
	return string(buffer[:n])
}

func identifyService(port int, banner string) (string, string) {
	// Try to identify the service based on the port number
	service := "Unknown"

	if s, exists := utils.CommonPorts[port]; exists {
		service = s
	}

	version := "Unknown"

	lowerBanner := strings.ToLower(banner)

	// SSH version detection
	if strings.Contains(lowerBanner, "ssh") {
		service = "SSH"
		parts := strings.Split(lowerBanner, " ")

		if len(parts) >= 2 {
			version = parts[1]
		}
	}

	// HTTP version detection
	if strings.Contains(lowerBanner, "http") || strings.Contains(lowerBanner, "apache") || strings.Contains(lowerBanner, "nginx") {
		if port == 443 {
			service = "HTTPS"
		} else {
			service = "HTTP"
		}

		// Try to find server info in format "Server: Apache/2.4.29"
		if strings.Contains(lowerBanner, "Server:") {
			parts := strings.Split(lowerBanner, "Server:")

			if len(parts) > 2 {
				version = strings.TrimSpace(parts[1])
			} else {
				parts = strings.Split(lowerBanner, " ")

				if len(parts) > 1 {
					version = strings.TrimSpace(parts[1])
				}
			}
		}

	}

	return service, version
}

func ScanOpenTcpPorts(host string, startPort, endPort int) map[int]utils.ServiceInfo {
	var wg sync.WaitGroup
	results := make(chan map[int]utils.ServiceInfo, endPort-startPort+1) // Buffered channel to hold open ports
	sem := make(chan struct{}, 100)                                      // Semaphore to limit concurrent goroutines

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)

		// Acquire a semaphore slot
		sem <- struct{}{}

		go func(port int) {
			defer func() { <-sem }() // Release the semaphore slot
			scanPort(host, port, results, &wg)
		}(port)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(results) // Close the results channel when done
	}()

	// Collect open ports from the results channel
	openPorts := make(map[int]utils.ServiceInfo)
	for result := range results {
		for port, info := range result {
			openPorts[port] = info
		}
	}

	return openPorts
}
