// Package cmd contains all the CLI commands for the Pork application.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/abraira85/pork/internal/output"
	"github.com/abraira85/pork/internal/ports"
	"github.com/abraira85/pork/internal/process"
)

// killCmd represents the "kill" command.
// It finds the process(es) listening on the specified port and terminates them,
// asking for confirmation first and flagging critical system processes.
var killCmd = &cobra.Command{
	Use:   "kill [port]",
	Short: "Kill the process using a specific port",
	Long: `Finds the process occupying the specified port and terminates it.

Pork asks the process to shut down cleanly first and only forces the kill if it
is still alive a few seconds later, so sockets and buffers get a chance to close.

Processes that look like critical system components are flagged with an extra
warning before you confirm.`,
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		output.PrintBanner()

		port, err := parsePort(args[0])
		if err != nil {
			output.PrintError("%v", err)
			os.Exit(1)
		}

		scanner := ports.NewScanner()
		listeners, err := scanner.GetPortProcesses(port)
		if err != nil {
			output.PrintError("Failed to inspect port: %v", err)
			os.Exit(1)
		}

		if len(listeners) == 0 {
			output.PrintSuccess("Port %d is already free", port)
			return
		}

		for _, info := range listeners {
			output.PrintInfo("Port %d is used by %s PID %d", port, info.Process, info.PID)
			if info.Command != "" {
				fmt.Printf("Command:\n%s\n", info.Command)
			}
			fmt.Println()
		}

		critical := criticalListeners(listeners)
		if len(critical) > 0 {
			output.PrintWarning("This may be a critical system process:")
			for _, info := range critical {
				fmt.Printf("\nPID      %d\nProcess  %s\nCommand  %s\n", info.PID, info.Process, info.Command)
			}
			fmt.Println()
		}

		if !confirmKill(len(critical) > 0) {
			fmt.Println("Aborted.")
			return
		}

		failures := 0
		for _, info := range listeners {
			if err := process.Terminate(info.PID); err != nil {
				output.PrintError("Failed to kill %s PID %d: %v", info.Process, info.PID, err)
				failures++
				continue
			}
			fmt.Printf("Killed %s process PID %d\n", info.Process, info.PID)
		}

		if failures > 0 {
			os.Exit(1)
		}

		output.PrintSuccess("Freed port %d", port)
	},
}

// criticalListeners returns the subset of listeners whose owning process looks
// like a critical system component.
func criticalListeners(listeners []*ports.PortInfo) []*ports.PortInfo {
	var critical []*ports.PortInfo
	for _, info := range listeners {
		if process.IsCritical(info.Process, info.PID) {
			critical = append(critical, info)
		}
	}
	return critical
}

// confirmKill asks the user to confirm the termination.
//
// A closed or empty stdin counts as "no", so a non-interactive run aborts
// instead of killing something unattended.
func confirmKill(critical bool) bool {
	if critical {
		fmt.Print("Kill anyway? [y/N]: ")
	} else {
		fmt.Print("Kill process? [y/N]: ")
	}

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil && response == "" {
		fmt.Println()
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// init registers the killCmd as a subcommand of rootCmd.
func init() {
	rootCmd.AddCommand(killCmd)
}
