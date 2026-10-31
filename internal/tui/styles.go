package tui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/kirre02/net-cli/internal/domain"
)

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	sidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1).
			Width(26)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("252"))

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

var methodStyles = map[domain.HTTPMethod]lipgloss.Style{
	domain.MethodGet: lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")),
	domain.MethodPost: lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")),
	domain.MethodPut: lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")),
	domain.MethodPatch: lipgloss.NewStyle().
		Foreground(lipgloss.Color("44")),
	domain.MethodDelete: lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")),
	domain.MethodHead: lipgloss.NewStyle().
		Foreground(lipgloss.Color("141")),
	domain.MethodOptions: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")),
}

func renderMethod(method domain.HTTPMethod) string {
	style, ok := methodStyles[method]
	if !ok {
		return string(method)
	}

	return style.Render(string(method))
}
