package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/abraira85/pork/internal/ports"
	"github.com/abraira85/pork/internal/process"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		tableHeight := msg.Height - 8
		if tableHeight < 5 {
			tableHeight = 5
		}
		m.table.SetHeight(tableHeight)
		return m, nil
	case tickMsg:
		m.reloadPorts()
		return m, tickCmd()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.state == stateActionMenu || m.state == stateActionResult {
				m.state = stateTableView
				m.actionIndex = 0
			}
		case "enter":
			if m.state == stateTableView {
				if len(m.activePorts) > 0 {
					// Get selected port
					m.selectedPort = m.activePorts[m.table.Cursor()]
					m.state = stateActionMenu
					m.actionIndex = 0
				}
			} else if m.state == stateActionMenu {
				return m.handleActionSelection()
			} else if m.state == stateActionResult {
				m.state = stateTableView
			}
		case "up", "k":
			if m.state == stateActionMenu {
				if m.actionIndex > 0 {
					m.actionIndex--
				}
			} else if m.state == stateTableView {
				var cmd tea.Cmd
				m.table, cmd = m.table.Update(msg)
				return m, cmd
			}
		case "down", "j":
			if m.state == stateActionMenu {
				if m.actionIndex < 2 { // 3 options: Kill, Inspect, Cancel
					m.actionIndex++
				}
			} else if m.state == stateTableView {
				var cmd tea.Cmd
				m.table, cmd = m.table.Update(msg)
				return m, cmd
			}
		default:
			if m.state == stateTableView {
				var cmd tea.Cmd
				m.table, cmd = m.table.Update(msg)
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m *Model) handleActionSelection() (tea.Model, tea.Cmd) {
	switch m.actionIndex {
	case 0:
		// Kill process
		if m.selectedPort.PID == 0 {
			m.resultMsg = "Cannot kill process: PID is unknown (try with sudo)."
			m.state = stateActionResult
			return *m, nil
		}

		if process.IsCritical(m.selectedPort.Process, m.selectedPort.PID) {
			m.resultMsg = fmt.Sprintf("Safety check failed: Process '%s' is critical.", m.selectedPort.Process)
			m.state = stateActionResult
			return *m, nil
		}

		p, err := os.FindProcess(int(m.selectedPort.PID))
		if err != nil {
			m.resultMsg = fmt.Sprintf("Error finding process: %v", err)
			m.state = stateActionResult
			return *m, nil
		}

		err = p.Kill()
		if err != nil {
			m.resultMsg = fmt.Sprintf("Failed to kill process: %v", err)
		} else {
			m.resultMsg = fmt.Sprintf("Process %d (%s) killed successfully.", m.selectedPort.PID, m.selectedPort.Process)
			m.reloadPorts()
		}

		m.state = stateActionResult
		return *m, nil
	case 1:
		// Inspect
		m.resultMsg = fmt.Sprintf("Port: %d\nPID: %d\nProcess: %s\nCommand: %s",
			m.selectedPort.Port, m.selectedPort.PID, m.selectedPort.Process, m.selectedPort.Command)
		m.state = stateActionResult
		return *m, nil
	case 2:
		// Cancel
		m.state = stateTableView
		return *m, nil
	}
	return *m, nil
}

func (m *Model) reloadPorts() {
	scanner := ports.NewScanner()
	activePorts, _ := scanner.GetActivePorts()
	m.activePorts = activePorts
	rows := m.GenerateRows()
	m.table.SetRows(rows)

	// Clamp cursor if the list shrunk
	if m.table.Cursor() >= len(rows) {
		if len(rows) > 0 {
			m.table.SetCursor(len(rows) - 1)
		} else {
			m.table.SetCursor(0)
		}
	}
}
