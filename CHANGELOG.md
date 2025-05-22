# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

- Initial public release of the network scanner tool.
- Supports subnet scanning with concurrency.
- Detects live devices by pinging IPs.
- Estimates OS based on TTL values.
- Retrieves hostnames via reverse DNS lookup.
- Gathers MAC addresses from ARP tables.
- Exports scan results to a formatted JSON file.
- Cross-platform support (Windows, Linux, macOS).

---

## [0.1.0] - 2025-05-22

### Added
- Core scanning functionality for IPv4 /24 subnets.
- TTL-based OS fingerprinting.
- Hostname resolution and MAC address retrieval.
- User prompt for subnet selection.
- JSON export of scan results.
- Terminal output formatting with colors for better readability.

### Fixed
- N/A (initial release)

### Deprecated
- N/A

### Removed
- N/A
