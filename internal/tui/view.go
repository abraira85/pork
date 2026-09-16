package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/abraira85/pork/internal/output"
)

func (m Model) View() string {
	switch m.state {
	case stateTableView:
		return m.viewTable()
	case stateActionMenu:
		return m.viewActionMenu()
	case stateConfirmKill:
		return m.viewConfirmKill()
	case stateActionResult:
		return m.viewActionResult()
	}
	return ""
}

// helpLine renders the key hints shown at the bottom of every screen.
func helpLine(text string) string {
	return lipgloss.NewStyle().Foreground(output.MutedColor).MarginTop(1).Render(text)
}

func (m Model) viewTable() string {
	return fmt.Sprintf(
		"\n  %s\n\n%s\n%s\n",
		lipgloss.NewStyle().Foreground(output.PrimaryColor).Bold(true).Render("🐷 Pork Interactive Shell"),
		m.table.View(),
		helpLine("↑/↓: Navigate • Enter: Select • q: Quit"),
	)
}

func (m Model) viewActionMenu() string {
	titleStyle := lipgloss.NewStyle().Foreground(output.PrimaryColor).Bold(true).MarginBottom(1)
	title := titleStyle.Render(fmt.Sprintf("Actions for Port %d (PID: %d)", m.selectedPort.Port, m.selectedPort.PID))

	options := []string{"Kill Process", "Inspect Details", "Cancel"}

	var sb strings.Builder
	sb.WriteString("\n  ")
	sb.WriteString(title)
	sb.WriteString("\n\n")

	for i, opt := range options {
		cursor := "  "
		style := lipgloss.NewStyle()

		if m.actionIndex == i {
			cursor = ">>"
			style = style.Foreground(output.SecondaryColor).Bold(true)
			if i == actionKill { // Kill is red when hovered
				style = style.Foreground(output.DangerColor)
			}
		}

		fmt.Fprintf(&sb, "  %s %s\n", cursor, style.Render(opt))
	}

	sb.WriteString(helpLine("\n  ↑/↓: Navigate • Enter: Select • Esc: Back"))
	sb.WriteString("\n")

	return sb.String()
}

func (m Model) viewConfirmKill() string {
	titleStyle := lipgloss.NewStyle().Foreground(output.WarningColor).Bold(true).MarginBottom(1)
	title := titleStyle.Render("Critical process")

	bodyStyle := lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(output.WarningColor)
	body := bodyStyle.Render(fmt.Sprintf(
		"%s (PID %d) looks like a critical system process.\nKilling it may destabilise your machine.\n\nKill anyway?",
		m.selectedPort.Process, m.selectedPort.PID,
	))

	return fmt.Sprintf("\n  %s\n%s\n\n  %s\n", title, body, helpLine("y: Kill • n/Esc: Back"))
}

func (m Model) viewActionResult() string {
	titleStyle := lipgloss.NewStyle().Foreground(output.PrimaryColor).Bold(true).MarginBottom(1)
	title := titleStyle.Render("Result")

	bodyStyle := lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(output.MutedColor)
	body := bodyStyle.Render(m.resultMsg)

	return fmt.Sprintf("\n  %s\n%s\n\n  %s\n", title, body, helpLine("Enter/Esc: Back"))
}
