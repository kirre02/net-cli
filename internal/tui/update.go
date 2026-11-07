package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kirre02/net-cli/internal/domain"
	"github.com/kirre02/net-cli/internal/storage"
)

type collectionsLoadedMsg struct {
	collections []domain.Collection
	err         error
}

type requestsLoadedMsg struct {
	collectionID string
	requests     []domain.Request
	err          error
}

func (m model) Init() tea.Cmd {
	return loadCollections(m.store)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case collectionsLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.collections = msg.collections
		m.selectedCollection = 0
		m.err = nil

	case requestsLoadedMsg:
		m.loading = false

		// Ignore results if the selected collection has changed.
		if len(m.collections) == 0 ||
			m.collections[m.selectedCollection].ID != msg.collectionID {
			return m, nil
		}

		m.err = msg.err

		if msg.err != nil {
			return m, nil
		}

		m.requests = msg.requests
		m.selectedRequestIndex = 0
		m.activeRequest = nil
		m.sidebarMode = sidebarRequests

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.sidebarMode == sidebarRequests {
				m.sidebarMode = sidebarCollections
				m.requests = nil
				m.selectedRequestIndex = 0
				m.activeRequest = nil
				m.err = nil
			}

		case "j", "down":
			if m.focusedPanel == focusSidebar {
				m.moveSelection(1)
			}

		case "k", "up":
			if m.focusedPanel == focusSidebar {
				m.moveSelection(-1)
			}

		case "h", "left":
			if m.focusedPanel > focusSidebar {
				m.focusedPanel--
			}

		case "l", "right":
			if m.focusedPanel < focusResponse {
				m.focusedPanel++
			}

		case "enter":
			if m.focusedPanel == focusSidebar {
				if m.sidebarMode == sidebarCollections {
					if len(m.collections) == 0 || m.loading {
						return m, nil
					}

					collectionID := m.collections[m.selectedCollection].ID
					m.loading = true
					m.err = nil

					return m, loadRequests(m.store, collectionID)
				}

				if m.sidebarMode == sidebarRequests {
					if len(m.requests) == 0 {
						return m, nil
					}

					m.activeRequest = &m.requests[m.selectedRequestIndex]
				}
			}
		}
	}

	return m, nil
}

func loadCollections(store *storage.SQLite) tea.Cmd {
	return func() tea.Msg {
		collections, err := store.GetCollections(context.Background())

		return collectionsLoadedMsg{
			collections: collections,
			err:         err,
		}
	}
}

func loadRequests(store *storage.SQLite, collectionID string) tea.Cmd {
	return func() tea.Msg {
		requests, err := store.GetRequestsByCollection(
			context.Background(),
			collectionID,
		)

		return requestsLoadedMsg{
			collectionID: collectionID,
			requests:     requests,
			err:          err,
		}
	}
}

func (m *model) moveSelection(delta int) {
	var count int

	if m.sidebarMode == sidebarCollections {
		count = len(m.collections)
	} else {
		count = len(m.requests)
	}

	if count == 0 {
		return
	}

	if m.sidebarMode == sidebarCollections {
		m.selectedCollection += delta

		if m.selectedCollection < 0 {
			m.selectedCollection = 0
		}

		if m.selectedCollection >= count {
			m.selectedCollection = count - 1
		}

		return
	}

	m.selectedRequestIndex += delta

	if m.selectedRequestIndex < 0 {
		m.selectedRequestIndex = 0
	}

	if m.selectedRequestIndex >= count {
		m.selectedRequestIndex = count - 1
	}
}
