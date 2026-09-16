// Package ports provides data structures and logic for scanning local network ports.
package ports

import (
	"fmt"
	"sort"

	gopsnet "github.com/shirou/gopsutil/v3/net"

	porkprocess "github.com/abraira85/pork/internal/process"
)

// MinPort and MaxPort bound the range of addressable TCP ports. Port 0 is
// excluded on purpose: it means "let the kernel pick one" and can never be
// inspected, freed or killed.
const (
	MinPort uint32 = 1
	MaxPort uint32 = 65535
)

// Scanner handles scanning for active ports and mapping them to processes.
type Scanner struct{}

// NewScanner creates a new instance of Scanner.
func NewScanner() *Scanner {
	return &Scanner{}
}

// listenerKey identifies a distinct listener. Deduplicating on the port alone
// would hide a second process bound to the same port on a different interface;
// keying on the owning PID as well collapses the common case (one process
// listening over both IPv4 and IPv6) without losing that information.
type listenerKey struct {
	port uint32
	pid  int32
}

// GetActivePorts retrieves a list of all currently listening TCP ports.
// It maps the network connections to their respective system processes.
// Returns a slice of PortInfo sorted by port number, then by PID.
func (s *Scanner) GetActivePorts() ([]*PortInfo, error) {
	// "tcp" means ipv4 and ipv6 TCP connections
	conns, err := gopsnet.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	var results []*PortInfo
	seen := make(map[listenerKey]bool)

	// Process metadata is looked up once per PID rather than once per socket:
	// a process listening on many ports is common and the lookup is the
	// expensive part of a scan.
	procs := make(map[int32]*porkprocess.Info)

	for _, conn := range conns {
		// We are only interested in ports that are currently listening
		if conn.Status != "LISTEN" {
			continue
		}

		key := listenerKey{port: conn.Laddr.Port, pid: conn.Pid}
		if seen[key] {
			continue
		}
		seen[key] = true

		info := &PortInfo{
			Port:     conn.Laddr.Port,
			Protocol: "TCP",
			Address:  fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port),
			Status:   conn.Status,
			PID:      conn.Pid,
		}

		// If we successfully found a PID associated with this port, fetch its details.
		// Note: Requires elevated privileges on some OS to see PIDs of other users.
		if conn.Pid > 0 {
			procInfo, cached := procs[conn.Pid]
			if !cached {
				procInfo, _ = porkprocess.GetInfo(conn.Pid)
				procs[conn.Pid] = procInfo
			}
			if procInfo != nil {
				info.Process = procInfo.Name
				info.Command = procInfo.Command
				info.User = procInfo.User
			}
		}

		results = append(results, info)
	}

	// Sort for stable, readable output.
	sort.Slice(results, func(i, j int) bool {
		if results[i].Port != results[j].Port {
			return results[i].Port < results[j].Port
		}
		return results[i].PID < results[j].PID
	})

	return results, nil
}

// GetPortMap returns every listening port indexed by port number, using a
// single system scan.
//
// Commands that need to test many ports (free, range) must use this instead of
// calling GetPortInfo in a loop: each GetPortInfo call is a full enumeration of
// the machine's sockets and their owning processes, which turns a range scan
// into thousands of redundant system walks.
//
// When several processes share a port, the entry with the lowest PID wins;
// GetPortProcesses exposes the full set.
func (s *Scanner) GetPortMap() (map[uint32][]*PortInfo, error) {
	activePorts, err := s.GetActivePorts()
	if err != nil {
		return nil, err
	}

	byPort := make(map[uint32][]*PortInfo, len(activePorts))
	for _, p := range activePorts {
		byPort[p.Port] = append(byPort[p.Port], p)
	}

	return byPort, nil
}

// GetPortProcesses retrieves every listener bound to a specific port.
// Returns an empty slice if the port is not currently in use.
func (s *Scanner) GetPortProcesses(port uint32) ([]*PortInfo, error) {
	activePorts, err := s.GetActivePorts()
	if err != nil {
		return nil, err
	}

	var matches []*PortInfo
	for _, p := range activePorts {
		if p.Port == port {
			matches = append(matches, p)
		}
	}

	return matches, nil
}
