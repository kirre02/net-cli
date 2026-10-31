package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	if m.width < 60 || m.height < 15 {
		return "Terminal is too small. Resize to at least 60x15."
	}

	sidebar := sidebarStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Collections"),
			"",
			mutedStyle.Render("No collections yet"),
		),
	)

	requestPanel := panelStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Request"),
			"",
			mutedStyle.Render("Select a request to get started."),
		),
	)

	responsePanel := panelStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Response"),
			"",
			mutedStyle.Render("Response data will appear here."),
		),
	)

	main := lipgloss.JoinVertical(
		lipgloss.Left,
		requestPanel,
		responsePanel,
	)

	header := titleStyle.Render("NetPilot") +
		mutedStyle.Render("  HTTP client")

	footer := mutedStyle.Render("q  Quit")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			sidebar,
			main,
		),
		"",
		footer,
	)

	return appStyle.Render(content)
}
