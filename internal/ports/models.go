// Package ports provides data structures and logic for scanning local network ports.
// It handles mapping network connections to operating system processes.
package ports

// PortInfo encapsulates all relevant details about a specific network port
// and the underlying OS process that is currently using it.
type PortInfo struct {
	// Port is the numerical port number (e.g., 3000, 8080).
	Port uint32
	// Protocol is the networking protocol used (currently always "TCP").
	Protocol string
	// Address is the local IP address the port is bound to (e.g., "127.0.0.1", "0.0.0.0").
	Address string
	// Status represents the connection state (e.g., "LISTEN").
	Status string
	// PID is the Process ID of the process holding the port.
	PID int32
	// Process is the short name of the executing program (e.g., "node", "postgres").
	Process string
	// Command is the full command-line invocation used to start the process.
	Command string
	// User is the system username that owns the process.
	User string
}
