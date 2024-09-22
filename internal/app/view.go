package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("\nOcurrió un error: %v\n\n", m.Err)
	}

	headerStyle := m.Renderer.NewStyle().
		MarginTop(1).
		Width(m.Width).
		Align(lipgloss.Center).
		Bold(true).
		Foreground(lipgloss.Color("205"))

	header := headerStyle.Render(fmt.Sprintf("Chacabuco en Red TUI - Page %d", m.CurrentPage))

	var content string
	if m.Fetching {
		content = m.Spinner.View() + " Obteniendo entradas..."
	} else if m.Quitting {
		content = "Bye!"
	} else if m.Status > 0 && len(m.Entries) > 0 {
		content = m.List.View()
	} else {
		content = m.Spinner.View() + " Obteniendo entradas..."
	}

	// add left margin to the help view
	helpView := m.Renderer.NewStyle().MarginLeft(1).Render(m.Help.View(m.Keys))

	// Calculate available height for content
	contentHeight := m.Height - 10 // Subtract space for header, help, and margins
	if m.Help.ShowAll {
		contentHeight -= strings.Count(helpView, "\n") + 1
	}

	// Ensure the content doesn't exceed the available height
	contentLines := strings.Split(content, "\n")
	if len(contentLines) > contentHeight {
		content = strings.Join(contentLines[:contentHeight], "\n")
	}

	// Pad the content to ensure consistent height
	for len(strings.Split(content, "\n")) < contentHeight {
		content += "\n"
	}

	navigationStyles := m.Renderer.NewStyle().
		MarginLeft(2)


	navigationText := ""
	if (m.CanGoBack && m.CanContinue) {
		navigationText = "← →"
	} else if (m.CanGoBack) {
		navigationText = "←"
	} else if (m.CanContinue) {
		navigationText = "→"
	}

	navigation := navigationStyles.Render(navigationText)

	if (m.IsFirstFetch) {
		loadingContent := lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			"Obteniendo entradas...",
		)
		return lipgloss.JoinVertical(
			lipgloss.Left,
			loadingContent,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"\n",
		content,
		"\n",
		navigation,
		"\n",
		helpView,
	)
}