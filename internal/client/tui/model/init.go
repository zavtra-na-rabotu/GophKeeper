package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/components"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/style"
	"strings"
)

const (
	initAuthIndex      = 0
	initExitIndex      = 1
	initAuthButtonText = "[ Login / Register ]"
	initExitButtonText = "[ Exit ]"
)

type InitModel struct {
	focusIndex int
	buttons    []*components.Button
}

func NewInitModel() *InitModel {
	return &InitModel{
		focusIndex: 0,
		buttons: []*components.Button{
			{initAuthIndex, initAuthButtonText},
			{initExitIndex, initExitButtonText},
		},
	}
}

func (m InitModel) Init() tea.Cmd {
	return nil
}

func (m InitModel) Update(_ tui.TUIContext, msg tea.Msg) (tui.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.focusIndex > 0 {
				m.focusIndex--
			}
		case "down":
			if m.focusIndex < initExitIndex {
				m.focusIndex++
			}
		case "enter":
			if m.focusIndex == initAuthIndex {
				return NewAuthModel(), nil
			}
			if m.focusIndex == initExitIndex {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m InitModel) View() string {
	var b strings.Builder

	for _, btn := range m.buttons {
		btnStyle := style.BlurredStyle
		if m.focusIndex == btn.Index {
			btnStyle = style.FocusedStyle
		}
		b.WriteString(btnStyle.Render(btn.Text) + "\n")
	}

	return b.String()
}
