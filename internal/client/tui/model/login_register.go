package model

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/style"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/utils"
	"strings"
)

const (
	buttonsNumber       = 2
	loginButtonIndex    = 2
	registerButtonIndex = 3
)

var (
	cursorStyle           = style.FocusedStyle
	registerButton        = fmt.Sprintf("[ %s ]", style.BlurredStyle.Render("Register"))
	registerButtonFocused = style.FocusedStyle.Render("[ Register ]")
	loginButton           = fmt.Sprintf("[ %s ]", style.BlurredStyle.Render("Login"))
	loginButtonFocused    = style.FocusedStyle.Render("[ Login ]")
)

type LoginRegisterModel struct {
	focusIndex int
	inputs     []textinput.Model
	error      string
}

func NewLoginRegisterModel() LoginRegisterModel {
	m := LoginRegisterModel{
		inputs: make([]textinput.Model, 3),
	}

	loginInput := utils.NewInput(utils.InputSettings{
		Placeholder: "Login",
		Focus:       true,
		CharLimit:   255,
		Style:       style.FocusedStyle,
	})

	passwordInput := utils.NewInput(utils.InputSettings{
		Placeholder: "Password",
		Focus:       false,
		CharLimit:   50,
	})
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'

	m.inputs[0] = loginInput
	m.inputs[1] = passwordInput

	return m
}

func (m LoginRegisterModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LoginRegisterModel) Update(app tui.TUIContext, msg tea.Msg) (tui.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// TODO: Сделать возврат на init menu

		case "ctrl+c", "esc":
			return NewInitModel(Choices, 0), nil

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == loginButtonIndex {
				err := m.loginUser(app)
				if err != nil {
					m.error = err.Error()
				} else {
					return NewMainModel(app), nil
				}
			}

			if s == "enter" && m.focusIndex == registerButtonIndex {
				err := m.registerUser(app)
				if err != nil {
					m.error = err.Error()
				} else {
					// TODO: GOTO Main menu
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-buttonsNumber; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = style.FocusedStyle
					m.inputs[i].TextStyle = style.FocusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = style.NoStyle
				m.inputs[i].TextStyle = style.NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m LoginRegisterModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m LoginRegisterModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-buttonsNumber {
			b.WriteRune('\n')
		}
	}

	logButton := &loginButton
	if m.focusIndex == loginButtonIndex {
		logButton = &loginButtonFocused
	}

	regButton := &registerButton
	if m.focusIndex == registerButtonIndex {
		regButton = &registerButtonFocused
	}

	fmt.Fprintf(&b, "\n\n%s\n%s\n\n", *logButton, *regButton)

	if len(m.error) > 0 {
		b.WriteString(style.ErrorStyle.Render(m.error))
	}

	return b.String()
}

func (m LoginRegisterModel) loginUser(ctx tui.TUIContext) error {
	token, err := ctx.UserService.Login(m.inputs[0].Value(), m.inputs[1].Value())
	if err != nil {
		return err
	}

	ctx.SecretService.SetToken(token)
	ctx.SecretService.SetPassword(m.inputs[1].Value())

	return nil
}

func (m LoginRegisterModel) registerUser(ctx tui.TUIContext) error {
	token, err := ctx.UserService.Register(m.inputs[0].Value(), m.inputs[1].Value())
	if err != nil {
		return err
	}

	ctx.SecretService.SetToken(token)
	ctx.SecretService.SetPassword(m.inputs[1].Value())

	return nil
}
