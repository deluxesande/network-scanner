package udp

import (
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/deluxesande/network-scanner/utils"
)

func identifyUdpService(port int) string {
	// Check if the port is in the map
	if service, exists := utils.UdpServices[port]; exists {
		return service
	}

	// Return "Unknown" if the port is not in the map
	return "Unknown"
}

func scanUdpPort(host string, port int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the WaitGroup counter when the goroutine completes

	address := host + ":" + strconv.Itoa(port)

	conn, err := net.DialTimeout("udp", address, time.Second*1)

	if err != nil {
		return
	}

	defer conn.Close()

	_, err = conn.Write([]byte("ping"))

	if err != nil {
		return
	}

	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	_, err = conn.Read(buffer)

	if err == nil {
		results <- port
	}
}

func ScanOpenUdpPorts(host string, startPort int, endPort int) map[int]string {
	var wg sync.WaitGroup
	results := make(chan int)
	sem := make(chan struct{}, 100) // Limit concurrency to 100 goroutines

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)

		// Acquire a semaphore slot
		sem <- struct{}{}

		go func(port int) {
			defer func() { <-sem }() // Release the semaphore slot
			scanUdpPort(host, port, results, &wg)
		}(port)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	openPorts := make(map[int]string)

	for result := range results {
		service := identifyUdpService(result)
		openPorts[result] = service
	}

	return openPorts
}
