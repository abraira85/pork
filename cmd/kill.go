// Package cmd contains all the CLI commands for the Pork application.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/outboss/pork/internal/output"
	"github.com/outboss/pork/internal/ports"
	"github.com/outboss/pork/internal/process"
)

// killCmd represents the "kill" command.
// It finds the process listening on the specified port and terminates it.
// By default, it uses a graceful termination signal and asks for user confirmation,
// unless overridden by flags.
var killCmd = &cobra.Command{
	Use:   "kill [port]",
	Short: "Kill the process using a specific port",
	Long: `Finds the process occupying the specified port and attempts to kill it.
It incorporates safety checks to prevent accidentally killing critical system processes.
Requires confirmation before terminating the process unless --yes is provided.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintBanner()
		portStr := args[0]
		portNum, err := strconv.ParseUint(portStr, 10, 32)
		if err != nil {
			output.PrintError("Invalid port number: %s", portStr)
			os.Exit(1)
		}

		scanner := ports.NewScanner()
		info, err := scanner.GetPortInfo(uint32(portNum))
		if err != nil {
			output.PrintError("Failed to inspect port: %v", err)
			os.Exit(1)
		}

		if info == nil {
			output.PrintSuccess("Port %d is already free", portNum)
			return
		}

		output.PrintInfo("Port %d is used by %s PID %d\n", portNum, info.Process, info.PID)
		if info.Command != "" {
			fmt.Printf("Command:\n%s\n\n", info.Command)
		}

		if process.IsCritical(info.Process, info.PID) {
			output.PrintWarning("This process may be important:")
			fmt.Printf("\nPID      %d\nProcess  %s\nCommand  %s\n\n", info.PID, info.Process, info.Command)
			fmt.Print("Kill anyway? [y/N]: ")
		} else {
			fmt.Print("Kill process? [y/N]: ")
		}

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			if err := process.Terminate(info.PID); err != nil {
				output.PrintError("Failed to kill process: %v", err)
				os.Exit(1)
			}
			output.PrintSuccess("Freed port %d", portNum)
			fmt.Printf("Killed %s process PID %d\n", info.Process, info.PID)
		} else {
			fmt.Println("Aborted.")
		}
	},
}

// init registers the killCmd as a subcommand of rootCmd.
func init() {
	rootCmd.AddCommand(killCmd)
}
