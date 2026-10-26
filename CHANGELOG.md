# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

- N/A

---

## [v0.3.2] - 2026-10-26

### Added
- `--concurrency`, `--timeout`, `--no-dns` and `--force` flags for controlling subnet sweeps.
- Real CIDR support: any IPv4 prefix up to `/16` is swept correctly, and a bare IP is treated as `/32`.
- Hosts that drop ICMP are still reported, merged in from the kernel ARP cache.
- Host discovery reports the hop count when a TTL is below its OS default.
- `--help` renders the live flag list instead of a hand-maintained copy.
- CI workflow (vet, race tests, build) and a goreleaser release workflow.

### Changed
- `subnet.ScanSubnets` replaces `ScanSubnet(wg, chan)` and takes a context, so Ctrl-C stops a sweep.
- `subnet.ScanOptions` lives in `utils/utils.go` with the other shared types.
- Reverse DNS is bounded by a timeout and no longer holds a scan slot indefinitely.
- Linux reads `/proc/net/arp` directly, so `net-tools` is no longer required.
- Device lists are sorted by IP, so repeat scans produce identical JSON.
- Port scans now exit non-zero on invalid arguments instead of silently succeeding.
- `--ip` takes its website as a flag value, so later flags are no longer swallowed.

### Fixed
- `--subnet foo` panicked instead of returning an error.
- Networks wider than `/24` were silently truncated to a `/24` sweep.
- OS detection returned `Unknown` for every host one hop away from its default TTL.
- Live hosts were dropped when the ping output had no parseable TTL line.
- macOS ARP entries with single-digit octets (`0:1c:42:…`) were never matched.
- HTTP server versions were never parsed because the `Server:` header was compared against lowercase input.
- TCP banner grabbing opened a second connection per open port.
- Scan reports are written atomically and readable by the team (`0644`).

---

## [v0.3.1] - 2026-10-25

### Changed
- `--version` now reports the version and build date injected from the release tag at build time, instead of hardcoded values.

---

## [v0.3.0] - 2026-10-25

### Added
- `--ip <website>` resolves a domain name to its IP address.
- `--version` prints version, build date, author, and repository information.
- Running `--tcp` and `--udp` together performs both scans in one invocation.
- A unit test covering OS-specific TTL parsing in `PingWithTTL`.

### Fixed
- macOS `ping` timeouts are now passed in milliseconds (`-W 1000`), so reachable hosts are no longer reported as offline.
- TTL detection accepts `TTL=`, `TTL:`, and mixed casing across platforms.
- ARP table parsing handles Windows output and separators other than spaces.
- Local subnet detection skips networks narrower than /24 instead of scanning the wrong range.
- Subnet choices print their real prefix length instead of assuming a /24.

### Changed
- Linux and macOS share a single installer that detects the OS and architecture (x86_64, arm64) and downloads the latest release.
- The Windows installer detects architecture, installs from a temp directory, and updates only the user-scope `PATH`.
- Release archives now contain `netscanner` / `netscanner.exe`, matching what the installers expect.
- Removed `install/netscanner_install_mac.sh` in favour of the unified script.
- README updated for `--udp`, `--ip`, comma-separated subnets, the current module path, and new uninstall steps.

---

## [v0.2.0] - 2025-05-23

### Added
- Introduced CLI options for TCP and UDP port scanning.
- Added the ability for users to specify which subnet to scan.
- Implemented service identification for TCP and UDP ports.
- Enhanced concurrency for faster port scanning.
- Added banner grabbing for TCP ports to identify services and versions.

### Fixed
- Improved handling of invalid subnet inputs.

### Changed
- Redesigned the application to function as a full-featured CLI tool.
- Updated the help menu to reflect new CLI options.

---

## [v0.1.1] - 2025-05-22

### Added
- An option to rerun the scan again or quit the program.

### Fixed
- Resolved an issue where the binaries exit abruptly on scan completion.

### Changed
- N/A

---

## [v0.1.0] - 2025-05-22

- Initial public release of the network scanner tool.
- Supports subnet scanning with concurrency.
- Detects live devices by pinging IPs.
- Estimates OS based on TTL values.
- Retrieves hostnames via reverse DNS lookup.
- Gathers MAC addresses from ARP tables.
- Exports scan results to a formatted JSON file.
- Cross-platform support (Windows, Linux, macOS).

---