# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

- N/A

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