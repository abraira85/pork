# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Project documentation and open-source structure (CONTRIBUTING, CODE_OF_CONDUCT, SECURITY).
- CI pipeline (lint, vet, test, build) running on a self-hosted runner.
- Automated releases with Semantic Versioning: every push to `main` bumps the
  version tag from Conventional Commits (`patch` by default, `feat!`/`feat:`
  bump `major`/`minor`) and GoReleaser publishes tagged binaries.
- Version injected into the binary (`pork --version`).

### Fixed

- `pork <port>` failed with `unknown command` — the root command rejected positional
  arguments because it also has subcommands. Now `pork 3000` inspects the port correctly.
- `pork free` never considered port 65535 (off-by-one). It now scans up to and includes 65535.
- Migrated the linter to golangci-lint v2 and fixed all reported issues (dead code,
  unchecked errors, unused parameters, and staticcheck/style findings).
- Release pipeline: GoReleaser now runs inside the semantic-release workflow so that
  GitHub Releases are published even when the tag is pushed by GITHUB_TOKEN.

## [0.1.0] - 2026-06-27

### Added

- `pork <port>` — inspect a port and show the process occupying it.
- `pork list` — table of all listening ports.
- `pork kill <port>` — safely free a port with confirmation and critical-process protection.
- `pork free <port>` — find the next available port.
- `pork range <start> <end>` — visual map of busy/free ports.
- `pork shell` — interactive TUI to browse, inspect, and kill.