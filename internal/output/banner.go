// Package output provides terminal styling and formatting utilities for Pork.
package output

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

// PrintBanner displays a styled ASCII art banner for the tool.
func PrintBanner() {
	banner := `
  ____             _    
 |  _ \ ___  _ __ | | __
 | |_) / _ \| '__|| |/ /
 |  __/ (_) | |   |   < 
 |_|   \___/|_|   |_|\_\
`
	
	style := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true)
		
	subtitle := lipgloss.NewStyle().
		Foreground(MutedColor).
		Italic(true).
		Render("  The modern port inspector")

	fmt.Println(style.Render(banner))
	fmt.Println(subtitle)
	fmt.Println()
}
