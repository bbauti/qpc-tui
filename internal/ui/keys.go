package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Left  key.Binding
	Right  key.Binding
	Next   key.Binding
	Prev   key.Binding
	Up     key.Binding
	Down   key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right},
		{k.Up, k.Down},
		{k.Next, k.Prev},
		{k.Help, k.Quit},
	}
}

var Keys = KeyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "j"),
		key.WithHelp("←/j", "pagina anterior"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "k"),
		key.WithHelp("→/k", "siguiente pagina"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "articulo anterior "),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "siguiente articulo "),
	),
	Next: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "bajar pagina "),
	),
	Prev: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "subir pagina "),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "mostrar ayuda"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "salir"),
	),
}