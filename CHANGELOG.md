# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic
Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Allow showing commit logs between versions

### Changed

- The command line interface for running checks

## [v1.4.1] - Mar 01, 2025

### Fixed

- Logic for determining if versions are different in scenarios where one or more
  envs are missing

### Changed

- Add an upper bound to maximum number of concurrent fetches

## [v1.4.0] - Feb 14, 2025

### Added

- Show task definition registration time for each version
- Allow configuring html title via `--html-title`

### Changed

- Go and dependency upgrades

## [v1.3.1] - Jan 07, 2024

### Changed

- Go and dependency upgrades

## [v1.3.0] - Aug 20, 2024

### Added

- Add tabular output using `-f table`

### Changed

- Config file default location on darwin changed `~/Library/Application
  Support/ecsv/ecsv.yml`

[unreleased]: https://github.com/dhth/ecsv/compare/v1.4.1...HEAD
[v1.4.1]: https://github.com/dhth/ecsv/compare/v1.4.0...v1.4.1
[v1.4.0]: https://github.com/dhth/ecsv/compare/v1.3.1...v1.4.0
[v1.3.1]: https://github.com/dhth/ecsv/compare/v1.3.0...v1.3.1
[v1.3.0]: https://github.com/dhth/ecsv/compare/v1.2.2...v1.3.0
