// Package cmd contains all the CLI commands for the Pork application.
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/abraira85/pork/internal/output"
	"github.com/abraira85/pork/internal/ports"
)

// freeCmd represents the "free" command.
// It checks if a specified port is free, and if it's busy, it automatically
// scans upwards for the next available port.
var freeCmd = &cobra.Command{
	Use:   "free [port]",
	Short: "Find the next free port",
	Long: `Checks if the specified port is currently available.
If the port is occupied, it scans upwards to find the next available port.
Useful when starting development servers and needing to quickly find an open port.`,
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		output.PrintBanner()

		port, err := parsePort(args[0])
		if err != nil {
			output.PrintError("%v", err)
			os.Exit(1)
		}

		// One scan covers the whole search: probing each candidate port
		// individually would re-enumerate every socket on the machine.
		scanner := ports.NewScanner()
		busy, err := scanner.GetPortMap()
		if err != nil {
			output.PrintError("Failed to check port: %v", err)
			os.Exit(1)
		}

		if len(busy[port]) == 0 {
			output.PrintSuccess("Port %d is free", port)
			return
		}

		output.PrintError("Port %d is busy", port)

		for candidate := port + 1; candidate <= ports.MaxPort; candidate++ {
			if len(busy[candidate]) == 0 {
				output.PrintInfo("Next free port: %d", candidate)
				return
			}
		}

		output.PrintWarning("No free port above %d", port)
		os.Exit(1)
	},
}

// init registers the freeCmd as a subcommand of rootCmd.
func init() {
	rootCmd.AddCommand(freeCmd)
}
