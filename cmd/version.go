package cmd

// version is the current release of pork.
//
// It is injected at build time via -ldflags
// (e.g. -X github.com/abraira85/pork/cmd.version=1.2.3) and defaults to "dev"
// for local builds.
var version = "dev"
