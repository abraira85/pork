// Package cmd contains all the CLI commands for the Pork application.
package cmd

import (
	"fmt"
	"strconv"

	"github.com/abraira85/pork/internal/ports"
)

// parsePort converts a CLI argument into a port number.
//
// ParseUint alone is not enough: it happily accepts values far beyond the
// 16-bit port space, which made `pork 99999` report a non-existent port as free
// and let `pork range 1 4000000000` spin for billions of iterations.
func parsePort(arg string) (uint32, error) {
	n, err := strconv.ParseUint(arg, 10, 32)
	if err != nil || n < uint64(ports.MinPort) || n > uint64(ports.MaxPort) {
		return 0, fmt.Errorf("invalid port number: %s (expected %d-%d)", arg, ports.MinPort, ports.MaxPort)
	}
	return uint32(n), nil
}

// isNumeric reports whether the argument even looks like a port number.
//
// The root command doubles as `pork <port>`, so a mistyped subcommand lands
// here too; telling the user "invalid port number: lst" would be misleading.
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
