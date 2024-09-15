package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
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

	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	fetching   bool
	quitting   bool
	fetchCmd   tea.Cmd
	spinner    spinner.Model
}

type keyMap struct {
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "j"),
		key.WithHelp("←/j", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "k"),
		key.WithHelp("→/k", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
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
		keys:        keys,
		help:        help.New(),
		inputStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
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
		switch {
		case key.Matches(msg, m.keys.Left):
			if !m.canGoBack || m.fetching {
				return m, nil
			}
			m.fetching = true
			m.fetchCmd = fetchEntries(m.currentPage - 1)
			m.lastKey = "←"
			return m, m.fetchCmd
		case key.Matches(msg, m.keys.Right):
			if !m.canContinue || m.fetching {
				return m, nil
			}
			m.fetching = true
			m.fetchCmd = fetchEntries(m.currentPage + 1)
			m.lastKey = "→"
			return m, tea.Batch(m.spinner.Tick, m.fetchCmd)
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) View() string {
	_, h, _ := term.GetSize(int(os.Stdout.Fd()))

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

		for index, entry := range m.entries {
			s += fmt.Sprintf("%d. %s\n", index+1, entry.Title)
		}
	}

	helpView := m.help.View(m.keys)

	remainingLines := h - strings.Count(s, "\n") - strings.Count(helpView, "\n") - 1

	s += strings.Repeat("\n", remainingLines)

	s += helpView

	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
