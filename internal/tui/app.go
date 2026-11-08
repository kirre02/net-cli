package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/kirre02/net-cli/internal/storage"
)

func Run(store *storage.SQLite) error {
	p := tea.NewProgram(initialModel(store))

	_, err := p.Run()
	return err
}
