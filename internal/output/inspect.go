// Package output provides terminal styling and formatting utilities for Pork.
package output

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/abraira85/pork/internal/ports"
)

// PrintPortInspection renders the details of a specific port in a clear, formatted layout.
// If the port is free, it prints a success message.
// If the port is busy, it lists the details of the process holding it.
func PrintPortInspection(port uint32, info *ports.PortInfo) {
	if info == nil {
		PrintSuccess("Port %d is free", port)
		return
	}

	PrintError("Port %d is busy\n", port)

	// Format fields
	fields := []struct {
		Label string
		Value string
	}{
		{"PID", fmt.Sprintf("%d", info.PID)},
		{"Process", info.Process},
		{"User", info.User},
		{"Command", info.Command},
		{"Address", info.Address},
		{"Protocol", info.Protocol},
		{"Status", info.Status},
	}

	var builder strings.Builder
	for _, f := range fields {
		if f.Value != "" {
			builder.WriteString(LabelStyle.Render(f.Label) + ValueStyle.Render(f.Value) + "\n")
		}
	}

	fmt.Println(lipgloss.NewStyle().PaddingLeft(2).Render(builder.String()))
}
