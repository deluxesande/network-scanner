// package main

// import (
// 	"fmt"
// 	"net"
// 	"os/exec"
// 	"strings"
// 	"sync"
// )

// const (
// 	subnet      = "192.168.100." // Change this to your local subnet
// 	startIP     = 1
// 	endIP       = 254
// 	concurrency = 100
// )

// // ping sends a single ping to the IP and checks if it responds
// func ping(ip string) bool {
// 	// -n 1 = send 1 echo request, -w 1000 = wait 1s max
// 	cmd := exec.Command("ping", "-n", "1", "-w", "1000", ip)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return false
// 	}
// 	return strings.Contains(string(output), "Reply from")
// }

// // scanNetwork pings all IPs in the subnet and returns the ones that are active
// func scanNetwork() []string {
// 	var wg sync.WaitGroup
// 	var mu sync.Mutex
// 	activeIPs := []string{}
// 	sem := make(chan struct{}, concurrency)

// 	for i := startIP; i <= endIP; i++ {
// 		ip := fmt.Sprintf("%s%d", subnet, i)
// 		wg.Add(1)
// 		sem <- struct{}{}

// 		go func(ip string) {
// 			defer wg.Done()
// 			if ping(ip) {
// 				mu.Lock()
// 				activeIPs = append(activeIPs, ip)
// 				mu.Unlock()
// 			}
// 			<-sem
// 		}(ip)
// 	}

// 	wg.Wait()
// 	return activeIPs
// }

// func main() {
// 	fmt.Println("Scanning network...")
// 	devices := scanNetwork()

// 	fmt.Println("Active devices found:")
// 	for _, ip := range devices {
// 		names, _ := net.LookupAddr(ip)
// 		if len(names) > 0 {
// 			fmt.Printf("%s (%s)\n", ip, names[0])
// 		} else {
// 			fmt.Println(ip)
// 		}
// 	}
// }
