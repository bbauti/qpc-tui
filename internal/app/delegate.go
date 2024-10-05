/*
	Custom delegate for the List bubble
	We need a custom delegate to render the list with the model renderer, since the
	list bubble uses lipgloss to render the list the styles weren't applied correctly
	in the ssh server.
*/

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
	model Model
}

func NewCustomDelegate(renderer *lipgloss.Renderer, model Model) list.ItemDelegate {
	return &customDelegate{renderer: renderer, model: model}
}

func (d customDelegate) Height() int                               { return 2 }
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

	indexStr := fmt.Sprintf("%d. ", index+1)
	titleStr := i.title
	var subtitle string
	if d.model.Entries != nil {
		for _, entry := range d.model.Entries {
			if entry.Title == i.title && entry.Date == i.desc {
				subtitle = fmt.Sprintf("%s | %s", entry.Category, entry.Date)
				break
			}
		}
	}
	if index < 9 {
		indexStr = " " + indexStr
	}

	indexStyle := d.renderer.NewStyle().
		MarginLeft(4)

	titleStyle := d.renderer.NewStyle()

	if index == m.Index() {
		titleStyle = titleStyle.
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("102"))
	}

	subtitleStyle := d.renderer.NewStyle().
		Foreground(lipgloss.Color("8")).
		MarginLeft(8)

	fmt.Fprint(w, indexStyle.Render(indexStr))
	fmt.Fprint(w, titleStyle.Render(titleStr))
	if subtitle != "" {
		fmt.Fprint(w, "\n"+subtitleStyle.Render(subtitle))
	}
}