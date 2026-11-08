package tui

import (
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	if m.width < 60 || m.height < 15 {
		return "Terminal is too small. Resize to at least 60x15."
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
	sidebarBorder := lipgloss.Color("240")
	if m.focusedPanel == focusSidebar {
		sidebarBorder = lipgloss.Color("42")
	}

	sidebarContentWidth := sidebarWidth - sidebarStyle.GetHorizontalFrameSize()
	sidebarContentHeight := availableHeight - sidebarStyle.GetVerticalFrameSize()

	sidebar := sidebarStyle.
		Width(sidebarContentWidth).
		Height(sidebarContentHeight).
		BorderForeground(sidebarBorder).
		Render(fitContent(
			lipgloss.JoinVertical(lipgloss.Left, items...),
			sidebarContentWidth,
			sidebarContentHeight,
		))

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

		requestContent = append(
			requestContent,
			renderMethod(request.Method),
			"",
			mutedStyle.Render("URL"),
			request.URL,
			"",
			mutedStyle.Render("Headers"),
		)

		if len(request.Headers) == 0 {
			requestContent = append(
				requestContent,
				mutedStyle.Render("No headers"),
			)
		} else {
			keys := make([]string, 0, len(request.Headers))

			for key := range request.Headers {
				keys = append(keys, key)
			}

			sort.Strings(keys)

			for _, key := range keys {
				requestContent = append(
					requestContent,
					key+": "+request.Headers[key],
				)
			}
		}

		requestContent = append(
			requestContent,
			"",
			mutedStyle.Render("Body"),
		)

		if request.Body == "" {
			requestContent = append(
				requestContent,
				mutedStyle.Render("No body"),
			)
		} else {
			requestContent = append(requestContent, request.Body)
		}
	}

	// Request panel.
	requestBorder := lipgloss.Color("240")
	if m.focusedPanel == focusRequest {
		requestBorder = lipgloss.Color("42")
	}

	panelContentWidth := mainWidth - panelStyle.GetHorizontalFrameSize()
	requestContentHeight := requestHeight - panelStyle.GetVerticalFrameSize()

	requestPanel := panelStyle.
		Width(panelContentWidth).
		Height(requestContentHeight).
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

	responseContentHeight := lipgloss.Height(sidebar) -
		lipgloss.Height(requestPanel) -
		panelStyle.GetVerticalFrameSize()/2

	responsePanel := panelStyle.
		Width(panelContentWidth).
		Height(responseContentHeight).
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

	bindings := []string{
		renderKeyBinding("j/k", "Navigate"),
		renderKeyBinding("h/l", "Switch panel"),
		renderKeyBinding("enter", "Select"),
	}

	if m.sidebarMode == sidebarRequests {
		bindings = append(bindings, renderKeyBinding("esc", "Back"))
	}

	bindings = append(bindings, renderKeyBinding("q", "Quit"))

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

	return appStyle.Render(content)
}
