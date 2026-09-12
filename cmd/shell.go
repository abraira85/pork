package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/outboss/pork/internal/output"
	"github.com/outboss/pork/internal/ports"
	"github.com/outboss/pork/internal/tui"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Launch the interactive TUI",
	Long:  `Launches an interactive terminal user interface to browse, inspect, and kill processes occupying local ports.`,
	Run: func(cmd *cobra.Command, args []string) {
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
