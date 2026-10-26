package tcp

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/deluxesande/network-scanner/utils"
)

const bannerTimeout = 2 * time.Second

// Ports whose server waits for a request before saying anything
var httpProbePorts = map[int]bool{80: true, 8000: true, 8080: true, 8888: true}

func scanPort(host string, port int, results chan<- utils.PortResult, wg *sync.WaitGroup) {
	defer wg.Done()

	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	banner := grabBanner(conn, port)
	service, version := identifyService(port, banner)
	results <- utils.PortResult{Port: port, Service: utils.ServiceInfo{Service: service, Version: version}}
}

// Reuse the open connection; a second dial doubled the scan time
func grabBanner(conn net.Conn, port int) string {
	conn.SetDeadline(time.Now().Add(bannerTimeout))

	if httpProbePorts[port] {
		conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return ""
	}
	return string(buffer[:n])
}

func identifyService(port int, banner string) (string, string) {
	service := "Unknown"
	if s, exists := utils.CommonPorts[port]; exists {
		service = s
	}

	version := "Unknown"
	lowerBanner := strings.ToLower(banner)

	// SSH speaks first: SSH-2.0-OpenSSH_8.9p1
	if strings.HasPrefix(lowerBanner, "ssh-") {
		service = "SSH"
		if fields := strings.Fields(banner); len(fields) > 0 {
			version = strings.TrimPrefix(strings.TrimPrefix(fields[0], "SSH-2.0-"), "SSH-1.99-")
		}
		return service, version
	}

	if strings.Contains(lowerBanner, "http") || strings.Contains(lowerBanner, "apache") || strings.Contains(lowerBanner, "nginx") {
		if port == 443 {
			service = "HTTPS"
		} else {
			service = "HTTP"
		}

		// Server header of a HEAD response, sliced from the original for its casing
		if idx := strings.Index(lowerBanner, "server:"); idx >= 0 {
			value := banner[idx+len("server:"):]
			if end := strings.IndexAny(value, "\r\n"); end >= 0 {
				value = value[:end]
			}
			version = strings.TrimSpace(value)
		}
		return service, version
	}

	return service, version
}

func ScanOpenTcpPorts(host string, startPort, endPort int) map[int]utils.ServiceInfo {
	var wg sync.WaitGroup
	results := make(chan utils.PortResult, 256)
	sem := make(chan struct{}, 100) // Semaphore to limit concurrent goroutines

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
		openPorts[result.Port] = result.Service
	}

	return openPorts
}
