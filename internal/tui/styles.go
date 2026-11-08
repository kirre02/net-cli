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
			Padding(1)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("252"))

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	selectedCollectionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")).
				Bold(true).
				PaddingLeft(1)

	collectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			PaddingLeft(1)

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	footerDescriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245"))

	selectedFieldStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("42")).
				Padding(0, 1)

	urlCursorStyle = lipgloss.NewStyle().
			Reverse(true)
)

var methodStyles = map[domain.HTTPMethod]lipgloss.Style{
	domain.MethodGet:     lipgloss.NewStyle().Foreground(lipgloss.Color("42")),
	domain.MethodPost:    lipgloss.NewStyle().Foreground(lipgloss.Color("214")),
	domain.MethodPut:     lipgloss.NewStyle().Foreground(lipgloss.Color("39")),
	domain.MethodPatch:   lipgloss.NewStyle().Foreground(lipgloss.Color("51")),
	domain.MethodDelete:  lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
	domain.MethodHead:    lipgloss.NewStyle().Foreground(lipgloss.Color("141")),
	domain.MethodOptions: lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
}

func renderMethod(method domain.HTTPMethod) string {
	style, ok := methodStyles[method]
	if !ok {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Render(string(method))
	}

	return style.Render(string(method))
}

func renderKeyBinding(key, description string) string {
	return keyStyle.Render(key) + " " +
		footerDescriptionStyle.Render(description)
}
