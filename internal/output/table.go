// Package output provides terminal styling and formatting utilities for Pork.
package output

import (
	"fmt"
	"os"

	"golang.org/x/term"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/outboss/pork/internal/ports"
	"github.com/outboss/pork/internal/process"
)

// getTerminalWidth returns the width of the terminal, or a default if it cannot be determined.
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		return 80 // fallback width
	}
	return width
}

// PrintActivePortsTable renders a beautiful terminal table displaying the given active ports.
func PrintActivePortsTable(activePorts []*ports.PortInfo) {
	if len(activePorts) == 0 {
		PrintInfo("No active local ports found or all filtered.")
		return
	}

	termWidth := getTerminalWidth()
	
	progWidth := termWidth - 25
	if progWidth < 15 {
		progWidth = 15
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
		programStr := "-"
		
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
