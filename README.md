# NetScanner - Network Scanner in Go

![netscanner](https://img.shields.io/badge/language-Go-blue?logo=go)

A concurrent network scanner built with Go that detects active devices on your local subnets. It identifies devices’ IP addresses, hostnames, operating systems (estimated via TTL), and MAC addresses, then exports results to a JSON file.

---

## Features

- Detects local subnets automatically
- Scans devices in subnet concurrently for faster results
- Estimates device OS based on TTL values
- Resolves hostnames via reverse DNS
- Retrieves MAC addresses from ARP tables
- Exports device info to `results.json`
- Cross-platform (Windows, Linux, macOS)

---

## How It Works

1. Detects local network subnets by reading system network interfaces.
2. Prompts user to select subnet(s) or scan all.
3. Pings IPs within selected subnet(s) concurrently.
4. Parses TTL from ping response to estimate OS.
5. Resolves hostname via reverse DNS lookup.
6. Fetches MAC addresses from ARP table.
7. Displays results and exports to JSON.

---

## Prerequisites

- Go installed (for building from source): [Install Go](https://golang.org/doc/install)
- Internet connectivity on your machine (for network scanning)
- Administrative privileges may be required on some OS to access ARP tables.

---

## Cloning and Running

### 1. Clone the repo

```bash
git clone https://github.com/yourusername/netscanner.git
cd netscanner
```

### 2. Run the code

```bash
go run subnet_scanner.go
```
---

## Example Terminal Session

```bash
$ ./netscanner
🔍 Detecting local subnets...
Choose a subnet to scan:
  [1] 192.168.1.0/24
  [2] 10.0.0.0/24
  [0] All subnets
Enter choice (e.g. 0 or 1,2): 1

🌐 Scanning selected subnet(s)...

🔄 Retrieving MAC addresses...

📋 Active Devices Found:
-------------------------------------------------------------------
IP Address       Hostname                       OS              MAC Address       
-------------------------------------------------------------------
192.168.1.1      router.local                   Windows         00:11:22:33:44:55
192.168.1.10     device1.local                  Linux/macOS     AA:BB:CC:DD:EE:FF
192.168.1.25                                     Windows         11:22:33:44:55:66
-------------------------------------------------------------------
✅ Done. 3 device(s) detected.

📁 Results saved to results.json
```

---

## Output Files

- `results.json` — JSON file containing an array of detected devices, with fields:

```json
[
  {
    "ip": "192.168.1.10",
    "hostname": "device1.local",
    "os": "Linux/macOS",
    "mac": "AA:BB:CC:DD:EE:FF"
  },
  ...
]
```

---

## How to Use the Precompiled Binaries
Simply download the binary for your platform from the Release page and run it from your terminal or command prompt.

### Running on Linux/macOS
```bash
./netscanner
```
### Running on Windows (PowerShell)
```bash
.\netscanner.exe
```

---

## Notes

- On some systems, administrative rights may be needed to run the `arp` command or to ping successfully.
- Estimated OS is based on common TTL values but might not always be accurate.
- Hostnames may be empty if reverse DNS is not set up on the network.
- The scanner pings IPs from `.1` to `.254` in the subnet.

---
