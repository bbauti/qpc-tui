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
func (d customDelegate) Spacing() int                              { return 0 }
func (d customDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d customDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := d.renderer.NewStyle().
		Foreground(lipgloss.Color("205")).
		Render

	if index == m.Index() {
		fn = d.renderer.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("205")).
			Render
	}

	fmt.Fprint(w, fn(str))
}