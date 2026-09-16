// Package cmd contains all the CLI commands for the Pork application.
// It uses the Cobra framework to define and route commands, arguments, and flags.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/abraira85/pork/internal/output"
	"github.com/abraira85/pork/internal/ports"
)

// rootCmd represents the base command when called without any subcommands.
// It also acts as the "inspect" command if a port number is provided directly (e.g. `pork 3000`).
var rootCmd = &cobra.Command{
	Use:   "pork [port]",
	Short: "Pork is a tiny terminal tool to inspect, visualize and free local ports.",
	Long: `Pork is a tiny terminal tool to inspect, visualize and free local ports.
It provides a beautiful and simple interface over standard tools like lsof or netstat.`,
	Version: version,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintBanner()
		if len(args) == 0 {
			// If no port is provided, show the default help message
			_ = cmd.Help()
			return
		}

		if !isNumeric(args[0]) {
			output.PrintError("Unknown command %q. Run 'pork --help' to see the available commands.", args[0])
			os.Exit(1)
		}

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

		output.PrintPortInspection(port, listeners)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Add custom colored help template
	customHelpTemplate := `
{{.Long}}

` + "\033[1;38;2;255;42;117m" + `Usage:` + "\033[0m" + `
  {{.UseLine}}
{{if gt (len .Aliases) 0}}
` + "\033[1;38;2;255;42;117m" + `Aliases:` + "\033[0m" + `
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

` + "\033[1;38;2;255;42;117m" + `Examples:` + "\033[0m" + `
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

` + "\033[1;38;2;255;42;117m" + `Available Commands:` + "\033[0m" + `{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  ` + "\033[38;2;255;126;179m" + `{{rpad .Name .NamePadding }}` + "\033[0m" + ` {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

` + "\033[1;38;2;255;42;117m" + `Flags:` + "\033[0m" + `
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

` + "\033[1;38;2;255;42;117m" + `Global Flags:` + "\033[0m" + `
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

` + "\033[1;38;2;255;42;117m" + `Additional help topics:` + "\033[0m" + `{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "` + "\033[38;2;255;126;179m" + `{{.CommandPath}} [command] --help` + "\033[0m" + `" for more information about a command.{{end}}
`

	rootCmd.SetHelpTemplate(customHelpTemplate)
	rootCmd.SetUsageTemplate(customHelpTemplate)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
