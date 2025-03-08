package model

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
)

const (
	TitleText = "Choose using arrow keys and 'Enter':\n\n"
)

var (
	Choices = []string{"Login/Register", "Exit"}
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

func (m InitModel) Init() tea.Cmd {
	return nil
}

func (m InitModel) Update(_ tui.TUIContext, msg tea.Msg) (tui.Model, tea.Cmd) {
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
				return NewLoginRegisterModel(), nil
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
		title += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return title
}
