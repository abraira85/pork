// Package process provides cross-platform utilities for fetching process metadata,
// analyzing resource usage, and performing safety checks before terminating processes.
package process

import (
	"fmt"
	"os"
)

// Terminate gracefully attempts to kill the process with the given PID.
// It relies on standard os.Process behavior.
// On Unix-like systems, os.Process.Kill() sends SIGKILL.
// In the future, this can be expanded with build tags to send SIGTERM first.
func Terminate(pid int32) error {
	p, err := os.FindProcess(int(pid))
	if err != nil {
		return fmt.Errorf("could not find process: %w", err)
	}

	err = p.Kill()
	if err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	return nil
}
