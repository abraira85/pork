// Package cmd contains all the CLI commands for the Pork application.
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/abraira85/pork/internal/output"
	"github.com/abraira85/pork/internal/ports"
	"github.com/abraira85/pork/internal/process"
)

// rangeCmd represents the "range" command.
// It scans a specified range of ports and visualizes their status (free or busy)
// in a clean, terminal-friendly map.
var rangeCmd = &cobra.Command{
	Use:   "range [start] [end]",
	Short: "Visual map of a port range",
	Long: `Scans an inclusive range of ports (from start to end) and displays a visual map.
It shows which ports are free and which are busy, along with the process information
for the busy ports.`,
	Args: cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		output.PrintBanner()

		start, errStart := parsePort(args[0])
		end, errEnd := parsePort(args[1])

		if errStart != nil || errEnd != nil || start > end {
			output.PrintError("Invalid range: %s to %s (ports must be %d-%d, and start must not exceed end)",
				args[0], args[1], ports.MinPort, ports.MaxPort)
			os.Exit(1)
		}

		output.PrintInfo("Pork range scan\n")

		// A single scan feeds the whole map; scanning per port made a wide
		// range take one full system walk per port.
		scanner := ports.NewScanner()
		busy, err := scanner.GetPortMap()
		if err != nil {
			output.PrintError("Failed to scan ports: %v", err)
			os.Exit(1)
		}

		var nextFree uint32
		foundFree := false

		for port := start; port <= end; port++ {
			listeners := busy[port]
			if len(listeners) == 0 {
				fmt.Printf("%d  ○ free\n", port)
				if !foundFree {
					nextFree = port
					foundFree = true
				}
				continue
			}

			for _, info := range listeners {
				// Mirror the wording used by `pork list`: without elevated
				// privileges the socket is visible but its owner is not.
				program, pid := "Unknown (sudo)", "-"
				if info.PID > 0 {
					program = process.Identify(info.Process, info.Command)
					pid = strconv.Itoa(int(info.PID))
				}
				fmt.Printf("%d  ● busy   %-15s PID %s\n", port, program, pid)
			}
		}

		fmt.Printf("\n● busy   ○ free\n\n")
		if foundFree {
			output.PrintSuccess("Next free port: %d", nextFree)
		} else {
			output.PrintWarning("No free ports in this range")
		}
	},
}

// init registers the rangeCmd as a subcommand of rootCmd.
func init() {
	rootCmd.AddCommand(rangeCmd)
}
