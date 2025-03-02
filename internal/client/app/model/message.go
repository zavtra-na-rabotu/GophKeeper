package model

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/app"
)

type MessageModel struct {
	Message string
}

func NewMessageModel(message string) *MessageModel {
	return &MessageModel{
		Message: message,
	}
}

func (m MessageModel) Update(app app.App, msg tea.Msg) (app.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return NewInitModel(Choices, 0), nil
	}
	return m, nil
}

func (m MessageModel) View() string {
	return fmt.Sprintf(m.Message)
}
