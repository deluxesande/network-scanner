# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

- N/A

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