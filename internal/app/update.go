package app

import (
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

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
		m.FetchCmd = nil
		return m, tea.Batch(cmd, m.Spinner.Tick)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Left):
			if !m.CanGoBack || m.Fetching {
				return m, nil
			}
			m.Fetching = true
			m.FetchCmd = fetchEntries(m.CurrentPage - 1)
			m.LastKey = "←"
			return m, m.FetchCmd
		case key.Matches(msg, m.Keys.Right):
			if !m.CanContinue || m.Fetching {
				return m, nil
			}
			m.Fetching = true
			m.FetchCmd = fetchEntries(m.CurrentPage + 1)
			m.LastKey = "→"
			return m, tea.Batch(m.Spinner.Tick, m.FetchCmd)
		case key.Matches(msg, m.Keys.Help):
			m.Help.ShowAll = !m.Help.ShowAll
		case key.Matches(msg, m.Keys.Quit):
			m.Quitting = true
			return m, tea.Quit
		}
	}

	return m, cmd
}