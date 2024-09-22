package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/list"

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
					tabItems = append(tabItems, lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(item))
			} else if i == m.CurrentCategory {
					tabItems = append(tabItems, lipgloss.NewStyle().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("0")).Render(item))
			} else {
					tabItems = append(tabItems, lipgloss.NewStyle().Render(item))
			}
	}

	title := lipgloss.NewStyle().
			Width(35).
			Foreground(lipgloss.Color("8")).
			Render(fmt.Sprintf("Chacabuco en Red TUI - Page %d", m.CurrentPage))

	// Create a style for each tab item with vertical borders
	tabStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("8")).
			PaddingLeft(1).
			PaddingRight(1)

	// Apply the style to each tab item
	styledTabs := make([]string, len(tabItems))
	for i, tab := range tabItems {
			styledTabs[i] = tabStyle.Render(tab)
	}

	// Join the styled tabs
	navigationMenuContent := lipgloss.JoinHorizontal(lipgloss.Center, styledTabs...)

	// Join title and navigationMenuContent horizontally
	titleAndNavigation := lipgloss.JoinHorizontal(
			lipgloss.Center,
			title,
			navigationMenuContent,
	)

	// Add a border around the entire titleAndNavigation
	titleAndNavigation = lipgloss.NewStyle().
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
			content = fmt.Sprintf("Title: %v\n", m.SelectedEntry.Title)
			content += fmt.Sprintf("Date: %v\n", m.SelectedEntry.Date)
			content += fmt.Sprintf("Category: %v\n", m.SelectedEntry.Category)
			content += fmt.Sprintf("CategoryId: %v\n", m.SelectedEntry.CategoryId)
			content += fmt.Sprintf("Body: %v\n", m.SelectedEntry.Body)
			content += fmt.Sprintf("Link: %v\n", m.SelectedEntry.Link)
	} else if m.Status > 0 && len(m.Entries) > 0 {
			filteredEntries := filterEntriesByCategory(m.Entries, m.CurrentCategory)
			m.List.SetItems(entriesToListItems(filteredEntries))
			content = m.List.View()
	} else {
			content = m.Spinner.View() + " Obteniendo entradas..."
	}

	// add left margin to the help view
	helpView := lipgloss.NewStyle().MarginLeft(1).Render(m.Help.View(m.Keys))

	// Calculate available height for content
	contentHeight := m.Height - 11 // Subtract space for header, help, and margins
	if m.Help.ShowAll {
		contentHeight -= strings.Count(helpView, "\n") + 1
	}
	if (m.SelectedEntry != nil) {
		contentHeight += strings.Count(content, "\n") + 2
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

	navigationStyles := lipgloss.NewStyle().
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

	if m.SelectedEntry != nil {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleAndNavigation,
			content,
			"\n",
			helpView,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleAndNavigation,
		content,
		"\n",
		navigation,
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