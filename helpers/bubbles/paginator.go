package bubbles

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/lipgloss"
	"github.com/xsevy/terrapi/helpers"
)

func NewPaginator(perpage int, itemsLength int) paginator.Model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = perpage
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.Color(helpers.Colors.Purple)).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.Color(helpers.Colors.Grey)).Render("•")
	p.SetTotalPages(itemsLength)

	return p
}
