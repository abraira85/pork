// Package ports provides data structures and logic for scanning local network ports.
package ports

import (
	"fmt"
	"sort"

	gopsnet "github.com/shirou/gopsutil/v3/net"
	porkprocess "github.com/outboss/pork/internal/process"
)

// Scanner handles scanning for active ports and mapping them to processes.
type Scanner struct{}

// NewScanner creates a new instance of Scanner.
func NewScanner() *Scanner {
	return &Scanner{}
}

// GetActivePorts retrieves a list of all currently listening TCP ports.
// It maps the network connections to their respective system processes.
// Returns a slice of PortInfo sorted by port number.
func (s *Scanner) GetActivePorts() ([]*PortInfo, error) {
	// "tcp" means ipv4 and ipv6 TCP connections
	conns, err := gopsnet.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	var results []*PortInfo
	seen := make(map[uint32]bool)

	for _, conn := range conns {
		// We are only interested in ports that are currently listening
		if conn.Status != "LISTEN" {
			continue
		}

		port := conn.Laddr.Port
		
		// Avoid duplicate entries if multiple interfaces listen on the same port
		if seen[port] {
			continue
		}
		seen[port] = true

		info := &PortInfo{
			Port:     port,
			Protocol: "TCP",
			Address:  fmt.Sprintf("%s:%d", conn.Laddr.IP, port),
			Status:   conn.Status,
			PID:      conn.Pid,
		}

		// If we successfully found a PID associated with this port, fetch its details.
		// Note: Requires elevated privileges on some OS to see PIDs of other users.
		if conn.Pid > 0 {
			if procInfo, err := porkprocess.GetInfo(conn.Pid); err == nil {
				info.Process = procInfo.Name
				info.Command = procInfo.Command
				info.User = procInfo.User
				info.Path = procInfo.Path
				info.CPU = procInfo.CPU
				info.MemoryMB = procInfo.MemoryMB
			}
		}

		results = append(results, info)
	}

	// Sort the results by port number for better readability
	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results, nil
}

// GetPortInfo retrieves information for a specific port.
// Returns nil if the port is not currently listening.
func (s *Scanner) GetPortInfo(port uint32) (*PortInfo, error) {
	ports, err := s.GetActivePorts()
	if err != nil {
		return nil, err
	}

	for _, p := range ports {
		if p.Port == port {
			return p, nil
		}
	}

	return nil, nil // Port is free
}

// IsFree checks if a specific port is currently free (not listening).
func (s *Scanner) IsFree(port uint32) (bool, error) {
	info, err := s.GetPortInfo(port)
	if err != nil {
		return false, err
	}
	return info == nil, nil
}
