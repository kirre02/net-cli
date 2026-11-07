package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kirre02/net-cli/internal/storage"
)

func Run(store *storage.SQLite) error {
	p := tea.NewProgram(
		initialModel(store),
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	return err
}
