package model

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/style"
	"strings"
)

const (
	initLoginRegisterIndex      = 0
	initExitIndex               = 1
	initLastElementIndex        = 1
	initLoginRegisterButtonText = "[ Login / Register ]"
	initExitButtonTest          = "[ Exit ]"
)

var (
	initLoginRegisterButton        = style.BlurredStyle.Render(initLoginRegisterButtonText)
	initLoginRegisterButtonFocused = style.FocusedStyle.Render(initLoginRegisterButtonText)
	initExitButton                 = style.BlurredStyle.Render(initExitButtonTest)
	initExitButtonFocused          = style.FocusedStyle.Render(initExitButtonTest)
)

type InitModel struct {
	focusIndex int
}

func NewInitModel() *InitModel {
	return &InitModel{
		focusIndex: 0,
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
			if m.focusIndex < initLastElementIndex {
				m.focusIndex++
			}
		case "enter":
			if m.focusIndex == initLoginRegisterIndex {
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

	loginRegisterBtn := initLoginRegisterButton
	if m.focusIndex == initLoginRegisterIndex {
		loginRegisterBtn = initLoginRegisterButtonFocused
	}

	exitBtn := initExitButton
	if m.focusIndex == initExitIndex {
		exitBtn = initExitButtonFocused
	}

	fmt.Fprintf(&b, "\n\n%s\n%s\n\n", loginRegisterBtn, exitBtn)

	return b.String()
}
