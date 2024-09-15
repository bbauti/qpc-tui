package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

func (m model) Init() tea.Cmd {
	return checkServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		m.status = int(msg)
		if m.status == 200 {
			entries, canContinue, canGoBack, err := scraper.ScrapePage(m.currentPage)
			if err != nil {
				fmt.Println("Error scraping page:", err)
			} else {
				m.entries = entries
				m.canContinue = canContinue
				m.canGoBack = canGoBack
			}
		}
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	s := fmt.Sprintf("Checking %s...\n\n", url)

	if m.status > 0 && len(m.entries) > 0 {
		s += fmt.Sprintf("Status: %d\n", m.status)
		s += fmt.Sprintf("Can continue: %v\n", m.canContinue)
		s += fmt.Sprintf("Can go back: %v\n", m.canGoBack)
		s += fmt.Sprintf("Current page: %d\n", m.currentPage)
		s += fmt.Sprintf("Entries: %d\n", len(m.entries))
	} else {
		s += "Loading..."
	}

	return "\n" + s + "\n\n"
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
