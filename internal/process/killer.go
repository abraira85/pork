// Package process provides cross-platform utilities for fetching process metadata,
// analyzing resource usage, and performing safety checks before terminating processes.
package process

import (
	"fmt"
	"os"
	"time"

	gopsprocess "github.com/shirou/gopsutil/v3/process"
)

const (
	// gracePeriod is how long a process is given to shut down cleanly after
	// being asked to, before Pork forces the issue.
	gracePeriod = 3 * time.Second

	// pollInterval is how often the process is re-checked while waiting.
	pollInterval = 50 * time.Millisecond
)

// Terminate stops the process with the given PID.
//
// It first asks the process to shut down cleanly (SIGTERM on Unix), giving it
// gracePeriod to close sockets and flush state, and only escalates to a forced
// kill (SIGKILL) if it is still alive afterwards. Windows has no deliverable
// equivalent of SIGTERM, so there the forced kill happens immediately.
//
// A process that disappears on its own while we are waiting is treated as
// success: the caller's goal is a free port, not a specific signal.
func Terminate(pid int32) error {
	p, err := os.FindProcess(int(pid))
	if err != nil {
		return fmt.Errorf("could not find process %d: %w", pid, err)
	}

	graceful, err := requestShutdown(p)
	if err != nil {
		// The process may have exited between the port scan and the signal.
		if !isRunning(pid) {
			return nil
		}
	} else if graceful && waitForExit(pid, gracePeriod) {
		return nil
	}

	if err := p.Kill(); err != nil {
		if !isRunning(pid) {
			return nil
		}
		return fmt.Errorf("failed to kill process %d: %w", pid, err)
	}

	return nil
}

// waitForExit polls until the process is gone or the timeout expires.
// It reports whether the process actually exited.
func waitForExit(pid int32, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !isRunning(pid) {
			return true
		}
		time.Sleep(pollInterval)
	}
	return !isRunning(pid)
}

// isRunning reports whether a process with the given PID still exists.
func isRunning(pid int32) bool {
	exists, err := gopsprocess.PidExists(pid)
	return err == nil && exists
}
