package tui

import (
	"context"
	"errors"
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"
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

type requestSavedMsg struct {
	request domain.Request
	err     error
}

func (m *model) Init() tea.Cmd {
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
		m.selectedField = fieldMethod
		m.sidebarMode = sidebarRequests
		m.updateActiveRequest()

	case requestSavedMsg:
		m.saving = false

		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		for i := range m.requests {
			if m.requests[i].ID == msg.request.ID {
				m.requests[i] = msg.request
				break
			}
		}

		m.updateActiveRequest()

		m.editingURL = false
		m.urlDraft = ""
		m.urlCursor = 0

		m.editingHeader = false
		m.editingHeaders = false
		m.headerDraftKey = ""
		m.headerDraftValue = ""
		m.headerOriginalKey = ""
		m.headerCursor = 0

		m.err = nil

	case tea.KeyPressMsg:
		if m.saving {
			return m, nil
		}

		if m.editingURL {
			return m, m.handleURLInput(msg)
		}

		if m.editingHeader {
			return m, m.handleHeaderInput(msg)
		}

		if m.editingHeaders {
			return m, m.handleHeaderNavigation(msg)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.sidebarMode == sidebarRequests {
				m.sidebarMode = sidebarCollections
				m.requests = nil
				m.selectedRequestIndex = 0
				m.activeRequest = nil
				m.selectedField = fieldMethod
				m.err = nil
			}

		case "j", "down":
			switch m.focusedPanel {
			case focusSidebar:
				m.moveSelection(1)
			case focusRequest:
				m.moveRequestField(1)
			}

		case "k", "up":
			switch m.focusedPanel {
			case focusSidebar:
				m.moveSelection(-1)
			case focusRequest:
				m.moveRequestField(-1)
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
			}

			if m.focusedPanel == focusRequest &&
				m.activeRequest != nil {
				switch m.selectedField {
				case fieldURL:
					m.startURLEditing()
				case fieldHeaders:
					m.startHeaderEditing()
				}
			}
		}
	}

	return m, nil
}

func (m *model) startURLEditing() {
	if m.activeRequest == nil {
		return
	}

	m.urlDraft = m.activeRequest.URL
	m.urlCursor = len([]rune(m.urlDraft))
	m.editingURL = true
	m.err = nil
}

func (m *model) handleURLInput(msg tea.KeyPressMsg) tea.Cmd {
	saveKey, _ := saveShortcut()

	switch msg.String() {
	case "esc":
		m.editingURL = false
		m.urlDraft = ""
		m.urlCursor = 0
		m.err = nil

	case saveKey:
		if m.activeRequest == nil {
			return nil
		}

		request := *m.activeRequest
		request.URL = m.urlDraft

		m.saving = true
		m.err = nil

		return saveRequest(m.store, request)

	case "left":
		if m.urlCursor > 0 {
			m.urlCursor--
		}

	case "right":
		if m.urlCursor < len([]rune(m.urlDraft)) {
			m.urlCursor++
		}

	case "home":
		m.urlCursor = 0

	case "end":
		m.urlCursor = len([]rune(m.urlDraft))

	case "backspace":
		if m.urlCursor > 0 {
			runes := []rune(m.urlDraft)
			runes = append(runes[:m.urlCursor-1], runes[m.urlCursor:]...)
			m.urlDraft = string(runes)
			m.urlCursor--
		}

	case "delete":
		runes := []rune(m.urlDraft)

		if m.urlCursor < len(runes) {
			runes = append(runes[:m.urlCursor], runes[m.urlCursor+1:]...)
			m.urlDraft = string(runes)
		}

	case "enter":
		// Keep editing until the user saves or cancels.
		return nil

	default:
		// Insert typed characters at the cursor position.
		if msg.Text != "" {
			runes := []rune(m.urlDraft)
			insert := []rune(msg.Text)

			updated := make([]rune, 0, len(runes)+len(insert))
			updated = append(updated, runes[:m.urlCursor]...)
			updated = append(updated, insert...)
			updated = append(updated, runes[m.urlCursor:]...)

			m.urlDraft = string(updated)
			m.urlCursor += len(insert)
		}
	}

	return nil
}

func (m *model) startHeaderEditing() {
	if m.activeRequest == nil {
		return
	}

	m.editingHeaders = true
	m.selectedHeader = 0
	m.err = nil
}

func (m *model) handleHeaderNavigation(msg tea.KeyPressMsg) tea.Cmd {
	keys := sortedHeaderKeys(m.activeRequest.Headers)

	switch msg.String() {
	case "esc":
		m.editingHeaders = false
		m.selectedHeader = 0
		m.err = nil

	case "j", "down":
		if m.selectedHeader < len(keys)-1 {
			m.selectedHeader++
		}

	case "k", "up":
		if m.selectedHeader > 0 {
			m.selectedHeader--
		}

	case "a":
		m.startHeaderInput(true)

	case "e":
		if len(keys) > 0 {
			m.startHeaderInput(false)
		}

	case "d":
		if len(keys) == 0 {
			return nil
		}

		request := *m.activeRequest
		request.Headers = cloneHeaders(request.Headers)

		delete(request.Headers, keys[m.selectedHeader])

		m.saving = true
		m.err = nil

		return saveRequest(m.store, request)
	}

	return nil
}

func (m *model) startHeaderInput(isNew bool) {
	m.headerDraftKey = ""
	m.headerDraftValue = ""
	m.headerOriginalKey = ""
	m.headerEditField = headerKeyField

	if !isNew {
		keys := sortedHeaderKeys(m.activeRequest.Headers)

		if m.selectedHeader < 0 || m.selectedHeader >= len(keys) {
			return
		}

		key := keys[m.selectedHeader]

		m.headerDraftKey = key
		m.headerDraftValue = m.activeRequest.Headers[key]
		m.headerOriginalKey = key
	}

	m.headerCursor = len([]rune(m.headerDraftKey))
	m.editingHeader = true
	m.err = nil
}

func (m *model) handleHeaderInput(msg tea.KeyPressMsg) tea.Cmd {
	saveKey, _ := saveShortcut()

	switch msg.String() {
	case "esc":
		m.editingHeader = false
		m.headerDraftKey = ""
		m.headerDraftValue = ""
		m.headerOriginalKey = ""
		m.headerCursor = 0
		m.err = nil

	case "tab":
		if m.headerEditField == headerKeyField {
			m.headerEditField = headerValueField
			m.headerCursor = len([]rune(m.headerDraftValue))
		} else {
			m.headerEditField = headerKeyField
			m.headerCursor = len([]rune(m.headerDraftKey))
		}

	case saveKey:
		key := strings.TrimSpace(m.headerDraftKey)

		if key == "" {
			m.err = errors.New("header name cannot be empty")
			return nil
		}

		if m.activeRequest == nil {
			return nil
		}

		request := *m.activeRequest
		request.Headers = cloneHeaders(request.Headers)

		if m.headerOriginalKey != "" {
			delete(request.Headers, m.headerOriginalKey)
		}

		for existingKey := range request.Headers {
			if strings.EqualFold(existingKey, key) {
				m.err = errors.New("header already exists")
				return nil
			}
		}

		request.Headers[key] = m.headerDraftValue

		m.saving = true
		m.err = nil

		return saveRequest(m.store, request)

	case "left":
		if m.headerCursor > 0 {
			m.headerCursor--
		}

	case "right":
		if m.headerCursor < len([]rune(m.currentHeaderDraft())) {
			m.headerCursor++
		}

	case "home":
		m.headerCursor = 0

	case "end":
		m.headerCursor = len([]rune(m.currentHeaderDraft()))

	case "backspace":
		if m.headerCursor > 0 {
			runes := []rune(m.currentHeaderDraft())
			runes = append(runes[:m.headerCursor-1], runes[m.headerCursor:]...)

			m.setCurrentHeaderDraft(string(runes))
			m.headerCursor--
		}

	case "delete":
		runes := []rune(m.currentHeaderDraft())

		if m.headerCursor < len(runes) {
			runes = append(runes[:m.headerCursor], runes[m.headerCursor+1:]...)
			m.setCurrentHeaderDraft(string(runes))
		}

	case "enter":
		return nil

	default:
		if msg.Text != "" {
			runes := []rune(m.currentHeaderDraft())
			insert := []rune(msg.Text)

			updated := make([]rune, 0, len(runes)+len(insert))
			updated = append(updated, runes[:m.headerCursor]...)
			updated = append(updated, insert...)
			updated = append(updated, runes[m.headerCursor:]...)

			m.setCurrentHeaderDraft(string(updated))
			m.headerCursor += len(insert)
		}
	}

	return nil
}

func (m *model) currentHeaderDraft() string {
	if m.headerEditField == headerKeyField {
		return m.headerDraftKey
	}

	return m.headerDraftValue
}

func (m *model) setCurrentHeaderDraft(value string) {
	if m.headerEditField == headerKeyField {
		m.headerDraftKey = value
		return
	}

	m.headerDraftValue = value
}

func cloneHeaders(headers map[string]string) map[string]string {
	cloned := make(map[string]string, len(headers))

	for key, value := range headers {
		cloned[key] = value
	}

	return cloned
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

func saveRequest(
	store *storage.SQLite,
	request domain.Request,
) tea.Cmd {
	return func() tea.Msg {
		err := store.UpdateRequest(context.Background(), request)

		return requestSavedMsg{
			request: request,
			err:     err,
		}
	}
}

func (m *model) updateActiveRequest() {
	if len(m.requests) == 0 ||
		m.selectedRequestIndex < 0 ||
		m.selectedRequestIndex >= len(m.requests) {
		m.activeRequest = nil
		return
	}

	m.activeRequest = &m.requests[m.selectedRequestIndex]
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

	m.updateActiveRequest()
}

func (m *model) moveRequestField(delta int) {
	m.selectedField += requestField(delta)

	if m.selectedField < fieldMethod {
		m.selectedField = fieldMethod
	}

	if m.selectedField >= requestFieldCount {
		m.selectedField = requestFieldCount - 1
	}
}

func saveShortcut() (key, label string) {
	switch runtime.GOOS {
	case "darwin":
		return "super+s", "⌘+S"
	case "linux":
		return "super+s", "Super+S"
	case "windows":
		return "ctrl+s", "Ctrl+S"
	default:
		return "ctrl+s", "Ctrl+S"
	}
}
