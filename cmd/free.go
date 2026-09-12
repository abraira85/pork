// Package cmd contains all the CLI commands for the Pork application.
package cmd

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/abraira85/pork/internal/output"
	"github.com/abraira85/pork/internal/ports"
)

// freeCmd represents the "free" command.
// It checks if a specified port is free, and if it's busy, it automatically
// scans for the next available port sequentially or within a specified range.
var freeCmd = &cobra.Command{
	Use:   "free [port]",
	Short: "Find the next free port",
	Long: `Checks if the specified port is currently available.
If the port is occupied, it scans upwards to find the next available port.
Useful when starting development servers and needing to quickly find an open port.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintBanner()
		portStr := args[0]
		portNum, err := strconv.ParseUint(portStr, 10, 32)
		if err != nil {
			output.PrintError("Invalid port number: %s", portStr)
			os.Exit(1)
		}

		scanner := ports.NewScanner()

		// Check the requested port
		isFree, err := scanner.IsFree(uint32(portNum))
		if err != nil {
			output.PrintError("Failed to check port: %v", err)
			os.Exit(1)
		}

		if isFree {
			output.PrintSuccess("Port %d is free", portNum)
			return
		}

		output.PrintError("Port %d is busy", portNum)

		// Find next free
		for i := portNum + 1; i < 65535; i++ {
			free, err := scanner.IsFree(uint32(i))
			if err != nil {
				continue
			}
			if free {
				output.PrintInfo("Next free port: %d", i)
				return
			}
		}
	},
}

// init registers the freeCmd as a subcommand of rootCmd.
func init() {
	rootCmd.AddCommand(freeCmd)
}
