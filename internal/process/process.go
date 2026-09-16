// Package process provides cross-platform utilities for fetching process metadata,
// analyzing resource usage, and performing safety checks before terminating processes.
package process

import (
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// maxDisplayLen is the widest a process identifier may be before it is elided.
const maxDisplayLen = 30

// Info encapsulates the metadata of an operating system process.
//
// It deliberately carries only the fields Pork actually renders. Resource
// counters such as CPU% or RSS are far more expensive to collect than the
// basics, and every extra field is paid once per listening socket on every
// scan, including the TUI's periodic refresh.
type Info struct {
	// PID is the unique Process ID.
	PID int32
	// Name is the short executable name (e.g., "nginx").
	Name string
	// Command is the full command line used to launch the process.
	Command string
	// User is the operating system user who owns the process.
	User string
}

// GetInfo retrieves metadata for a given Process ID (PID).
// It queries the operating system (using gopsutil) to gather the process name,
// command line arguments, and owner.
// Returns an error if the process does not exist or if permissions are insufficient.
func GetInfo(pid int32) (*Info, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}

	// Each accessor may fail independently (most often because the process is
	// owned by another user); an empty field is better than discarding the rest.
	name, _ := p.Name()
	cmd, _ := p.Cmdline()
	user, _ := p.Username()

	return &Info{
		PID:     pid,
		Name:    name,
		Command: cmd,
		User:    user,
	}, nil
}

// Identify returns a simplified, human-readable identifier for the process.
// It prefers the executable invoked on the command line and falls back to the
// process name reported by the OS.
func Identify(name, cmd string) string {
	if cmd != "" {
		parts := strings.Fields(cmd)
		if len(parts) > 0 {
			base := filepath.Base(parts[0])
			if base != "" && base != "." && base != string(filepath.Separator) {
				return truncate(base, maxDisplayLen)
			}
		}
	}
	return truncate(name, maxDisplayLen)
}

// truncate shortens s to at most maxLen characters, appending an ellipsis when
// it has to cut. It counts runes rather than bytes so that a multi-byte path is
// never sliced mid-character.
func truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}
