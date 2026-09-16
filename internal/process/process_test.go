package process

import (
	"strings"
	"testing"
)

func TestIdentify(t *testing.T) {
	tests := []struct {
		name     string
		procName string
		cmd      string
		expected string
	}{
		{"Prefers the command's executable", "node", "node /srv/app/server.js", "node"},
		{"Strips the executable's directory", "python3", "/usr/bin/python3 manage.py", "python3"},
		{"Falls back to the process name", "postgres", "", "postgres"},
		{"Ignores a whitespace-only command", "redis-server", "   ", "redis-server"},
		{"Collapses repeated spaces", "go", "go  run  ./main.go", "go"},
		{
			"Truncates an absurdly long executable",
			"x",
			strings.Repeat("a", 50),
			strings.Repeat("a", 27) + "...",
		},
		{
			"Truncates an absurdly long process name",
			strings.Repeat("b", 50),
			"",
			strings.Repeat("b", 27) + "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Identify(tt.procName, tt.cmd)
			if result != tt.expected {
				t.Errorf("Identify(%q, %q) = %q, expected %q", tt.procName, tt.cmd, result, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"Short string is untouched", "npm start", 40, "npm start"},
		{"Exact length is untouched", "1234567890", 10, "1234567890"},
		{"Long string is elided", "npm run dev --host 0.0.0.0", 20, "npm run dev --hos..."},
		{"Tiny budget has no room for an ellipsis", "abcdef", 3, "abc"},
		{"Zero budget yields an empty string", "abcdef", 0, ""},
		// Cutting by bytes would slice these multi-byte runes in half.
		{"Multi-byte runes are not split", "ñññññññ", 5, "ññ..."},
		{"Multi-byte runes fit by count, not bytes", "ñññññ", 5, "ñññññ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, expected %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}
