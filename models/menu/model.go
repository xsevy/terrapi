package menu

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xsevy/terrapi/helpers"
	"github.com/xsevy/terrapi/messages"
	"github.com/xsevy/terrapi/models/select_column"
	"github.com/xsevy/terrapi/models/setup_column"
)

type MenuModel struct {
	selectColumn *select_column.SelectColumnModel
	setupColumn  *setup_column.SetupColumnModel
	keys         helpers.KeyMap
}

func NewMenuModel(selectColumn *select_column.SelectColumnModel, setupColumn *setup_column.SetupColumnModel) *MenuModel {
	return &MenuModel{
		selectColumn: selectColumn,
		setupColumn:  setupColumn,
		keys:         helpers.Keys,
	}
}

func (m *MenuModel) Init() tea.Cmd {
	return nil
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var newSetupColumn tea.Model

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up, m.keys.Down, m.keys.Enter, m.keys.Escape):
			var newSelectColumn tea.Model

			if m.selectColumn.GetFocused() {
				newSelectColumn, cmd = m.selectColumn.Update(msg)
				m.selectColumn = newSelectColumn.(*select_column.SelectColumnModel)
			} else {
				newSetupColumn, cmd = m.setupColumn.Update(msg)
				m.setupColumn = newSetupColumn.(*setup_column.SetupColumnModel)
			}
		default:
			if m.setupColumn.GetFocused() {
				newSetupColumn, cmd = m.setupColumn.Update(msg)
				m.setupColumn = newSetupColumn.(*setup_column.SetupColumnModel)
			}
		}
	case messages.StartSetupMsg:
		m.switchColumn(msg.ID)
	case messages.CloseSetupMsg:
		m.switchColumn("")
	}
	return m, cmd
}

func (m *MenuModel) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.selectColumn.View(),
		m.setupColumn.View(),
	)
}

func (m *MenuModel) switchColumn(id string) {
	if id != "" {
		m.setupColumn.SetID(id)
	}
	m.selectColumn.SetFocused(!m.selectColumn.GetFocused())
	m.setupColumn.SetFocused(!m.setupColumn.GetFocused())
}
