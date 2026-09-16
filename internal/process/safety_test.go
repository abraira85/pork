package process

import (
	"testing"
)

func TestIsCritical(t *testing.T) {
	tests := []struct {
		name     string
		procName string
		pid      int32
		expected bool
	}{
		{"PID 1 is critical", "init", 1, true},
		{"PID 0 is critical", "sched", 0, true},
		{"Systemd is critical", "systemd", 123, true},
		{"Launchd is critical", "launchd", 456, true},
		{"Normal process is not critical", "node", 18422, false},
		{"Normal process uppercase is not critical", "NODE", 18422, false},
		{"Dockerd is critical", "dockerd", 999, true},

		// Windows names are matched with the .exe suffix stripped.
		{"lsass.exe is critical", "lsass.exe", 640, true},
		{"LSASS.EXE is critical", "LSASS.EXE", 640, true},

		// systemd spawns many differently-named helpers.
		{"systemd-resolved is critical", "systemd-resolved", 771, true},
		{"systemd-networkd is critical", "systemd-networkd", 772, true},

		// Names that merely contain a critical keyword must not be flagged:
		// this is what substring matching used to get wrong.
		{"initdb is not critical", "initdb", 5001, false},
		{"cronjob-runner is not critical", "cronjob-runner", 5002, false},
		{"my-init is not critical", "my-init", 5003, false},
		{"sshd-tunnel-helper is not critical", "sshd-tunnel-helper", 5004, false},
		{"containerd-shim-mine is not critical", "containerd-shim-mine", 5005, false},

		// Surrounding whitespace should not change the verdict.
		{"Padded name is normalised", "  sshd  ", 600, true},
		{"Empty name is not critical", "", 7000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCritical(tt.procName, tt.pid)
			if result != tt.expected {
				t.Errorf("IsCritical(%q, %d) = %v, expected %v", tt.procName, tt.pid, result, tt.expected)
			}
		})
	}
}
