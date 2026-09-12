// Package output provides terminal styling and formatting utilities for Pork.
package output

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/abraira85/pork/internal/ports"
	"github.com/abraira85/pork/internal/process"
)

// PrintActivePortsTable renders a beautiful terminal table displaying the given active ports.
func PrintActivePortsTable(activePorts []*ports.PortInfo) {
	if len(activePorts) == 0 {
		PrintInfo("No active local ports found or all filtered.")
		return
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(PrimaryColor)).
		Headers("PORT", "PID", "PROGRAM").
		StyleFunc(func(row, col int) lipgloss.Style {
			// Row 0 is the header in lipgloss table
			if row == 0 {
				return lipgloss.NewStyle().
					Bold(true).
					Foreground(SecondaryColor).
					Padding(0, 1)
			}

			// Row style for data (row >= 1)
			style := lipgloss.NewStyle().Padding(0, 1)

			if col == 0 {
				style = style.Foreground(PrimaryColor).Bold(true)
			}

			// Safe bounds check for activePorts
			dataIdx := row - 1
			if col == 1 && dataIdx >= 0 && dataIdx < len(activePorts) && activePorts[dataIdx].PID > 0 {
				style = style.Foreground(MutedColor)
			}
			return style
		})

	for _, p := range activePorts {
		pidStr := "-"
		var programStr string

		if p.PID > 0 {
			pidStr = fmt.Sprintf("%d", p.PID)
			programStr = process.Identify(p.Process, p.Command)
		} else {
			programStr = "Unknown (requires sudo)"
		}

		t.Row(fmt.Sprintf("%d", p.Port), pidStr, programStr)
	}

	fmt.Println(t.Render())
}
