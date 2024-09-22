package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Left  key.Binding
	Right  key.Binding
	Next   key.Binding
	Prev   key.Binding
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Help   key.Binding
	Quit   key.Binding
	Tab    key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Enter, k.Tab, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right},
		{k.Up, k.Down},
		{k.Next, k.Prev},
		{k.Enter, k.Tab},
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
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "ver articulo"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "mostrar ayuda"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "salir / volver atrás"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "cambiar categoria"),
	),
}
