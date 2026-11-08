package tui

import (
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing...")
	}

	if m.width < 60 || m.height < 15 {
		return tea.NewView("Terminal is too small. Resize to at least 60x15.")
	}

	// Calculate available dimensions.
	availableWidth := m.width - appStyle.GetHorizontalFrameSize() - 4
	availableHeight := m.height - appStyle.GetVerticalFrameSize() - 4

	// Allocate 25% to the sidebar and 75% to the main area.
	sidebarWidth := availableWidth / 4
	mainWidth := availableWidth - sidebarWidth

	requestHeight := availableHeight / 2

	// Build sidebar content.
	items := []string{}

	if m.sidebarMode == sidebarCollections {
		items = append(items, titleStyle.Render("Collections"), "")

		for i, collection := range m.collections {
			if i == m.selectedCollection {
				items = append(items,
					selectedCollectionStyle.Render("> "+collection.Name),
				)
			} else {
				items = append(items,
					collectionStyle.Render("  "+collection.Name),
				)
			}
		}

		if len(m.collections) == 0 {
			items = append(items, mutedStyle.Render("No collections yet"))
		}
	} else {
		items = append(items, titleStyle.Render("Requests"), "")

		for i, request := range m.requests {
			if i == m.selectedRequestIndex {
				items = append(items,
					selectedCollectionStyle.Render("> "+request.Name),
				)
			} else {
				items = append(items,
					collectionStyle.Render("  "+request.Name),
				)
			}
		}

		if m.loading {
			items = append(items, mutedStyle.Render("Loading requests..."))
		} else if len(m.requests) == 0 {
			items = append(items, mutedStyle.Render("No requests yet"))
		}
	}

	if m.err != nil {
		items = append(items, mutedStyle.Render(m.err.Error()))
	}

	// Keep panel content within its allocated area.
	fitContent := func(content string, width, height int) string {
		return lipgloss.NewStyle().
			Width(width).
			MaxHeight(height).
			Render(content)
	}

	// Sidebar panel.
	sidebarBorder := "240"
	if m.focusedPanel == focusSidebar {
		sidebarBorder = "42"
	}

	sidebarContentWidth := sidebarWidth - sidebarStyle.GetHorizontalFrameSize()
	sidebarContentHeight := availableHeight - sidebarStyle.GetVerticalFrameSize()

	sidebar := sidebarStyle.
		Width(sidebarContentWidth).
		Height(availableHeight).
		BorderForeground(lipgloss.Color(sidebarBorder)).
		Render(fitContent(
			lipgloss.JoinVertical(lipgloss.Left, items...),
			sidebarContentWidth,
			sidebarContentHeight,
		))

	panelContentWidth := mainWidth - panelStyle.GetHorizontalFrameSize()

	// Highlight the selected request field when the Request panel is focused.
	renderField := func(field requestField, content string) string {
		if m.focusedPanel != focusRequest || m.selectedField != field {
			return content
		}

		return selectedFieldStyle.
			Width(panelContentWidth - selectedFieldStyle.GetHorizontalFrameSize()).
			Render(content)
	}

	// Render the URL draft with a visible cursor.
	renderEditableURL := func() string {
		runes := []rune(m.urlDraft)

		before := string(runes[:m.urlCursor])
		after := string(runes[m.urlCursor:])

		cursor := " "
		if m.urlCursor < len(runes) {
			cursor = string(runes[m.urlCursor])
		}

		return before + urlCursorStyle.Render(cursor) + after
	}

	// Build request content.
	requestContent := []string{
		titleStyle.Render("Request"),
		"",
	}

	if m.activeRequest == nil {
		requestContent = append(
			requestContent,
			mutedStyle.Render("Select a request to get started."),
		)
	} else {
		request := m.activeRequest

		// Method field.
		methodContent := strings.Join([]string{
			mutedStyle.Render("Method"),
			renderMethod(request.Method),
		}, "\n")

		requestContent = append(
			requestContent,
			renderField(fieldMethod, methodContent),
			"",
		)

		// URL field.
		urlValue := request.URL

		if m.editingURL {
			urlValue = renderEditableURL()
		}

		urlContent := strings.Join([]string{
			mutedStyle.Render("URL"),
			urlValue,
		}, "\n")

		requestContent = append(
			requestContent,
			renderField(fieldURL, urlContent),
			"",
		)

		// Headers field.
		headerLines := []string{
			mutedStyle.Render("Headers"),
		}

		if len(request.Headers) == 0 {
			headerLines = append(
				headerLines,
				mutedStyle.Render("No headers"),
			)
		} else {
			keys := make([]string, 0, len(request.Headers))

			for key := range request.Headers {
				keys = append(keys, key)
			}

			sort.Strings(keys)

			for _, key := range keys {
				headerLines = append(
					headerLines,
					key+": "+request.Headers[key],
				)
			}
		}

		requestContent = append(
			requestContent,
			renderField(fieldHeaders, strings.Join(headerLines, "\n")),
			"",
		)

		// Body field.
		body := request.Body

		if body == "" {
			body = mutedStyle.Render("No body")
		}

		bodyContent := strings.Join([]string{
			mutedStyle.Render("Body"),
			body,
		}, "\n")

		requestContent = append(
			requestContent,
			renderField(fieldBody, bodyContent),
		)
	}

	// Request panel.
	requestBorder := lipgloss.Color("240")
	if m.focusedPanel == focusRequest {
		requestBorder = lipgloss.Color("42")
	}

	requestContentHeight := requestHeight - panelStyle.GetVerticalFrameSize()

	requestPanel := panelStyle.
		Width(panelContentWidth).
		Height(requestHeight).
		BorderForeground(requestBorder).
		Render(fitContent(
			lipgloss.JoinVertical(lipgloss.Left, requestContent...),
			panelContentWidth,
			requestContentHeight,
		))

	// Response panel.
	responseBorder := lipgloss.Color("240")
	if m.focusedPanel == focusResponse {
		responseBorder = lipgloss.Color("42")
	}

	responseHeight := availableHeight - requestHeight
	responseContentHeight := responseHeight - panelStyle.GetVerticalFrameSize()

	responsePanel := panelStyle.
		Width(panelContentWidth).
		Height(responseHeight).
		BorderForeground(responseBorder).
		Render(fitContent(
			lipgloss.JoinVertical(
				lipgloss.Left,
				titleStyle.Render("Response"),
				"",
				mutedStyle.Render("Response data will appear here."),
			),
			panelContentWidth,
			responseContentHeight,
		))
	// Combine main panels.
	main := lipgloss.JoinVertical(
		lipgloss.Left,
		requestPanel,
		responsePanel,
	)

	header := titleStyle.Render("NetPilot") +
		mutedStyle.Render("  HTTP client")

	var bindings []string

	switch {
	case m.saving:
		bindings = []string{
			mutedStyle.Render("Saving request..."),
		}

	case m.editingURL:
		_, saveLabel := saveShortcut()

		bindings = []string{
			renderKeyBinding("←/→", "Move cursor"),
			renderKeyBinding(saveLabel, "Save"),
			renderKeyBinding("esc", "Cancel"),
		}

	default:
		bindings = []string{
			renderKeyBinding("j/k", "Navigate"),
			renderKeyBinding("h/l", "Switch panel"),
		}

		if m.focusedPanel == focusRequest &&
			m.selectedField == fieldURL &&
			m.activeRequest != nil {
			bindings = append(
				bindings,
				renderKeyBinding("enter", "Edit URL"),
			)
		} else {
			bindings = append(
				bindings,
				renderKeyBinding("enter", "Select"),
			)
		}

		if m.sidebarMode == sidebarRequests {
			bindings = append(bindings, renderKeyBinding("esc", "Back"))
		}

		bindings = append(bindings, renderKeyBinding("q", "Quit"))
	}

	footer := strings.Join(
		bindings,
		mutedStyle.Render("  |  "),
	)

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

	return tea.View{
		Content:   appStyle.Render(content),
		AltScreen: true,
	}
}
