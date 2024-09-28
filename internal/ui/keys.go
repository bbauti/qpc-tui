package ui

import "github.com/charmbracelet/bubbles/key"

type KeyBinding struct {
	key.Binding
	Enabled bool
}

type KeyMap struct {
	Left  KeyBinding
	Right KeyBinding
	Next  KeyBinding
	Prev  KeyBinding
	Up    KeyBinding
	Down  KeyBinding
	Enter KeyBinding
	Help  KeyBinding
	Quit  KeyBinding
	Tab   KeyBinding
}

func (k KeyMap) ShortHelp() []key.Binding {
	bindings := []key.Binding{}
	for _, kb := range []KeyBinding{k.Left, k.Right, k.Enter, k.Tab, k.Help, k.Quit} {
		if kb.Enabled {
			bindings = append(bindings, kb.Binding)
		}
	}
	return bindings
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.enabledBindings(k.Left, k.Right),
		k.enabledBindings(k.Up, k.Down),
		k.enabledBindings(k.Next, k.Prev),
		k.enabledBindings(k.Enter, k.Tab),
		k.enabledBindings(k.Help, k.Quit),
	}
}

func (k KeyMap) enabledBindings(bindings ...KeyBinding) []key.Binding {
	enabled := []key.Binding{}
	for _, b := range bindings {
		if b.Enabled {
			enabled = append(enabled, b.Binding)
		}
	}
	return enabled
}

func (k *KeyMap) EnableKey(keyName string) {
	k.setKeyState(keyName, true)
}

func (k *KeyMap) DisableKey(keyName string) {
	k.setKeyState(keyName, false)
}

func (k *KeyMap) setKeyState(keyName string, enabled bool) {
	switch keyName {
	case "Left":
		k.Left.Enabled = enabled
	case "Right":
		k.Right.Enabled = enabled
	case "Next":
		k.Next.Enabled = enabled
	case "Prev":
		k.Prev.Enabled = enabled
	case "Up":
		k.Up.Enabled = enabled
	case "Down":
		k.Down.Enabled = enabled
	case "Enter":
		k.Enter.Enabled = enabled
	case "Help":
		k.Help.Enabled = enabled
	case "Quit":
		k.Quit.Enabled = enabled
	case "Tab":
		k.Tab.Enabled = enabled
	}
}

var Keys = KeyMap{
	Left: KeyBinding{
		Binding: key.NewBinding(
			key.WithKeys("left", "j"),
			key.WithHelp("←/j", "pagina anterior"),
		),
		Enabled: true,
	},
	Right: KeyBinding{
		Binding: key.NewBinding(
			key.WithKeys("right", "k"),
			key.WithHelp("→/k", "siguiente pagina"),
		),
		Enabled: true,
	},
	Up: KeyBinding{
		Binding: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "articulo anterior "),
		),
		Enabled: true,
	},
	Down: KeyBinding{
		Binding: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "siguiente articulo "),
		),
		Enabled: true,
	},
	Next: KeyBinding{
		Binding: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "bajar pagina "),
		),
		Enabled: true,
	},
	Prev: KeyBinding{
		Binding: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "subir pagina "),
		),
		Enabled: true,
	},
	Enter: KeyBinding{
		Binding: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "ver articulo"),
		),
		Enabled: true,
	},
	Help: KeyBinding{
		Binding: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "mostrar ayuda"),
		),
		Enabled: true,
	},
	Quit: KeyBinding{
		Binding: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q", "salir / volver atrás"),
		),
		Enabled: true,
	},
	Tab: KeyBinding{
		Binding: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "cambiar categoria"),
		),
		Enabled: true,
	},
}
