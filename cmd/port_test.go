package cmd

import "testing"

func TestParsePort(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		expected uint32
		wantErr  bool
	}{
		{"Typical dev port", "3000", 3000, false},
		{"Lowest valid port", "1", 1, false},
		{"Highest valid port", "65535", 65535, false},

		// Port 0 means "let the kernel choose" and can never be inspected.
		{"Zero is rejected", "0", 0, true},

		// ParseUint with a 32-bit size used to let these through, which made
		// pork report imaginary ports as free and range scans run forever.
		{"Above the port space is rejected", "65536", 0, true},
		{"Far above the port space is rejected", "4000000000", 0, true},

		{"Non-numeric is rejected", "lst", 0, true},
		{"Negative is rejected", "-1", 0, true},
		{"Empty is rejected", "", 0, true},
		{"Whitespace is rejected", " 3000", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePort(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parsePort(%q) error = %v, wantErr %v", tt.arg, err, tt.wantErr)
			}
			if got != tt.expected {
				t.Errorf("parsePort(%q) = %d, expected %d", tt.arg, got, tt.expected)
			}
		})
	}
}
