# Contributing to Network Scanner

Thank you for your interest in contributing! This document explains how you can help improve this project.

## How to Contribute

### Reporting Issues
- Please use GitHub Issues to report bugs, request features, or ask questions.
- Provide as much detail as possible, including your OS, Go version, and steps to reproduce.

### Code Contributions
- Fork the repository and create a feature branch (`git checkout -b feature/your-feature`).
- Follow Go idioms and formatting conventions (`go fmt` your code).
- Write clear, concise commit messages.
- Test your changes thoroughly on supported platforms (Windows, Linux, macOS).
- Submit a pull request (PR) with a clear description of your changes.

### Code Style & Standards
- Use idiomatic Go code.
- Keep functions small and focused.
- Document exported functions and types using GoDoc comments.
- Handle errors properly; avoid ignoring errors.

### Testing
- Add tests for new features or bug fixes where applicable.
- Run all tests before submitting your PR:
  ```bash
  go test ./...
