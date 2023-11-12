package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/xsevy/terrapi/helpers"
)

var (
	commonColumnStyle  = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).Height(28).PaddingRight(1).PaddingLeft(1)
	bluredColumnStyle  = commonColumnStyle.Copy().BorderForeground(lipgloss.Color(helpers.Colors.Grey))
	focusedColumnStyle = commonColumnStyle.Copy().BorderForeground(lipgloss.Color(helpers.Colors.Purple))

	SelectColumnStyleFocused = lipgloss.NewStyle().Width(25).PaddingLeft(2).Inherit(focusedColumnStyle)
	SelectColumnStyleBlured  = lipgloss.NewStyle().Width(25).PaddingLeft(2).Inherit(bluredColumnStyle)

	SetupColumnStyleFocused = lipgloss.NewStyle().Width(50).PaddingLeft(2).Inherit(focusedColumnStyle)
	SetupColumnStyleBlured  = lipgloss.NewStyle().Width(50).PaddingLeft(2).Inherit(bluredColumnStyle)

	ChoiceStyle         = lipgloss.NewStyle()
	SelectedChoiceStyle = lipgloss.NewStyle().Background(lipgloss.Color(helpers.Colors.Purple))
	DisabledChoiceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(helpers.Colors.Grey))

	focusedTitle = lipgloss.NewStyle().Underline(true)
)

func GetFocusedTitle(title string, focused bool) string {
	if focused {
		return focusedTitle.Render(title)
	}

	return title
}
