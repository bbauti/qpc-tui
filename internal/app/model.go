package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"

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

	CurrentCategory string
	SelectedEntry		*scraper.Article

	Keys         ui.KeyMap
	Help         help.Model
	InputStyle   lipgloss.Style
	LastKey      string
	Fetching     bool
	IsFirstFetch bool
	Quitting     bool
	FetchCmd     tea.Cmd
	Spinner      spinner.Model
	List         list.Model
	Viewport     viewport.Model

	renderer *lipgloss.Renderer
}

func InitialModel(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// The pty is the pseudo terminal that is created when the program starts,
	// it is used to get the size of the terminal.
	pty, _, _ := s.Pty()
	width, height := pty.Window.Width, pty.Window.Height

	// Since we use Wish, we need to use the MakeRenderer function to create the renderer,
	// since the one that lipgloss provides is not compatible with Wish.
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
		Width:     width,
		Height:    height,
		Bg:        bg,
		TxtStyle:  txtStyle,
		QuitStyle: quitStyle,

		CurrentCategory: "",
		SelectedEntry: nil,

		CurrentPage: 1,
		Spinner:     sp,
		Fetching:    true,
		IsFirstFetch: true,
		Keys:        ui.Keys,
		Help:        help.New(),
		InputStyle:  renderer.NewStyle().Foreground(lipgloss.Color("#FF75B7")),

		Viewport: viewport.New(width, height-8),

		renderer: renderer,
	}

	// To make the list work correctly with our custom renderer we need to use a custom
	// delegate and modify some styles.
	listItems := []list.Item{}
	l := list.New(listItems, NewCustomDelegate(renderer, m), width, height-6)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.Styles.PaginationStyle = renderer.NewStyle().PaddingLeft(2)
	l.Styles.ActivePaginationDot = renderer.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"}).
		SetString("•")
	l.Styles.InactivePaginationDot = renderer.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}).
		SetString("•")

	m.List = l

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}