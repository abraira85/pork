// Package main is the entry point for the Pork CLI tool.
// Pork is a modern, lightweight, and cross-platform terminal utility
// designed to inspect, visualize, and free local ports.
package main

import (
	"github.com/outboss/pork/cmd"
)

// main initializes and executes the root Cobra command.
// It delegates the entire execution flow to the cmd package.
func main() {
	cmd.Execute()
}
