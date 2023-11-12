package navigation

import tea "github.com/charmbracelet/bubbletea"

type Focusable interface {
	Blur()
	Focus() tea.Cmd
}

type FormField interface {
	tea.Model
	Focusable
	Value() string
}
