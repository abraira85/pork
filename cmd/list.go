// Package cmd contains all the CLI commands for the Pork application.
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/abraira85/pork/internal/output"
	"github.com/abraira85/pork/internal/ports"
)

// listCmd represents the "list" command.
// It retrieves all active local ports that are currently in a LISTEN state,
// along with the process information using those ports, and displays them in a formatted table.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List active local ports",
	Long: `Lists all active local ports currently listening on the machine.
It displays a clean table with port, PID, and program details.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintBanner()

		scanner := ports.NewScanner()
		activePorts, err := scanner.GetActivePorts()
		if err != nil {
			output.PrintError("Failed to scan active ports: %v", err)
			os.Exit(1)
		}

		output.PrintActivePortsTable(activePorts)
	},
}

// init registers the listCmd as a subcommand of rootCmd.
func init() {
	rootCmd.AddCommand(listCmd)
}
