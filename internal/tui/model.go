package tui

import (
	"github.com/kirre02/net-cli/internal/domain"
	"github.com/kirre02/net-cli/internal/storage"
)

type focus int

const (
	focusSidebar focus = iota
	focusRequest
	focusResponse
)

type sidebarMode int

const (
	sidebarCollections sidebarMode = iota
	sidebarRequests
)

type model struct {
	store *storage.SQLite

	collections        []domain.Collection
	selectedCollection int

	requests             []domain.Request
	selectedRequestIndex int
	activeRequest        *domain.Request

	sidebarMode  sidebarMode
	focusedPanel focus

	width   int
	height  int
	loading bool
	err     error
}

func initialModel(store *storage.SQLite) *model {
	return &model{
		store:        store,
		focusedPanel: focusSidebar,
		sidebarMode:  sidebarCollections,
	}
}
