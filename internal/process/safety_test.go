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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCritical(tt.procName, tt.pid)
			if result != tt.expected {
				t.Errorf("IsCritical(%s, %d) = %v, expected %v", tt.procName, tt.pid, result, tt.expected)
			}
		})
	}
}

func TestFormatCommand(t *testing.T) {
	tests := []struct {
		name      string
		cmd       string
		maxLength int
		expected  string
	}{
		{"Short command", "npm start", 40, "npm start"},
		{"Exact length command", "1234567890", 10, "1234567890"},
		{"Long command truncated", "npm run dev --host 0.0.0.0 --port 3000", 20, "npm run dev --hos..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCommand(tt.cmd, tt.maxLength)
			if result != tt.expected {
				t.Errorf("FormatCommand(%s, %d) = %s, expected %s", tt.cmd, tt.maxLength, result, tt.expected)
			}
		})
	}
}
