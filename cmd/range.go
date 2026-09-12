// Package cmd contains all the CLI commands for the Pork application.
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/abraira85/pork/internal/output"
	"github.com/abraira85/pork/internal/ports"
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
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintBanner()
		startStr := args[0]
		endStr := args[1]

		startNum, err1 := strconv.ParseUint(startStr, 10, 32)
		endNum, err2 := strconv.ParseUint(endStr, 10, 32)

		if err1 != nil || err2 != nil || startNum > endNum {
			output.PrintError("Invalid range: %s to %s", startStr, endStr)
			os.Exit(1)
		}

		output.PrintInfo("Pork range scan\n")

		scanner := ports.NewScanner()

		var nextFree uint32

		for i := startNum; i <= endNum; i++ {
			info, err := scanner.GetPortInfo(uint32(i))
			if err != nil || info == nil {
				fmt.Printf("%d  ○ free\n", i)
				if nextFree == 0 {
					nextFree = uint32(i)
				}
			} else {
				fmt.Printf("%d  ● busy   %-15s PID %d\n", i, info.Process, info.PID)
			}
		}

		fmt.Printf("\n● busy   ○ free\n\n")
		if nextFree != 0 {
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
