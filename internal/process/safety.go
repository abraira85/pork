// Package process provides cross-platform utilities for fetching process metadata,
// analyzing resource usage, and performing safety checks before terminating processes.
package process

import (
	"strings"
)

// IsCritical evaluates whether a given process name and PID belong to a critical
// system component. This is a safety measure to prevent the user from accidentally
// terminating vital services (like systemd, dockerd, or kernel tasks) which could
// cause system instability.
//
// It checks against well-known critical PIDs (0, 1) and a hardcoded list of
// daemon/service names common across Linux, macOS, and Windows.
// Returns true if the process is deemed critical and should not be killed.
func IsCritical(name string, pid int32) bool {
	// PID 1 is typically init/systemd on Linux, launchd on macOS.
	// PID 0 is the sched/kernel task.
	if pid == 1 || pid == 0 {
		return true
	}

	n := strings.ToLower(name)

	// Common critical daemons and system processes across multiple OS environments
	criticalNames := []string{
		// Linux/Unix
		"systemd", "init", "kernel_task", "dbus-daemon",
		"xorg", "wayland", "networkmanager", "sshd",
		"cron", "crond", "rsyslogd", "syslogd",
		"auditd", "kthreadd",

		// macOS
		"launchd",

		// Windows
		"svchost.exe", "csrss.exe", "wininit.exe",
		"smss.exe", "lsass.exe", "services.exe",
		"explorer.exe",

		// Container/Dev Infrastructure
		"dockerd", "containerd",
	}

	// Check if the process name contains any of the critical substrings
	for _, c := range criticalNames {
		if strings.Contains(n, c) {
			return true
		}
	}
	return false
}
