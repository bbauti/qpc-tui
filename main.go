package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"qpc-tui/scraper"
)

const url = "https://quepensaschacabuco.com/"

type model struct {
	status      int
	currentPage int
	entries     []scraper.Article
	canContinue bool
	canGoBack   bool
	err         error

	fetching bool
	quitting bool
	fetchCmd tea.Cmd
	spinner  spinner.Model
}

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

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		currentPage: 1,
		spinner:     s,
		fetching:    true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, checkServer, fetchEntries(m.currentPage))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case statusMsg:
		m.status = int(msg)
		if m.status == 200 && m.fetchCmd == nil {
			m.fetchCmd = fetchEntries(m.currentPage)
			return m, tea.Batch(m.spinner.Tick, m.fetchCmd)
		}
		return m, nil

	case errMsg:
		m.err = msg
		m.fetching = false
		m.fetchCmd = nil
		return m, tea.Quit

	case struct {
		entries     []scraper.Article
		canContinue bool
		canGoBack   bool
		page        int
	}:
		m.entries = msg.entries
		m.canContinue = msg.canContinue
		m.canGoBack = msg.canGoBack
		m.currentPage = msg.page
		m.fetching = false
		m.fetchCmd = nil
		return m, tea.Batch(cmd, m.spinner.Tick)

	case tea.KeyMsg:
		switch msg.String() {
		case "right", "l":
			if !m.canContinue || m.fetching {
				return m, nil
			}
			m.fetching = true
			m.fetchCmd = fetchEntries(m.currentPage + 1)
			return m, tea.Batch(m.spinner.Tick, m.fetchCmd)
		case "left", "k":
			if !m.canGoBack || m.fetching {
				return m, nil
			}
			m.fetching = true
			m.fetchCmd = fetchEntries(m.currentPage - 1)
			return m, m.fetchCmd
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	s := fmt.Sprintf("Chacabuco en Red TUI. Page %d\n\n", m.currentPage)

	if m.fetching {
		s += m.spinner.View() + " Loading...\n"
	} else if m.quitting {
		s += "Bye!\n"
	} else if m.status > 0 && len(m.entries) > 0 {
		s += fmt.Sprintf("Status: %d\n", m.status)
		s += fmt.Sprintf("Can continue: %v\n", m.canContinue)
		s += fmt.Sprintf("Can go back: %v\n", m.canGoBack)
		s += fmt.Sprintf("Current page: %d\n", m.currentPage)
		s += fmt.Sprintf("Entries: %d\n", len(m.entries))
	}

	// help section
	s += "\nControls:\n"
	if m.canGoBack {
		s += "  <- left\n"
	}
	if m.canContinue {
		s += "  -> right\n"
	}
	s += "  q quit\n"

	return "\n" + s + "\n\n"
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
