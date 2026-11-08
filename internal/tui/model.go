package tui

import (
	"sort"

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

type requestField int

const (
	fieldMethod requestField = iota
	fieldURL
	fieldHeaders
	fieldBody
)

type headerEditField int

const (
	headerKeyField headerEditField = iota
	headerValueField
)

const requestFieldCount = fieldBody + 1

type model struct {
	store *storage.SQLite

	collections        []domain.Collection
	selectedCollection int

	requests             []domain.Request
	selectedRequestIndex int
	activeRequest        *domain.Request

	selectedField requestField

	editingURL bool
	urlDraft   string
	urlCursor  int
	saving     bool

	editingHeaders bool
	editingHeader  bool

	selectedHeader int

	headerDraftKey    string
	headerDraftValue  string
	headerOriginalKey string

	headerEditField headerEditField
	headerCursor    int

	sidebarMode  sidebarMode
	focusedPanel focus

	width   int
	height  int
	loading bool
	err     error
}

func initialModel(store *storage.SQLite) *model {
	return &model{
		store:         store,
		focusedPanel:  focusSidebar,
		sidebarMode:   sidebarCollections,
		selectedField: fieldMethod,
	}
}

func sortedHeaderKeys(headers map[string]string) []string {
	keys := make([]string, 0, len(headers))

	for key := range headers {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
