package model

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/app"
)

const (
	TitleText = "Choose using arrow keys and 'Enter':\n\n"
)

var (
	Choices = []string{"Login", "Register", "Exit"}
)

type InitModel struct {
	Choices []string
	Cursor  int
}

func NewInitModel(choices []string, cursor int) *InitModel {
	return &InitModel{
		Choices: choices,
		Cursor:  cursor,
	}
}

func (m InitModel) Update(_ app.App, msg tea.Msg) (app.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "enter":
			if m.Cursor == 0 {
				return LoginModel{InputPos: 0}, nil
			}
		}
	}
	return m, nil
}

func (m InitModel) View() string {
	title := TitleText
	for i, choice := range m.Choices {
		cursor := " "
		if i == m.Cursor {
			cursor = ">"
		}
		title += fmt.Sprintf("%m %m\n", cursor, choice)
	}
	return title
}
