package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"

	"qpc-tui/internal/scraper"
	"qpc-tui/internal/ui"
)

type Model struct {
	Term      string
	Profile   string
	Width     int
	Height    int
	Bg        string
	TxtStyle  lipgloss.Style
	QuitStyle lipgloss.Style

	Status      int
	CurrentPage int
	Entries     []scraper.Article
	CanContinue bool
	CanGoBack   bool
	Err         error

	Keys       ui.KeyMap
	Help       help.Model
	InputStyle lipgloss.Style
	LastKey    string
	Fetching   bool
	Quitting   bool
	FetchCmd   tea.Cmd
	Spinner    spinner.Model
}

func InitialModel(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()

	renderer := bubbletea.MakeRenderer(s)
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	bg := "light"
	if renderer.HasDarkBackground() {
		bg = "dark"
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = renderer.NewStyle().Foreground(lipgloss.Color("205"))

	m := Model{
		Term:      pty.Term,
		Profile:   renderer.ColorProfile().Name(),
		Width:     pty.Window.Width,
		Height:    pty.Window.Height,
		Bg:        bg,
		TxtStyle:  txtStyle,
		QuitStyle: quitStyle,

		CurrentPage: 1,
		Spinner:     sp,
		Fetching:    true,
		Keys:        ui.Keys,
		Help:        help.New(),
		InputStyle:  renderer.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}