package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/abraira85/pork/internal/ports"
)

func key(s string) tea.KeyMsg {
	switch s {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

// quits reports whether the returned command would terminate the program.
func quits(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

func modelWith(p *ports.PortInfo) Model {
	m := NewModel([]*ports.PortInfo{p})
	return m
}

func ordinaryPort() *ports.PortInfo {
	return &ports.PortInfo{Port: 3000, PID: 4242, Process: "node", Command: "node server.js"}
}

func criticalPort() *ports.PortInfo {
	return &ports.PortInfo{Port: 22, PID: 4243, Process: "sshd", Command: "/usr/sbin/sshd -D"}
}

func TestQuitKeys(t *testing.T) {
	m := modelWith(ordinaryPort())

	_, cmd := m.Update(key("q"))
	if !quits(cmd) {
		t.Error("q should quit from the table view, which is what the help line promises")
	}

	// Opening the action menu and pressing q must not tear down the session.
	next, _ := m.Update(key("enter"))
	menu := next.(Model)
	if menu.state != stateActionMenu {
		t.Fatalf("enter should open the action menu, got state %v", menu.state)
	}

	_, cmd = menu.Update(key("q"))
	if quits(cmd) {
		t.Error("q should not quit while the action menu is open")
	}

	_, cmd = menu.Update(key("ctrl+c"))
	if !quits(cmd) {
		t.Error("ctrl+c must always quit")
	}
}

func TestEscReturnsFromMenu(t *testing.T) {
	m := modelWith(ordinaryPort())

	next, _ := m.Update(key("enter"))
	next, _ = next.(Model).Update(key("esc"))

	if got := next.(Model).state; got != stateTableView {
		t.Errorf("esc should return to the table view, got state %v", got)
	}
}

func TestCriticalProcessAsksInsteadOfRefusing(t *testing.T) {
	m := modelWith(criticalPort())

	next, _ := m.Update(key("enter"))              // open action menu
	next, cmd := next.(Model).Update(key("enter")) // choose "Kill Process"

	confirm := next.(Model)
	if confirm.state != stateConfirmKill {
		t.Fatalf("a critical process should prompt for confirmation, got state %v", confirm.state)
	}
	if cmd != nil {
		t.Error("nothing should be terminated before the user confirms")
	}

	// Declining returns to the menu without killing anything.
	back, cmd := confirm.Update(key("n"))
	if got := back.(Model).state; got != stateActionMenu {
		t.Errorf("n should return to the action menu, got state %v", got)
	}
	if cmd != nil {
		t.Error("declining must not schedule a termination")
	}
}

func TestKillRunsOffTheUILoop(t *testing.T) {
	m := modelWith(ordinaryPort())

	next, _ := m.Update(key("enter"))
	next, cmd := next.(Model).Update(key("enter"))

	started := next.(Model)
	if started.state != stateActionResult {
		t.Fatalf("killing should show a result screen immediately, got state %v", started.state)
	}
	if cmd == nil {
		t.Fatal("the termination must be handed to a background command, not run inline")
	}
	if started.resultMsg == "" {
		t.Error("the user should see progress while the kill is in flight")
	}
}

func TestUnknownPIDCannotBeKilled(t *testing.T) {
	m := modelWith(&ports.PortInfo{Port: 5432, PID: 0})

	next, _ := m.Update(key("enter"))
	next, cmd := next.(Model).Update(key("enter"))

	result := next.(Model)
	if result.state != stateActionResult {
		t.Fatalf("expected a result screen, got state %v", result.state)
	}
	if cmd != nil {
		t.Error("no termination should be scheduled for an unknown PID")
	}
}

func TestResizeAdjustsColumns(t *testing.T) {
	m := modelWith(ordinaryPort())

	next, _ := m.Update(tea.WindowSizeMsg{Width: 200, Height: 60})
	wide := next.(Model).table.Columns()

	next, _ = m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	narrow := next.(Model).table.Columns()

	if wide[2].Width <= narrow[2].Width {
		t.Errorf("the PROGRAM column should grow with the terminal: wide=%d narrow=%d",
			wide[2].Width, narrow[2].Width)
	}
	if narrow[2].Width < minProgWidth {
		t.Errorf("the PROGRAM column should never shrink below %d, got %d", minProgWidth, narrow[2].Width)
	}
}
