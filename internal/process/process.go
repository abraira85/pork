// Package process provides cross-platform utilities for fetching process metadata,
// analyzing resource usage, and performing safety checks before terminating processes.
package process

import (
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// Info encapsulates the metadata and resource usage of an operating system process.
// It is fetched on-demand to avoid heavy system calls when not needed.
type Info struct {
	// PID is the unique Process ID.
	PID int32
	// Name is the short executable name (e.g., "nginx").
	Name string
	// Command is the full command line used to launch the process.
	Command string
	// User is the operating system user who owns the process.
	User string
	// Path is the absolute file path to the executable binary.
	Path string
	// CPU represents the percentage of CPU currently utilized by the process.
	CPU float64
	// MemoryMB represents the Resident Set Size (RSS) memory used in Megabytes.
	MemoryMB float64
}

// GetInfo retrieves comprehensive metadata for a given Process ID (PID).
// It queries the operating system (using gopsutil) to gather details like
// the process name, command line arguments, owner, and resource usage.
// Returns an error if the process does not exist or if permissions are insufficient.
func GetInfo(pid int32) (*Info, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}

	name, _ := p.Name()
	cmd, _ := p.Cmdline()
	user, _ := p.Username()
	path, _ := p.Exe()

	// CPU usage calculation can block or be inaccurate without a time window.
	// For this tool, a quick percentage query is sufficient.
	cpu, _ := p.CPUPercent()

	memInfo, err := p.MemoryInfo()
	var memMB float64
	if err == nil && memInfo != nil {
		memMB = float64(memInfo.RSS) / 1024 / 1024
	}

	return &Info{
		PID:      pid,
		Name:     name,
		Command:  cmd,
		User:     user,
		Path:     path,
		CPU:      cpu,
		MemoryMB: memMB,
	}, nil
}

// Identify returns a simplified, human-readable identifier for the process.
func Identify(name, cmd string) string {
	if cmd != "" {
		parts := strings.Split(cmd, " ")
		if len(parts) > 0 {
			base := filepath.Base(parts[0])
			if base != "" && base != "." {
				// Truncate if it's absurdly long
				if len(base) > 30 {
					return base[:27] + "..."
				}
				return base
			}
		}
	}
	if len(name) > 30 {
		return name[:27] + "..."
	}
	return name
}

// FormatCommand truncates a long command line string to a specified maximum length,
// appending an ellipsis ("...") if truncation occurred. This is useful for UI display
// where long commands might break table formatting.
func FormatCommand(cmd string, maxLength int) string {
	if len(cmd) > maxLength && maxLength > 3 {
		return cmd[:maxLength-3] + "..."
	}
	return cmd
}
