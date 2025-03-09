package model

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
)

var (
	CreateChoices = []string{"Add credentials", "Add text", "Add binary", "Add card"}
)

type CreateModel struct {
	choices    []string
	focusIndex int
}

func (m CreateModel) Init() tea.Cmd {
	return nil
}

func (m CreateModel) Update(_ tui.TUIContext, msg tea.Msg) (tui.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up":
			if m.focusIndex > 0 {
				m.focusIndex--
			}
		case "down":
			if m.focusIndex < len(m.choices)-1 {
				m.focusIndex++
			}
		case "enter":
			if m.focusIndex == 0 {
				return NewLoginRegisterModel(), nil
			}
		}
	}
	return m, nil
}

func (m CreateModel) View() string {
	title := initTitleText
	for i, choice := range m.choices {
		cursor := " "
		if i == m.focusIndex {
			cursor = ">"
		}
		title += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return title
}
