package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/abraira85/pork/internal/output"
	"github.com/abraira85/pork/internal/ports"
	"github.com/abraira85/pork/internal/tui"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Launch the interactive TUI",
	Long:  `Launches an interactive terminal user interface to browse, inspect, and kill processes occupying local ports.`,
	Run: func(_ *cobra.Command, _ []string) {
		scanner := ports.NewScanner()
		activePorts, err := scanner.GetActivePorts()
		if err != nil {
			output.PrintError("Failed to scan ports: %v", err)
			os.Exit(1)
		}

		m := tui.NewModel(activePorts)
		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
