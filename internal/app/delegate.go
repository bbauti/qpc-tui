package app

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type customDelegate struct {
	renderer *lipgloss.Renderer
}

func NewCustomDelegate(renderer *lipgloss.Renderer) list.ItemDelegate {
	return &customDelegate{renderer: renderer}
}

func (d customDelegate) Height() int                               { return 1 }
func (d customDelegate) Spacing() int                              { return 1 }
func (d customDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

type Category int

const (
	Policiales Category = 8
	Sociedad Category = 48
	Automotores Category = 75
)


func (d customDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, i.title)
	if index < 9 {
		str = fmt.Sprintf(" %s", str)
	}

	fn := d.renderer.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginLeft(4).
		Render

	if index == m.Index() {
		fn = d.renderer.NewStyle().
			Foreground(lipgloss.Color("#1A1A1A")).
			MarginLeft(4).
			Background(lipgloss.Color("#888888")).
			Render
	}

	// You can use categoryId here if needed

	fmt.Fprint(w, fn(str))
}