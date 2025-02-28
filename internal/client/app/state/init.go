package state

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	TitleText = "Choose using arrow keys and 'Enter':\n\n"
)

var (
	Choices = []string{"Login", "Register", "Exit"}
)

type InitState struct {
	Choices []string
	Cursor  int
}

func (s InitState) Init() tea.Cmd {
	return nil
}

func (s InitState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return s, tea.Quit
		case "up":
			if s.Cursor > 0 {
				s.Cursor--
			}
		case "down":
			if s.Cursor < len(s.Choices)-1 {
				s.Cursor++
			}
		case "enter":
			if s.Cursor == 0 {
				return LoginState{InputPos: 0}, nil
			}
		}
	}
	return s, nil
}

func (s InitState) View() string {
	title := TitleText
	for i, choice := range s.Choices {
		cursor := " "
		if i == s.Cursor {
			cursor = ">"
		}
		title += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return title
}
