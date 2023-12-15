package main_model

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xsevy/terrapi/helpers"
	"github.com/xsevy/terrapi/messages"
	"github.com/xsevy/terrapi/models/menu"
	"github.com/xsevy/terrapi/templates"
)

type mainModel struct {
	menu *menu.MenuModel
	help help.Model
	keys helpers.KeyMap
}

func NewMainModel(menu *menu.MenuModel) *mainModel {
	return &mainModel{
		keys: helpers.Keys,
		menu: menu,
	}
}

func (m *mainModel) Init() tea.Cmd {
	return nil
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var newMenu tea.Model

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		default:
			newMenu, cmd = m.menu.Update(msg)
			m.menu = newMenu.(*menu.MenuModel)
		}

	case messages.StartSetupMsg:
		newMenu, cmd = m.menu.Update(msg)
		m.menu = newMenu.(*menu.MenuModel)
	case messages.CloseSetupMsg:
		newMenu, cmd = m.menu.Update(msg)
		m.menu = newMenu.(*menu.MenuModel)
	case messages.CreateResourceMsg:
		err := templates.CreateResources(msg.ID, "./", &msg)
		if err != nil {
			panic(err)
		}
		return m, tea.Quit
	}

	return m, cmd
}

func (m *mainModel) View() string {
	menuView := m.menu.View()
	helpView := m.help.View(m.keys)

	height := 8 - strings.Count(helpView, "\n")

	return menuView + strings.Repeat("\n", height) + helpView
}
