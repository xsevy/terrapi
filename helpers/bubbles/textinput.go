package bubbles

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xsevy/terrapi/styles"
)

type TextInputModel struct {
	t     textinput.Model
	title string
}

func NewTextInput(title string, placeholder string, charlimit int) *TextInputModel {
	ti := TextInputModel{
		title: title,
		t:     textinput.New(),
	}
	ti.t.Placeholder = placeholder
	ti.t.CharLimit = charlimit
	ti.t.Focus()
	return &ti
}

func (m *TextInputModel) Init() tea.Cmd {
	return nil
}

func (m *TextInputModel) Focused() bool {
	return m.t.Focused()
}

func (m *TextInputModel) Value() string {
	return m.t.Value()
}

func (m *TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.t, cmd = m.t.Update(msg)
	return m, cmd
}

func (m *TextInputModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, styles.GetFocusedTitle(m.title, m.t.Focused()), m.t.View())
}

func (m *TextInputModel) Focus() tea.Cmd {
	m.t.Focus()
	return nil
}

func (m *TextInputModel) Blur() {
	m.t.Blur()
}
