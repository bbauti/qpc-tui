package app

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func (m Model) View() string {
	_, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		h = 24
	}

	if m.Err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.Err)
	}

	s := fmt.Sprintf("Chacabuco en Red TUI. Page %d\n\n", m.CurrentPage)

	if m.Fetching {
		s += m.Spinner.View() + " Loading...\n"
	} else if m.Quitting {
		s += "Bye!\n"
	} else if m.Status > 0 && len(m.Entries) > 0 {
		s += fmt.Sprintf("Status: %d\n", m.Status)
		s += fmt.Sprintf("Can continue: %v\n", m.CanContinue)
		s += fmt.Sprintf("Can go back: %v\n", m.CanGoBack)
		s += fmt.Sprintf("Current page: %d\n", m.CurrentPage)
		s += fmt.Sprintf("Entries: %d\n", len(m.Entries))

		for index, entry := range m.Entries {
			s += fmt.Sprintf("%d. %s\n", index+1, entry.Title)
		}
	}

	helpView := m.Help.View(m.Keys)

	remainingLines := h - strings.Count(s, "\n") - strings.Count(helpView, "\n") - 1
	if remainingLines > 0 {
		s += strings.Repeat("\n", remainingLines)
	}

	s += helpView

	return s
}