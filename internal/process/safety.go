// Package process provides cross-platform utilities for fetching process metadata,
// analyzing resource usage, and performing safety checks before terminating processes.
package process

import (
	"strings"
)

// criticalNames lists processes whose termination is likely to destabilise the
// machine. Only daemons that plausibly hold a TCP port are worth listing here,
// since IsCritical is consulted exclusively for port owners.
var criticalNames = map[string]bool{
	// Linux/Unix
	"systemd":        true,
	"init":           true,
	"kernel_task":    true,
	"dbus-daemon":    true,
	"xorg":           true,
	"networkmanager": true,
	"sshd":           true,
	"cron":           true,
	"crond":          true,
	"rsyslogd":       true,
	"syslogd":        true,
	"auditd":         true,
	"kthreadd":       true,

	// macOS
	"launchd": true,

	// Windows (the .exe suffix is stripped before lookup)
	"svchost":  true,
	"csrss":    true,
	"wininit":  true,
	"smss":     true,
	"lsass":    true,
	"services": true,
	"explorer": true,

	// Container/Dev Infrastructure
	"dockerd":    true,
	"containerd": true,
}

// criticalPrefixes covers daemon families that spawn many differently-named
// helpers, such as systemd-resolved, systemd-networkd or systemd-logind.
var criticalPrefixes = []string{"systemd-"}

// IsCritical evaluates whether a given process name and PID belong to a critical
// system component. This is a safety measure so that the user is warned before
// terminating vital services (like systemd, dockerd, or kernel tasks) which could
// cause system instability.
//
// Matching is exact against the normalised name — lowercased, trimmed, and with
// any Windows .exe suffix removed — plus a short list of prefixes. Substring
// matching was tempting but flags ordinary user processes whose names merely
// contain a keyword, such as initdb, cronjob-runner or my-init.
func IsCritical(name string, pid int32) bool {
	// PID 1 is typically init/systemd on Linux, launchd on macOS.
	// PID 0 is the sched/kernel task. Anything lower is not a real process.
	if pid <= 1 {
		return true
	}

	n := strings.TrimSuffix(strings.ToLower(strings.TrimSpace(name)), ".exe")
	if n == "" {
		return false
	}

	if criticalNames[n] {
		return true
	}

	for _, prefix := range criticalPrefixes {
		if strings.HasPrefix(n, prefix) {
			return true
		}
	}

	return false
}
