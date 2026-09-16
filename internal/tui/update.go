package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/abraira85/pork/internal/ports"
	"github.com/abraira85/pork/internal/process"
)

// Action menu entries, in display order.
const (
	actionKill = iota
	actionInspect
	actionCancel
	actionCount
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(tableHeightFor(msg.Height))
		m.table.SetColumns(columnsFor(msg.Width))
		return m, nil
	case tickMsg:
		m.reloadPorts()
		return m, tickCmd()
	case killResultMsg:
		m.resultMsg = msg.text
		m.reloadPorts()
		m.state = stateActionResult
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// ctrl+c always exits. "q" only does so from the table, so that it cannot
	// tear down the session while a confirmation is on screen.
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "q":
		if m.state == stateTableView {
			return m, tea.Quit
		}
	}

	switch m.state {
	case stateTableView:
		return m.handleTableKey(msg)
	case stateActionMenu:
		return m.handleMenuKey(msg)
	case stateConfirmKill:
		return m.handleConfirmKey(msg)
	case stateActionResult:
		switch msg.String() {
		case "enter", "esc":
			m.state = stateTableView
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleTableKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "enter" {
		if len(m.activePorts) > 0 && m.table.Cursor() < len(m.activePorts) {
			m.selectedPort = m.activePorts[m.table.Cursor()]
			m.state = stateActionMenu
			m.actionIndex = 0
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateTableView
		m.actionIndex = 0
	case "up", "k":
		if m.actionIndex > 0 {
			m.actionIndex--
		}
	case "down", "j":
		if m.actionIndex < actionCount-1 {
			m.actionIndex++
		}
	case "enter":
		return m.handleActionSelection()
	}
	return m, nil
}

func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		return m.startKill()
	case "n", "N", "esc":
		m.state = stateActionMenu
	}
	return m, nil
}

func (m *Model) handleActionSelection() (tea.Model, tea.Cmd) {
	switch m.actionIndex {
	case actionKill:
		if m.selectedPort.PID == 0 {
			m.resultMsg = "Cannot kill process: PID is unknown (try with sudo)."
			m.state = stateActionResult
			return *m, nil
		}

		// Critical processes are gated behind an explicit confirmation rather
		// than refused outright, so the TUI and `pork kill` behave alike.
		if process.IsCritical(m.selectedPort.Process, m.selectedPort.PID) {
			m.state = stateConfirmKill
			return *m, nil
		}

		return m.startKill()
	case actionInspect:
		m.resultMsg = fmt.Sprintf("Port: %d\nPID: %d\nProcess: %s\nCommand: %s",
			m.selectedPort.Port, m.selectedPort.PID, m.selectedPort.Process, m.selectedPort.Command)
		m.state = stateActionResult
		return *m, nil
	case actionCancel:
		m.state = stateTableView
		return *m, nil
	}
	return *m, nil
}

// startKill hands the termination off to a background command and shows a
// placeholder result until it reports back.
func (m *Model) startKill() (tea.Model, tea.Cmd) {
	target := m.selectedPort
	m.resultMsg = fmt.Sprintf("Terminating %s (PID %d)...", target.Process, target.PID)
	m.state = stateActionResult
	return *m, terminateCmd(target.PID, target.Process)
}

func (m *Model) reloadPorts() {
	scanner := ports.NewScanner()
	activePorts, err := scanner.GetActivePorts()
	if err != nil {
		// Keep showing the last good snapshot rather than blanking the table.
		return
	}

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
