//go:build !windows

package process

import (
	"os"
	"syscall"
)

// requestShutdown asks the process to exit cleanly by sending SIGTERM.
// It reports whether the platform supports such a request at all.
func requestShutdown(p *os.Process) (bool, error) {
	return true, p.Signal(syscall.SIGTERM)
}
