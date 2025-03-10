package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/components"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/style"
	"strings"
)

const (
	addCredentialButtonIndex = 0
	addTextButtonIndex       = 1
	addBinaryButtonIndex     = 2
	addCardButtonIndex       = 3
	addBackButtonIndex       = 4
	addCredentialsButtonText = "[ Add credentials ]"
	addTextButtonText        = "[ Add text ]"
	addBinaryButtonText      = "[ Add binary ]"
	addCardButtonText        = "[ Add card ]"
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
			{addBackButtonIndex, components.BackButtonText},
		},
	}
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m AddModel) Update(ctx tui.TUIContext, msg tea.Msg) (tui.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = addBackButtonIndex
			}
		case "down":
			m.focusIndex++
			if m.focusIndex > addBackButtonIndex {
				m.focusIndex = 0
			}
		case "enter":
			if m.focusIndex == addCredentialButtonIndex {
				return NewCredentialSecretModel(), nil
			}
			if m.focusIndex == addTextButtonIndex {
				return NewTextSecretModel(), nil
			}
			if m.focusIndex == addBinaryButtonIndex {
				return NewBinarySecretModel(), nil
			}
			if m.focusIndex == addCardButtonIndex {

			}
			if m.focusIndex == addBackButtonIndex {
				return NewMainModel(ctx), nil
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
