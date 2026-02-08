# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.1] - 2024-05-24

### Added
- **Cross-Platform Support**: Unified interface (registry pattern) supporting both Dell iDRAC and Lenovo XClarity BMCs.
- **Power Management**: New `power` command to control server power state (`on`, `off`, `restart`, `status`).
- **System Information**: New `sysinfo` command to retrieve detailed hardware info (Model, Serial, BIOS version, SKU).
- **Boot Management**: New `boot` command to view boot order and set the next boot device (e.g., PXE, HDD).
- **Event Logs**: New `eventlog` command to inspect System Event Logs (SEL).
- **CI/CD**: Automated release pipeline using GitHub Actions and GoReleaser for macOS, Linux, and Windows builds.
- **Mock Client**: Internal mock client implementation for improved unit testing capabilities.

### Changed
- **Architecture**: Refactored `health` and `controller` commands to use the new plugin registry, removing hardcoded vendor logic.
- **Testing**: Enhanced test coverage for new commands and updated existing tests to use the new architecture.
- **Dependencies**: Upgraded project to Go 1.25 and updated all dependencies.

### Fixed
- **HTTP Client**: Improved HTTP client configuration to transparently handle JSON payloads and content types.
