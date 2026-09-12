// Package output provides terminal styling and formatting utilities for Pork.
// It leverages lipgloss for consistent and beautiful console output.
package output

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	// PrimaryColor is the main brand color for Pork (a vibrant pink/magenta).
	PrimaryColor = lipgloss.Color("#FF2A75")
	
	// SecondaryColor is used for accents.
	SecondaryColor = lipgloss.Color("#FF7EB3")
	
	// SuccessColor is used for positive outcomes (e.g., port is free).
	SuccessColor = lipgloss.Color("#00E676")
	
	// WarningColor is used for cautions (e.g., critical process).
	WarningColor = lipgloss.Color("#FFC400")
	
	// DangerColor is used for destructive actions (e.g., killing a process) or busy ports.
	DangerColor = lipgloss.Color("#FF1744")
	
	// MutedColor is used for secondary information.
	MutedColor = lipgloss.Color("#757575")

	// Pig is the standard Pork mascot emoji prefix.
	Pig = lipgloss.NewStyle().Foreground(PrimaryColor).Render("🐷")
	
	// WarningIcon is used for warnings.
	WarningIcon = lipgloss.NewStyle().Foreground(WarningColor).Render("⚠️")

	// TitleStyle is used for main section headers.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginBottom(1)
			
	// LabelStyle is used for field names (e.g., "PID", "Process").
	LabelStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Width(10)
			
	// ValueStyle is used for field values.
	ValueStyle = lipgloss.NewStyle().
			Bold(true)
)

// PrintSuccess prints a success message with the Pork mascot.
func PrintSuccess(msg string, args ...interface{}) {
	formatted := fmt.Sprintf(msg, args...)
	fmt.Printf("%s %s\n", Pig, lipgloss.NewStyle().Foreground(SuccessColor).Render(formatted))
}

// PrintInfo prints an informational message with the Pork mascot.
func PrintInfo(msg string, args ...interface{}) {
	formatted := fmt.Sprintf(msg, args...)
	fmt.Printf("%s %s\n", Pig, formatted)
}

// PrintWarning prints a warning message.
func PrintWarning(msg string, args ...interface{}) {
	formatted := fmt.Sprintf(msg, args...)
	fmt.Printf("%s %s\n", WarningIcon, lipgloss.NewStyle().Foreground(WarningColor).Render(formatted))
}

// PrintError prints an error message.
func PrintError(msg string, args ...interface{}) {
	formatted := fmt.Sprintf(msg, args...)
	fmt.Printf("%s\n", lipgloss.NewStyle().Foreground(DangerColor).Bold(true).Render("Error: "+formatted))
}
