package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"github.com/abraira85/pork/internal/output"
	"github.com/abraira85/pork/internal/ports"
	"github.com/abraira85/pork/internal/process"
)

type sessionState int

const (
	stateTableView sessionState = iota
	stateConfirmKill
	stateActionMenu
	stateActionResult
)

// Layout constants for the port table.
const (
	defaultWidth  = 80
	defaultHeight = 24
	chromeHeight  = 8
	minTableRows  = 5
	minProgWidth  = 15
	fixedColsWide = 25
)

type Model struct {
	table        table.Model
	activePorts  []*ports.PortInfo
	state        sessionState
	selectedPort *ports.PortInfo
	actionIndex  int
	resultMsg    string
}

func (m *Model) GenerateRows() []table.Row {
	var rows []table.Row
	for _, p := range m.activePorts {
		pidStr := "-"
		var programStr string

		if p.PID > 0 {
			pidStr = fmt.Sprintf("%d", p.PID)
			programStr = process.Identify(p.Process, p.Command)
		} else {
			programStr = "Unknown (try with sudo)"
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", p.Port),
			pidStr,
			programStr,
		})
	}
	return rows
}

// NewModel initializes the TUI model.
func NewModel(activePorts []*ports.PortInfo) Model {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 || height <= 0 {
		width = defaultWidth
		height = defaultHeight
	}

	m := Model{
		activePorts: activePorts,
		state:       stateTableView,
		actionIndex: 0,
	}

	t := table.New(
		table.WithColumns(columnsFor(width)),
		table.WithRows(m.GenerateRows()),
		table.WithFocused(true),
		table.WithHeight(tableHeightFor(height)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(output.PrimaryColor).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(output.PrimaryColor).
		Bold(true)
	t.SetStyles(s)

	m.table = t
	return m
}

// columnsFor builds the table columns for a given terminal width, giving the
// PROGRAM column whatever is left over.
func columnsFor(width int) []table.Column {
	progWidth := width - fixedColsWide
	if progWidth < minProgWidth {
		progWidth = minProgWidth
	}

	return []table.Column{
		{Title: "PORT", Width: 8},
		{Title: "PID", Width: 8},
		{Title: "PROGRAM", Width: progWidth},
	}
}

// tableHeightFor returns how many rows fit in a terminal of the given height,
// leaving room for the title and help line.
func tableHeightFor(height int) int {
	rows := height - chromeHeight
	if rows < minTableRows {
		rows = minTableRows
	}
	return rows
}

type tickMsg struct{}

// killResultMsg carries the outcome of a termination back to the UI loop.
type killResultMsg struct {
	text string
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

// terminateCmd kills a process off the UI goroutine.
//
// process.Terminate waits out a grace period before forcing the kill, which
// would otherwise freeze the interface for seconds while Update blocks.
func terminateCmd(pid int32, name string) tea.Cmd {
	return func() tea.Msg {
		if err := process.Terminate(pid); err != nil {
			return killResultMsg{text: fmt.Sprintf("Failed to kill process: %v", err)}
		}
		return killResultMsg{text: fmt.Sprintf("Process %d (%s) terminated.", pid, name)}
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}
