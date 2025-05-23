# NetScanner 

``` bash                                                    
	_                                           
	____   ____| |_   ___  ____ ____ ____  ____   ____  ____ 
   |  _ \ / _  )  _) /___)/ ___) _  |  _ \|  _ \ / _  )/ ___)
   | | | ( (/ /| |__|___ ( (__( ( | | | | | | | ( (/ /| |    
   |_| |_|\____)\___|___/ \____)_||_|_| |_|_| |_|\____)_|    														 
   

	Netscanner: Lightweight CLI tool to scan your local network for active devices,
	detect operating systems via TTL, and export results to JSON.

	Created by deluxesande
	GitHub: https://github.com/deluxesande/net-scanner
```

![netscanner](https://img.shields.io/badge/language-Go-blue?logo=go)

A concurrent network scanner built with Go that detects active devices on your local subnets. It identifies devices’ IP addresses, hostnames, operating systems (estimated via TTL), and MAC addresses, then exports results to a JSON file.


---

## Features

- Detects local subnets automatically
- Scans devices in subnet concurrently for faster results
- Estimates device OS based on TTL values
- Resolves hostnames via reverse DNS
- Retrieves MAC addresses from ARP tables
- **Scans specific ports to detect open services**
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
7. **Scans specified ports to detect open services.**
8. Displays results and exports to JSON.

---

## Prerequisites

- Go installed (for building from source): [Install Go](https://golang.org/doc/install)
- Internet connectivity on your machine (for network scanning)
- Administrative privileges may be required on some OS to access ARP tables.

---

## Help and Usage

The `netscanner` CLI tool provides a `--help` option (or `-h`) to display detailed usage instructions and examples. This is implemented using the `DisplayHelp` function, which outputs the following information:

### Usage
Usage: netscanner [options]

Options:
  -h                Show this help message and exit
  --version         Show version information and exit
  --output FILE     Specify output file for JSON results
  --subnet SUBNET   Specify a specific subnet to scan (e.g., 192.168.1.0/24)
  --tcp HOST STARTPORT ENDPORT
                    Scan for open TCP ports on a specific host
  --credits         Display program credits and exit
```

#### Scan TCP ports 80 to 100 on a specific host
```bash
netscanner --tcp 192.168.1.10 80 100
```
#### Scan a specific subnet and save results to a custom output file
```bash
netscanner --subnet 192.168.1.0/24 --output output.json
```

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

<pre>
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
</pre>

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

## Using as a Go Module

You can import this project into your Go application as a module. Add the following import statement to your Go code:

```go
import "github.com/yourusername/netscanner/netscanner"
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
## Install binaries into the system

This section provides instructions on how to run the `netscanner` binary on different operating systems and how to install it system-wide using one-line commands. The commands are tailored for Windows, Linux, and macOS, ensuring that the `netscanner` binary is downloaded, executed, and added to the system's `PATH` for easy access.

*The scripts accounts for x86_64 systems*

### Running the Binary
- **Linux/macOS**: The binary can be executed directly using `./netscanner`.
- **Windows**: The binary can be executed using `.\netscanner.exe` in PowerShell.

### Installing the Binary
The installation commands automate the process of downloading the appropriate installer script from the GitHub repository, running the script, and cleaning up after execution:
- **Windows**: Uses a PowerShell command to download and execute the `netscanner_install.ps1` script.
- **Linux**: Uses `curl` to fetch the `netscanner_install.sh` script, executes it, and removes it afterward.
- **macOS**: Similar to Linux, but uses the macOS-specific `netscanner_install_mac.sh` script.

These commands ensure that the `netscanner` binary is installed system-wide and can be run from anywhere on the system without requiring manual setup.

```bash
$script = "$env:TEMP\netscanner_install.ps1"; Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/deluxesande/network-scanner/main/install/netscanner_install.ps1' -OutFile $script; & $script; Remove-Item $script
```

```bash
curl -s https://raw.githubusercontent.com/deluxesande/network-scanner/main/install/netscanner_install.sh -o /tmp/netscanner_install.sh && bash /tmp/netscanner_install.sh && rm /tmp/netscanner_install.sh
```

```bash
curl -s https://raw.githubusercontent.com/deluxesande/network-scanner/main/install/netscanner_install_mac.sh -o /tmp/netscanner_install_mac.sh && bash /tmp/netscanner_install_mac.sh && rm /tmp/netscanner_install_mac.sh
```

---

## Notes

- On some systems, administrative rights may be needed to run the `arp` command or to ping successfully.
- Estimated OS is based on common TTL values but might not always be accurate.
- Hostnames may be empty if reverse DNS is not set up on the network.
- The scanner pings IPs from `.1` to `.254` in the subnet.

---

## Contributing

Contributions are welcome! If you'd like to contribute to this project, please follow these steps:

- Fork the repository.
- Create a new branch for your feature or bug fix.
- Commit your changes and push them to your fork.
- Submit a pull request with a detailed description of your changes.

For major changes, please open an issue first to discuss what you would like to change.

---

## Changelog

All notable changes to this project will be documented in the `CHANGELOG.md` file. Please refer to it for details on updates, fixes, and new features.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
