// Package output provides terminal styling and formatting utilities for Pork.
package output

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/abraira85/pork/internal/ports"
)

// PrintPortInspection renders the details of a specific port in a clear, formatted layout.
// If the port is free, it prints a success message. If more than one process is
// bound to the port on different interfaces, every listener is shown.
func PrintPortInspection(port uint32, listeners []*ports.PortInfo) {
	if len(listeners) == 0 {
		PrintSuccess("Port %d is free", port)
		return
	}

	PrintError("Port %d is busy\n", port)

	if len(listeners) > 1 {
		PrintWarning("%d processes are listening on port %d\n", len(listeners), port)
	}

	var builder strings.Builder
	for i, info := range listeners {
		if i > 0 {
			builder.WriteString("\n")
		}
		writePortFields(&builder, info)
	}

	fmt.Println(lipgloss.NewStyle().PaddingLeft(2).Render(builder.String()))
}

// writePortFields renders the non-empty fields of a single listener.
func writePortFields(builder *strings.Builder, info *ports.PortInfo) {
	// Without elevated privileges the listening socket is visible but its
	// owning process is not, which is worth saying rather than printing "0".
	pid := "Unknown (requires sudo)"
	if info.PID > 0 {
		pid = fmt.Sprintf("%d", info.PID)
	}

	fields := []struct {
		Label string
		Value string
	}{
		{"PID", pid},
		{"Process", info.Process},
		{"User", info.User},
		{"Command", info.Command},
		{"Address", info.Address},
		{"Protocol", info.Protocol},
		{"Status", info.Status},
	}

	for _, f := range fields {
		if f.Value != "" {
			builder.WriteString(LabelStyle.Render(f.Label) + ValueStyle.Render(f.Value) + "\n")
		}
	}
}
