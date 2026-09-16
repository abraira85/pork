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
- `pork kill --yes` to skip the confirmation prompt, and `--force` to allow
  terminating a critical process non-interactively.
- `pork kill` now requests a clean shutdown (SIGTERM) and only escalates to a
  forced kill (SIGKILL) if the process is still alive after a grace period.
- Markdown linting and a `go mod tidy` check wired into the Makefile
  (`make lint-docs`, `make tidy`) and CI, which already had a markdownlint
  config that nothing ever ran.
- Tests for port parsing, process identification, and the TUI state machine.

### Changed

- Port scans are performed once per command instead of once per probed port.
  `pork free` and `pork range` used to re-enumerate every socket on the machine
  for each candidate; `pork range 3000 3300` went from ~5.9s to ~0.02s.
- `IsCritical` matches process names exactly (case-insensitive, `.exe` stripped)
  plus a `systemd-` prefix rule, instead of substring matching that flagged
  ordinary processes such as `initdb`, `cronjob-runner` or `my-init`.
- The TUI asks for confirmation before killing a critical process instead of
  refusing outright, matching what `pork kill` does.
- The TUI terminates processes in a background command, so the grace period no
  longer freezes the interface, and it reuses the same kill path as the CLI.
- `pork range` reports unreadable processes as `Unknown (sudo)` rather than a
  blank name and `PID 0`, matching `pork list`.

### Removed

- Roadmap section in the README. Planned work is tracked in GitHub issues and
  milestones instead, so the README no longer drifts out of date between releases.
- Unused `PortScanResult` type, `PortInfo.String()`, and `FormatCommand`, which
  only its own test referenced.
- Per-process CPU, memory and executable-path collection. Nothing rendered these
  fields, but they were the most expensive part of every scan.

### Fixed

- `pork <port>` failed with `unknown command` — the root command rejected positional
  arguments because it also has subcommands. Now `pork 3000` inspects the port correctly.
- `pork free` never considered port 65535 (off-by-one). It now scans up to and includes 65535.
- Migrated the linter to golangci-lint v2 and fixed all reported issues (dead code,
  unchecked errors, unused parameters, and staticcheck/style findings).
- Release pipeline: GoReleaser now runs inside the semantic-release workflow so that
  GitHub Releases are published even when the tag is pushed by GITHUB_TOKEN.
- Release pipeline: the created tag is fetched explicitly before publishing, and
  `GORELEASER_CURRENT_TAG` pins the version GoReleaser reports.
- Port arguments are validated against the 1-65535 range. `pork 99999` used to
  report an impossible port as free, and `pork range 1 4000000000` would spin for
  billions of iterations.
- `pork free` now reports when there is no free port above the one requested,
  instead of exiting silently with a success status.
- `pork kill` documented a `--yes` flag that did not exist, and described SIGKILL
  as a graceful termination.
- `pork kill` frees every process bound to the port, not just the first one.
  Listeners are now deduplicated by port *and* PID, so two processes sharing a
  port on different interfaces are both reported.
- `pork <port> <extra>` no longer silently ignores the extra argument, and a
  mistyped subcommand reports an unknown command instead of an invalid port.
- The TUI help line advertised `Esc: Quit` in the table view, where Esc did
  nothing. `q` no longer quits from inside a menu or confirmation.
- The TUI no longer blanks the table when a refresh scan fails; it keeps the last
  good snapshot.
- Table columns are recalculated on terminal resize, not just at startup.
- Process names are truncated by rune instead of by byte, so multi-byte names are
  never cut mid-character.
- `go.mod` listed `bubbles` and `bubbletea` as indirect dependencies although both
  are imported directly.
- Every tracked text file now ends with a newline, as `.editorconfig` requires.

## [0.1.0] - 2026-06-27

### Added

- `pork <port>` — inspect a port and show the process occupying it.
- `pork list` — table of all listening ports.
- `pork kill <port>` — safely free a port with confirmation and critical-process protection.
- `pork free <port>` — find the next available port.
- `pork range <start> <end>` — visual map of busy/free ports.
- `pork shell` — interactive TUI to browse, inspect, and kill.
