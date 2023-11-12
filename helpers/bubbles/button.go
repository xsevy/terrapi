package bubbles

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xsevy/terrapi/helpers"
)

var (
	buttonStyle        = lipgloss.NewStyle().PaddingTop(1).PaddingBottom(1).PaddingRight(4).PaddingLeft(4)
	bluredButtonStyle  = buttonStyle.Copy().Background(lipgloss.Color(helpers.Colors.Grey))
	focusedButtonStyle = buttonStyle.Copy().Background(lipgloss.Color(helpers.Colors.Purple))
)

type ButtonModel struct {
	text    string
	focused bool
}

func NewButtonModel(text string, focused bool) *ButtonModel {
	return &ButtonModel{
		text:    text,
		focused: focused,
	}
}

func (m *ButtonModel) Init() tea.Cmd {
	return nil
}

func (m *ButtonModel) View() string {
	if m.focused {
		return focusedButtonStyle.Render(m.text)
	}

	return bluredButtonStyle.Render(m.text)
}

func (m *ButtonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *ButtonModel) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *ButtonModel) Blur() {
	m.focused = false
}

func (m *ButtonModel) Value() string {
	return ""
}
