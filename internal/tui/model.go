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
	stateActionMenu
	stateActionResult
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
	if err != nil || width <= 0 {
		width = 80
		height = 24
	}

	progWidth := width - 25
	if progWidth < 15 {
		progWidth = 15
	}

	tableHeight := height - 8
	if tableHeight < 5 {
		tableHeight = 5
	}

	columns := []table.Column{
		{Title: "PORT", Width: 8},
		{Title: "PID", Width: 8},
		{Title: "PROGRAM", Width: progWidth},
	}

	m := Model{
		activePorts: activePorts,
		state:       stateTableView,
		actionIndex: 0,
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(m.GenerateRows()),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
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

type tickMsg struct{}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}
