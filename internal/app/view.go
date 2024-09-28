package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/glamour"
	"github.com/muesli/reflow/wordwrap"

	"qpc-tui/internal/scraper"
)

func (m Model) View() string {
	if m.Err != nil {
			return fmt.Sprintf("\nOcurrió un error: %v\n\n", m.Err)
	}

	navigationMenuItems := []string{
			"Todas",
			"Policiales",
			"Sociedad",
			"Automotores",
	}

	var tabItems []string
	for i, item := range navigationMenuItems {
			if m.SelectedEntry != nil {
					tabItems = append(tabItems, m.renderer.NewStyle().Foreground(lipgloss.Color("8")).Render(item))
			} else if i == m.CurrentCategory {
					tabItems = append(tabItems, m.renderer.NewStyle().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("0")).Render(item))
			} else {
					tabItems = append(tabItems, m.renderer.NewStyle().Render(item))
			}
	}

	title := m.renderer.NewStyle().
			Width(35).
			Foreground(lipgloss.Color("8")).
			Render(fmt.Sprintf("Chacabuco en Red TUI - Page %d", m.CurrentPage))

	tabStyle := m.renderer.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("8")).
			PaddingLeft(1).
			PaddingRight(1)

	styledTabs := make([]string, len(tabItems))
	for i, tab := range tabItems {
			styledTabs[i] = tabStyle.Render(tab)
	}

	navigationMenuContent := lipgloss.JoinHorizontal(lipgloss.Center, styledTabs...)

	titleAndNavigation := lipgloss.JoinHorizontal(
			lipgloss.Center,
			title,
			navigationMenuContent,
	)

	titleAndNavigation = m.renderer.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("8")).
			Margin(1, 2).
			Width(m.Width-4).
			Align(lipgloss.Center).
			Render(titleAndNavigation)

	var content string
	if m.Fetching {
			content = m.Spinner.View() + " Obteniendo entradas..."
	} else if m.Quitting {
			content = "Bye!"
	} else if m.SelectedEntry != nil {
		wrappedBody := wordwrap.String(m.SelectedEntry.Body, m.Width-8)

    bodyRendered, err := glamour.Render(wrappedBody, "dark")
    if err != nil {
        content += fmt.Sprintf("Error rendering body: %v\n", err)
    } else {
        content += bodyRendered
    }
    content += m.renderer.NewStyle().Width(m.Width-4).MarginTop(m.Height-8).Align(lipgloss.Center).Foreground(lipgloss.Color("8")).Render(fmt.Sprintf(m.SelectedEntry.Link))
	} else if m.Status > 0 && len(m.Entries) > 0 {
			filteredEntries := filterEntriesByCategory(m.Entries, m.CurrentCategory)
			m.List.SetItems(entriesToListItems(filteredEntries))
			content = m.List.View()
	} else {
			content = m.Spinner.View() + " Obteniendo entradas..."
	}

	helpView := m.renderer.NewStyle().MarginLeft(1).Render(m.Help.View(m.Keys))

	contentHeight := m.Height-8
	if m.Help.ShowAll {
		contentHeight += 1
	}

	contentLines := strings.Split(content, "\n")
	if len(contentLines) > contentHeight {
		content = strings.Join(contentLines[:contentHeight], "\n")
	}

	for len(strings.Split(content, "\n")) < contentHeight {
		content += "\n"
	}

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

	if m.SelectedEntry != nil {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleAndNavigation,
			content,
			"\n",
			helpView,
		)
	}

	if m.Help.ShowAll {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleAndNavigation,
			content,
			helpView,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleAndNavigation,
		content,
		"\n",
		helpView,
	)
}

func filterEntriesByCategory(entries []scraper.Article, currentCategory int) []scraper.Article {
	if currentCategory == 0 {
			return entries
	}

	var filteredEntries []scraper.Article
	for _, entry := range entries {
			if entry.CategoryId == getCategoryIdByCategory(currentCategory) {
					filteredEntries = append(filteredEntries, entry)
			}
	}
	return filteredEntries
}

func getCategoryIdByCategory(category int) int {
	switch category {
	case 1:
			return 8 // Policiales
	case 2:
			return 48 // Sociedad
	case 3:
			return 75 // Automotores
	default:
			return 0 // All
	}
}

func entriesToListItems(entries []scraper.Article) []list.Item {
	items := make([]list.Item, len(entries))
	for i, entry := range entries {
			items[i] = item{title: entry.Title, desc: entry.Date}
	}
	return items
}