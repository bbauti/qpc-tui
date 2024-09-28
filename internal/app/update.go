package app

import (
	"net/http"
	"time"
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/glamour"

	"qpc-tui/internal/scraper"
)

const url = "https://quepensaschacabuco.com/"

type statusMsg int

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)

	if err != nil {
		return errMsg{err}
	}
	return statusMsg(res.StatusCode)
}

func fetchEntries(page int) tea.Cmd {
	return func() tea.Msg {
		entries, canContinue, canGoBack, err := scraper.ScrapePage(page)
		if err != nil {
			return errMsg{err}
		}
		return struct {
			entries     []scraper.Article
			canContinue bool
			canGoBack   bool
			page        int
		}{entries, canContinue, canGoBack, page}
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, checkServer, fetchEntries(m.CurrentPage))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case statusMsg:
		m.Status = int(msg)
		if m.Status == 200 && m.FetchCmd == nil {
			m.FetchCmd = fetchEntries(m.CurrentPage)
			return m, tea.Batch(m.Spinner.Tick, m.FetchCmd)
		}
		return m, nil

	case errMsg:
		m.Err = msg
		m.Fetching = false
		m.IsFirstFetch = false
		m.FetchCmd = nil
		return m, tea.Quit

	case struct {
		entries     []scraper.Article
		canContinue bool
		canGoBack   bool
		page        int
	}:
		m.Entries = msg.entries
		m.CanContinue = msg.canContinue
		m.CanGoBack = msg.canGoBack
		m.CurrentPage = msg.page
		m.Fetching = false
		m.IsFirstFetch = false
		m.FetchCmd = nil

		if (m.CanGoBack) {
			m.Keys.Left.Enabled = true
		} else {
			m.Keys.Left.Enabled = false
		}
		if (m.CanContinue) {
			m.Keys.Right.Enabled = true
		} else {
			m.Keys.Right.Enabled = false
		}

		sort.Slice(m.Entries, func(i, j int) bool {
			return m.Entries[i].Date > m.Entries[j].Date
		})

		items := make([]list.Item, len(m.Entries))
		for i, entry := range m.Entries {
			items[i] = item{title: entry.Title, desc: entry.Date}
		}
		m.List.SetItems(items)

		m.List.SetShowPagination(true)
		m.List.ResetSelected()
		m.List.ResetFilter()

		return m, tea.Batch(cmd, m.Spinner.Tick)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.List.SetSize(msg.Width, msg.Height-8)
		m.Viewport.Width = msg.Width
		m.Viewport.Height = msg.Height - 8

		cmd = m.List.NewStatusMessage(fmt.Sprintf("Window resized to %dx%d", msg.Width, msg.Height))

		return m, cmd

	case tea.KeyMsg:

		switch {
		case key.Matches(msg, m.Keys.Up.Binding) && m.Keys.Up.Enabled:
			m.Viewport.LineUp(1)
		case key.Matches(msg, m.Keys.Down.Binding) && m.Keys.Down.Enabled:
			m.Viewport.LineDown(1)
		case key.Matches(msg, m.Keys.Left.Binding) && m.Keys.Left.Enabled:
			if m.Fetching {
				return m, nil
			}
			m.Fetching = true
			m.FetchCmd = fetchEntries(m.CurrentPage - 1)
			m.LastKey = "←"
			return m, m.FetchCmd
		case key.Matches(msg, m.Keys.Right.Binding) && m.Keys.Right.Enabled:
			if m.Fetching {
				return m, nil
			}
			m.Fetching = true
			m.FetchCmd = fetchEntries(m.CurrentPage + 1)
			m.LastKey = "→"
			return m, tea.Batch(m.Spinner.Tick, m.FetchCmd)
		case key.Matches(msg, m.Keys.Help.Binding) && m.Keys.Help.Enabled:
			m.Help.ShowAll = !m.Help.ShowAll
			return m, nil
		case key.Matches(msg, m.Keys.Tab.Binding) && m.Keys.Tab.Enabled:
			m.CurrentCategory = (m.CurrentCategory + 1) % 4
			return m, nil
		case key.Matches(msg, m.Keys.Enter.Binding) && m.Keys.Enter.Enabled:
			selectedItem := m.List.SelectedItem().(item)
			for _, entry := range m.Entries {
				if entry.Title == selectedItem.Title() && entry.Date == selectedItem.Description() {
					m.SelectedEntry = &entry
					m.Keys.Quit.SetHelp("q", "volver atrás ")
					m.Keys.Up.SetHelp("↑", "subir ")
					m.Keys.Down.SetHelp("↓", "bajar ")
					m.Keys.Enter.Enabled = false
					m.Keys.Tab.Enabled = false
					m.Keys.Left.Enabled = false
					m.Keys.Right.Enabled = false

					r, _ := glamour.NewTermRenderer(
						glamour.WithAutoStyle(),
						glamour.WithWordWrap(m.Width-8),
					)

					bodyRendered, err := r.Render(m.SelectedEntry.Body)
					if err != nil {
						m.Err = err
						return m, tea.Quit
					}

					m.Viewport.SetContent(bodyRendered)
					m.Viewport.GotoTop()
					break
				}
			}
			return m, nil
		case key.Matches(msg, m.Keys.Quit.Binding) && m.Keys.Quit.Enabled:
			if m.SelectedEntry != nil {
				m.SelectedEntry = nil
				m.Keys.Quit.SetHelp("q", "salir")
				m.Keys.Up.SetHelp("↑", "articulo anterior ")
				m.Keys.Down.SetHelp("↓", "siguiente articulo ")
				m.Keys.Enter.Enabled = true
				m.Keys.Tab.Enabled = true
				m.Keys.Left.Enabled = true
				m.Keys.Right.Enabled = true
				return m, nil
			}
			m.Quitting = true
			return m, tea.Quit
		}
	}

	var listCmd tea.Cmd
	m.List, listCmd = m.List.Update(msg)
	return m, tea.Batch(cmd, listCmd)
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }