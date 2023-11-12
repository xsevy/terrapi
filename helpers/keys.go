package helpers

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Help     key.Binding
	Quit     key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up "),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down "),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help "),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit "),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "confirm "),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back "),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next "),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev "),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Help, k.Quit},
		{k.Enter, k.Escape},
		{k.Tab, k.ShiftTab},
	}
}
