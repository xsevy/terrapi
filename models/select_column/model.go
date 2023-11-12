package select_column

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xsevy/terrapi/helpers"
	"github.com/xsevy/terrapi/helpers/models"
	"github.com/xsevy/terrapi/models/select_column_choices"
	"github.com/xsevy/terrapi/styles"
)

type SelectColumnModel struct {
	choices *select_column_choices.SelectColumnChoicesModel
	keys    helpers.KeyMap
	models.ColumnModel
}

func NewSelectColumnModel(choices *select_column_choices.SelectColumnChoicesModel, focused bool) *SelectColumnModel {
	m := &SelectColumnModel{
		choices: choices,
		keys:    helpers.Keys,
	}
	m.SetFocused(focused)
	return m
}

func (m *SelectColumnModel) Init() tea.Cmd {
	return nil
}

func (m *SelectColumnModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var newSelectColumnChoices tea.Model

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up, m.keys.Down, m.keys.Enter, m.keys.Escape):
			newSelectColumnChoices, cmd = m.choices.Update(msg)
			m.choices = newSelectColumnChoices.(*select_column_choices.SelectColumnChoicesModel)
		}
	}
	return m, cmd
}

func (m *SelectColumnModel) View() string {
	var content string
	choicesView := m.choices.View()

	if m.GetFocused() {
		content = styles.SelectColumnStyleFocused.Render(choicesView)
	} else {
		content = styles.SelectColumnStyleBlured.Render(choicesView)
	}
	return content
}
