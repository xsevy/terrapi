package messages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type StartSetupMsg struct {
	ID string
}
type CloseSetupMsg struct{}

type startSetupOption func(*StartSetupMsg)

func SwitchColumn(to string, options ...startSetupOption) tea.Cmd {
	return func() tea.Msg {
		switch to {
		case "setup_column":
			msg := &StartSetupMsg{}
			for _, opt := range options {
				opt(msg)
			}
			return *msg
		case "select_column":
			msg := &CloseSetupMsg{}
			return *msg
		default:
			panic(fmt.Sprintf("Unknown value of 'to': %s", to))
		}
	}
}

func WithID(id string) startSetupOption {
	return func(msg *StartSetupMsg) {
		msg.ID = id
	}
}
