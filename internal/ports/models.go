// Package ports provides data structures and logic for scanning local network ports.
// It handles mapping network connections to operating system processes.
package ports

import (
	"fmt"
)

// PortInfo encapsulates all relevant details about a specific network port
// and the underlying OS process that is currently using it.
type PortInfo struct {
	// Port is the numerical port number (e.g., 3000, 8080).
	Port uint32
	// Protocol is the networking protocol used (e.g., "TCP", "UDP").
	Protocol string
	// Address is the local IP address the port is bound to (e.g., "127.0.0.1", "0.0.0.0").
	Address string
	// Status represents the connection state (e.g., "LISTEN", "ESTABLISHED").
	Status string
	// PID is the Process ID of the process holding the port.
	PID int32
	// Process is the short name of the executing program (e.g., "node", "postgres").
	Process string
	// Command is the full command-line invocation used to start the process.
	Command string
	// User is the system username that owns the process.
	User string
	// Path is the absolute file path to the process executable.
	Path string
	// CPU is the current CPU utilization percentage of the process.
	CPU float64
	// MemoryMB is the current memory (RSS) usage of the process in Megabytes.
	MemoryMB float64
}

// String returns a human-readable string representation of the PortInfo.
func (p *PortInfo) String() string {
	return fmt.Sprintf("%d (%s) - %s (PID: %d)", p.Port, p.Protocol, p.Process, p.PID)
}

// PortScanResult represents the outcome of querying a specific port.
// It indicates whether the port is currently occupied, and if so,
// provides the associated PortInfo.
type PortScanResult struct {
	// Port is the numerical port that was scanned.
	Port int
	// Busy is true if the port is currently bound/listening, false if it's free.
	Busy bool
	// Info contains the process details if Busy is true. It is nil if the port is free.
	Info *PortInfo
}
