package bubbles

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xsevy/terrapi/helpers"
	"github.com/xsevy/terrapi/helpers/functions"
	"github.com/xsevy/terrapi/helpers/navigation"
	"github.com/xsevy/terrapi/styles"
)

var (
	itemStyle         = lipgloss.NewStyle()
	selectedItemStyle = itemStyle.Copy().Background(lipgloss.Color(helpers.Colors.Purple))
)

type ListModel struct {
	keys      helpers.KeyMap
	title     string
	items     []string
	paginator paginator.Model
	selected  navigation.Selected
	focused   bool
}

func NewListModel(title string, items []string, sorted bool, focused bool, initialSelected int) *ListModel {
	lm := ListModel{
		keys:      helpers.Keys,
		title:     title,
		selected:  navigation.Selected(initialSelected),
		paginator: NewPaginator(5, len(items)),
		focused:   focused,
	}
	if sorted {
		lm.items = functions.SortSliceCaseInsensitive(items)
	} else {
		lm.items = items
	}

	return &lm
}

func (m *ListModel) Init() tea.Cmd {
	return nil
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	itemsLength := len(m.items)

	start, end := m.paginator.GetSliceBounds(itemsLength)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.selected.Prev()
			if int(m.selected) < start {
				m.paginator.PrevPage()
			}
		case key.Matches(msg, m.keys.Down):
			m.selected.Next(itemsLength - 1)
			if int(m.selected) == end {
				m.paginator.NextPage()
			}
		}
	}

	m.paginator, cmd = m.paginator.Update(msg)

	return m, cmd
}

func (m *ListModel) View() string {
	var b strings.Builder

	title := styles.GetFocusedTitle(m.title, m.focused)
	b.WriteString(title + "\n\n")

	start, end := m.paginator.GetSliceBounds(len(m.items))
	for i := start; i < end; i++ {
		item := m.items[i]
		if i == int(m.selected) {
			b.WriteString(selectedItemStyle.Render(item) + "\n\n")
		} else {
			b.WriteString(itemStyle.Render(item) + "\n\n")
		}
	}

	b.WriteString(m.paginator.View() + "\n")

	return b.String()
}

func (m *ListModel) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *ListModel) Blur() {
	m.focused = false
}

func (m *ListModel) Value() string {
	return m.items[m.selected]
}
