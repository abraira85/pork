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
	case stateActionResult:
		return m.viewActionResult()
	}
	return ""
}

func (m Model) viewTable() string {
	helpStyle := lipgloss.NewStyle().Foreground(output.MutedColor).MarginTop(1)
	helpText := helpStyle.Render("↑/↓: Navigate • Enter: Select • q/Esc: Quit")

	return fmt.Sprintf(
		"\n  %s\n\n%s\n%s\n",
		lipgloss.NewStyle().Foreground(output.PrimaryColor).Bold(true).Render("🐷 Pork Interactive Shell"),
		m.table.View(),
		helpText,
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
			if i == 0 { // Kill is red when hovered
				style = style.Foreground(output.DangerColor)
			}
		}

		fmt.Fprintf(&sb, "  %s %s\n", cursor, style.Render(opt))
	}

	helpStyle := lipgloss.NewStyle().Foreground(output.MutedColor).MarginTop(1)
	sb.WriteString(helpStyle.Render("\n  ↑/↓: Navigate • Enter: Select • Esc: Back"))
	sb.WriteString("\n")

	return sb.String()
}

func (m Model) viewActionResult() string {
	titleStyle := lipgloss.NewStyle().Foreground(output.PrimaryColor).Bold(true).MarginBottom(1)
	title := titleStyle.Render("Result")

	bodyStyle := lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(output.MutedColor)
	body := bodyStyle.Render(m.resultMsg)

	helpStyle := lipgloss.NewStyle().Foreground(output.MutedColor).MarginTop(1)
	helpText := helpStyle.Render("Enter/Esc: Back")

	return fmt.Sprintf("\n  %s\n%s\n\n  %s\n", title, body, helpText)
}
