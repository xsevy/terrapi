package select_column_choices

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xsevy/terrapi/helpers"
	"github.com/xsevy/terrapi/helpers/navigation"
	"github.com/xsevy/terrapi/messages"
	"github.com/xsevy/terrapi/styles"
)

type selectColumnChoice struct {
	id       string
	name     string
	children []selectColumnChoice
	disabled bool
}

func (c selectColumnChoice) GetName() string {
	return c.name
}

func (c selectColumnChoice) GetChildren() []navigation.NavigableItem {
	var items []navigation.NavigableItem
	for _, child := range c.children {
		items = append(items, child)
	}
	return items
}

func (c selectColumnChoice) IsDisabled() bool {
	return c.disabled
}

func (c selectColumnChoice) GetID() string {
	return c.id
}

type SelectColumnChoicesModel struct {
	keys  helpers.KeyMap
	stack navigation.NavigationStack
}

func getInitialItems() []navigation.NavigableItem {
	choices := [2]selectColumnChoice{
		{name: "AppSync", children: []selectColumnChoice{
			{name: helpers.ResourceNames[helpers.ResourceIDs.CreateAppSyncDataSource], id: helpers.ResourceIDs.CreateAppSyncDataSource},
			{name: helpers.ResourceNames[helpers.ResourceIDs.CreateAppSyncAPI], id: helpers.ResourceIDs.CreateAppSyncAPI},
		}},
		{name: "API Gateway", disabled: true},
	}
	initialItems := make([]navigation.NavigableItem, len(choices))
	for i, choice := range choices {
		initialItems[i] = choice
	}

	return initialItems
}

func NewSelectColumnChoicesModel() *SelectColumnChoicesModel {
	initialItems := getInitialItems()

	return &SelectColumnChoicesModel{
		keys:  helpers.Keys,
		stack: navigation.NewNavigationStack(initialItems),
	}
}

func (m *SelectColumnChoicesModel) Init() tea.Cmd {
	return nil
}

func (m *SelectColumnChoicesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentNav := m.stack.CurrentItem()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			currentNav.Selected.Prev()
		case key.Matches(msg, m.keys.Down):
			currentNav.Selected.Next(len(currentNav.Items) - 1)
		case key.Matches(msg, m.keys.Enter):
			selectedItem := currentNav.Items[currentNav.Selected]
			if !selectedItem.IsDisabled() {
				if len(selectedItem.GetChildren()) > 0 {
					m.stack.Push(selectedItem)
				} else {
					return m, messages.SwitchColumn("setup_column", messages.WithID(selectedItem.GetID()))
				}
			}
		case key.Matches(msg, m.keys.Escape):
			m.stack.Pop()
		}
	}
	return m, nil
}

func (m *SelectColumnChoicesModel) View() string {
	currentNav := m.stack.CurrentItem()
	var lines []string

	for i, choice := range currentNav.Items {
		var choiceName string
		if choice.(selectColumnChoice).IsDisabled() {
			choiceName = styles.DisabledChoiceStyle.Render(choice.GetName())
		} else {
			choiceName = choice.GetName()
		}

		if int(currentNav.Selected) == i {
			lines = append(lines, styles.SelectedChoiceStyle.Render(choiceName))
		} else {
			lines = append(lines, styles.ChoiceStyle.Render(choiceName))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Top, lines...)
}
