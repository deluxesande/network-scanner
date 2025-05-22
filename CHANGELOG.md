# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

- Added support for additional subnet selection options.
- Improved error handling for invalid IP ranges.
- Enhanced JSON export formatting.

---

## [v0.1.1] - 2025-05-22

### Added
- An Option to rerun the scan again or quit the program.

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