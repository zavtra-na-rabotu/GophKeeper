package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/components"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/style"
	"strings"
)

const (
	addLastElementIndex      = 4
	addCredentialButtonIndex = 0
	addTextButtonIndex       = 1
	addBinaryButtonIndex     = 2
	addCardButtonIndex       = 3
	addBackButtonIndex       = 4
	addCredentialsButtonText = "[ Add credentials ]"
	addTextButtonText        = "[ Add text ]"
	addBinaryButtonText      = "[ Add binary ]"
	addCardButtonText        = "[ Add card ]"
	addBackButtonText        = "[ Back ]"
)

type AddModel struct {
	focusIndex int
	buttons    []*components.Button
}

func NewAddModel() *AddModel {
	return &AddModel{
		focusIndex: 0,
		buttons: []*components.Button{
			{addCredentialButtonIndex, addCredentialsButtonText},
			{addTextButtonIndex, addTextButtonText},
			{addBinaryButtonIndex, addBinaryButtonText},
			{addCardButtonIndex, addCardButtonText},
			{addBackButtonIndex, addBackButtonText},
		},
	}
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m AddModel) Update(_ tui.TUIContext, msg tea.Msg) (tui.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.focusIndex > 0 {
				m.focusIndex--
			}
		case "down":
			if m.focusIndex < addLastElementIndex {
				m.focusIndex++
			}
		case "enter":
			if m.focusIndex == addLastElementIndex {

			}
		}
	}

	return m, nil
}

func (m AddModel) View() string {
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
