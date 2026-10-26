# NetScanner 

![netscanner](https://img.shields.io/badge/language-Go-blue?logo=go)

A concurrent network scanner built with Go that detects active devices on your local subnets. It identifies devices’ IP addresses, hostnames, operating systems (estimated via TTL), and MAC addresses, then exports results to a JSON file.


## Features

- Detects local subnets automatically
- Scans devices in subnet concurrently for faster results
- Estimates device OS based on TTL values
- Resolves hostnames via reverse DNS
- Retrieves MAC addresses from ARP tables
- **Scans specific ports to detect open services**
- Exports device info to `results.json`
- Cross-platform (Windows, Linux, macOS)

## How It Works

1. Detects local network subnets by reading system network interfaces.
2. Prompts user to select subnet(s) or scan all.
3. Pings every host in the selected subnet(s) concurrently.
4. Parses TTL from ping response to estimate OS, including how many hops away it is.
5. Resolves hostname via reverse DNS lookup.
6. Merges the ARP cache, so hosts that drop ICMP still appear with their MAC address.
7. **Scans specified ports to detect open services.**
8. Displays results and exports to JSON.

## Prerequisites

- A machine connected to a local network (Wi-Fi or Ethernet). Internet access is not needed to scan.
- `ping` must be available on your PATH. On minimal Linux installs:
  - Debian/Ubuntu: `sudo apt install iputils-ping`
  - Fedora: `sudo dnf install iputils`
  - macOS and Windows ship `ping` already.
- Go 1.24 or newer, only if you build from source.
- Administrator/root rights may be needed on some systems to use `ping`.

## Help and Usage

The `netscanner` CLI tool provides a `-h` option to display detailed usage instructions and examples. This is implemented using the `DisplayHelp` function, which outputs the following information:

### Usage
```bash
Usage: netscanner [options]

Options:
  -h                Show this help message and exit
  --version         Show version information and exit
  --output FILE     Specify output file for JSON results
  --subnet SUBNETS  Subnet(s) to scan, comma-separated (e.g., 192.168.1.0/24,10.0.0.0/24)
  --tcp HOST STARTPORT ENDPORT
                    Scan for open TCP ports on a specific host
  --udp HOST STARTPORT ENDPORT
                    Scan for open UDP ports on a specific host
  --ip WEBSITE      Resolve the IP address of a website (e.g., google.com)
  --concurrency N   Hosts probed in parallel (default: 100)
  --timeout DURATION
                    Per-host probe timeout (default: 1s)
  --no-dns          Skip reverse DNS lookups
  --force           Allow scanning non-private ranges
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

#### Scan UDP ports 50 to 100 on a specific host
```bash
netscanner --udp 192.168.1.10 50 100
```
#### Resolve a website's IP address
```bash
netscanner --ip google.com
```

---

## Cloning and Running

This section explains how to clone the repository and run the netscanner application using Go. It provides simple steps to get the code and execute it directly from the source.


```bash
git clone https://github.com/deluxesande/network-scanner.git
cd network-scanner
```

```bash
go run . --subnet 192.168.1.0/24
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

🔍 254/254 hosts probed...
📋 Active Devices Found:
--------------------------------------------------------------------------------------------------------
IP Address       Hostname                   OS                                         MAC Address
--------------------------------------------------------------------------------------------------------
192.168.1.1      router.local               Linux/macOS/iOS/Android/BSD (~1 hop)       00:11:22:33:44:55
192.168.1.10     device1.local              Linux/macOS/iOS/Android/BSD                AA:BB:CC:DD:EE:FF
192.168.1.25                                Windows NT and later                       11:22:33:44:55:66
--------------------------------------------------------------------------------------------------------
✅ Done. 3 device(s) detected.

📁 Results saved to results.json
</pre>

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

## Using as a Go Module

```bash
go get github.com/deluxesande/network-scanner@latest
```

```go
package main

import (
    "fmt"
    "log"

    "github.com/deluxesande/network-scanner/subnet"
    "github.com/deluxesande/network-scanner/tcp"
)

func main() {
    subnets, err := subnet.GetLocalSubnets()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(subnets)
    fmt.Println(tcp.ScanOpenTcpPorts("192.168.1.10", 80, 100))
}
```

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

The one-line installers download the latest release for your OS and CPU (x86_64 or arm64), put the `netscanner` binary in a user folder, and add that folder to your `PATH`. Open a new terminal after installing.

- **Windows**: runs `netscanner_install.ps1` in PowerShell and installs to `%USERPROFILE%\netscanner`.
- **Linux and macOS**: runs `netscanner_install.sh`, which detects your OS and architecture, and installs to `~/.local/bin`.

These commands ensure that the `netscanner` binary is installed system-wide and can be run from anywhere on the system without requiring manual setup.

#### Windows (PowerShell)
```powershell
$script = "$env:TEMP\netscanner_install.ps1"; Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/deluxesande/network-scanner/main/install/netscanner_install.ps1' -OutFile $script; & $script; Remove-Item $script
```

#### Linux and macOS (Intel and Apple Silicon, x86_64 and arm64)
```bash
curl -fsSL https://raw.githubusercontent.com/deluxesande/network-scanner/main/install/netscanner_install.sh | bash
```

## Uninstalling the binaries

This section provides commands to uninstall the netscanner binary from your system. The commands remove the binary from its installation directory and ensure it is no longer included in the system's PATH. After running the commands, restart your terminal to apply the changes.

#### Windows
```powershell
$installPath = "$env:USERPROFILE\netscanner"
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$newPath = ($userPath -split ";" | Where-Object { $_ -and $_ -ne $installPath }) -join ";"
[Environment]::SetEnvironmentVariable("Path", $newPath, "User")
Remove-Item -Recurse -Force $installPath
Write-Host "netscanner uninstalled successfully."
```

#### Linux and macOS

```bash
rm -f "$HOME/.local/bin/netscanner"
for rc in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile"; do
  [ -f "$rc" ] && sed -i.bak '/# added by netscanner installer/d' "$rc" && rm -f "$rc.bak"
done
echo "netscanner uninstalled. Restart your terminal."
```

If you installed with an older version of the installer, also delete the `export PATH="$HOME/.local/bin:$PATH"` line from your shell rc file by hand.

## Notes

- On some systems, administrative rights may be needed to run the `arp` command or to ping successfully.
- Estimated OS is based on initial TTL values; lower TTLs are reported with an approximate hop count.
- Hostnames may be empty if reverse DNS is not set up on the network.
- The scanner sweeps every usable host in the subnet, and merges the ARP cache so hosts that drop ICMP still appear.
- Ranges outside RFC1918 need `--force`. Anything wider than a `/16` is refused.

## Contributing

Contributions are welcome! If you'd like to contribute to this project, please follow these steps:

- Fork the repository.
- Create a new branch for your feature or bug fix.
- Commit your changes and push them to your fork.
- Submit a pull request with a detailed description of your changes.

For major changes, please open an issue first to discuss what you would like to change.

## Changelog

All notable changes to this project will be documented in the `CHANGELOG.md` file. Please refer to it for details on updates, fixes, and new features.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
