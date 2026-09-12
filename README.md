# Pork 🐷

A tiny terminal tool to inspect, visualize and free local ports.

Stop fighting with `lsof`, `ss`, `netstat` and random kill commands. Pork is a modern, cross-platform alternative that gives you clear, visual answers.

```bash
pork 3000
pork kill 3000
pork list
pork range 3000 3010
```

## Features

- **Inspect ports**: See exactly which process is using a port.
- **List active ports**: View all listening ports in a clean table.
- **Kill processes safely**: Free up ports with safety checks for critical system processes.
- **Find free ports**: Quickly scan for the next available port.
- **Visualize ranges**: See a map of busy and free ports in a given range.

## Installation

Since Pork is built with Go, you can easily install it using:

```bash
go install github.com/outboss/pork@latest
```

## Usage

### Inspect a Port

```bash
pork 3000
```
*Output:*
```text
🐷 Port 3000 is busy

PID       18422
Process   node
User      rober
Command   npm run dev
Address   127.0.0.1:3000
Protocol  TCP
Status    LISTEN
```

### List Active Ports

```bash
pork list
```
Displays a clean, formatted table of all ports currently in a `LISTEN` state.

### Free a Port (Kill Process)

```bash
pork kill 3000
```
Pork will find the process occupying the port and ask for your confirmation before safely terminating it. It prevents you from accidentally killing critical system processes.

### Find Next Free Port

```bash
pork free 3000
```
If port 3000 is busy, Pork will automatically scan upwards to find the next available port for you to use.

### Visualize a Port Range

```bash
pork range 3000 3010
```
Displays a visual map of the specified range, indicating which ports are free (`○`) and which are busy (`●`), along with the processes occupying the busy ones.

## Roadmap

- [x] **v0.1.0 (Current)**: Core CLI functionality (inspect, list, kill, free, range).
- [ ] **v0.2.0**: Advanced filtering (`pork list --json`, `--process node`), and flags (`--force`, `--yes`).
- [ ] **v0.3.0**: Real-time `watch` mode and interactive TUI (`pork ui`).
- [ ] **v1.0.0**: Stable cross-platform releases, automated binaries, and full documentation.

## License

MIT
