//go:build windows

package process

import "os"

// requestShutdown reports that Windows offers no signal os.Process can deliver
// to ask a process to exit cleanly, so callers go straight to the forced kill.
func requestShutdown(_ *os.Process) (bool, error) {
	return false, nil
}
