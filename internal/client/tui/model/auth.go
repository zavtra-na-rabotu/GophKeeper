package model

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/components"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/style"
	"strings"
)

const (
	authLastElementIndex    = 4
	authLastInputIndex      = 1
	authLoginInputIndex     = 0
	authPasswordInputIndex  = 1
	authLoginButtonIndex    = 2
	authRegisterButtonIndex = 3
	authBackButtonIndex     = 4
	authLoginLimit          = 255
	authPasswordLimit       = 50
	authLoginButtonText     = "[ Login ]"
	authRegisterButtonText  = "[ Register ]"
	authBackButtonText      = "[ Back ]"
)

type AuthModel struct {
	focusIndex int
	error      string
	buttons    []*components.Button
	inputs     []textinput.Model
}

func NewAuthModel() AuthModel {
	m := AuthModel{
		focusIndex: 0,
		inputs:     make([]textinput.Model, 2),
		buttons: []*components.Button{
			{authLoginButtonIndex, authLoginButtonText},
			{authRegisterButtonIndex, authRegisterButtonText},
			{authBackButtonIndex, authBackButtonText},
		},
	}

	loginInput := components.NewInput(components.InputSettings{
		Placeholder: "Login",
		Focus:       true,
		CharLimit:   authLoginLimit,
		Style:       style.FocusedStyle,
	})

	passwordInput := components.NewInput(components.InputSettings{
		Placeholder: "Password",
		Focus:       false,
		CharLimit:   authPasswordLimit,
	})
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'

	m.inputs[authLoginInputIndex] = loginInput
	m.inputs[authPasswordInputIndex] = passwordInput

	return m
}

func (m AuthModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AuthModel) Update(ctx tui.TUIContext, msg tea.Msg) (tui.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return NewInitModel(), nil

		case "enter":
			// Handle login button
			if m.focusIndex == authLoginButtonIndex {
				err := m.loginUser(ctx)
				if err != nil {
					m.error = err.Error()
				} else {
					return NewMainModel(ctx), nil
				}
			}

			// Handle register button
			if m.focusIndex == authRegisterButtonIndex {
				err := m.registerUser(ctx)
				if err != nil {
					m.error = err.Error()
				} else {
					return NewMainModel(ctx), nil
				}
			}

			// Handle back button
			if m.focusIndex == authBackButtonIndex {
				return NewInitModel(), nil
			}

		case "up", "down":
			s := msg.String()

			if s == "up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > authLastElementIndex {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = authLastElementIndex
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= authLastInputIndex; i++ {
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

func (m AuthModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m AuthModel) View() string {
	var b strings.Builder

	// Render inputs
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < authLastInputIndex {
			b.WriteRune('\n')
		}
	}

	b.WriteString("\n\n")

	for _, btn := range m.buttons {
		btnStyle := style.BlurredStyle
		if m.focusIndex == btn.Index {
			btnStyle = style.FocusedStyle
		}
		b.WriteString(btnStyle.Render(btn.Text) + "\n")
	}

	if len(m.error) > 0 {
		b.WriteString(style.ErrorStyle.Render(m.error))
	}

	return b.String()
}

func (m AuthModel) loginUser(ctx tui.TUIContext) error {
	token, err := ctx.UserService.Login(m.inputs[0].Value(), m.inputs[1].Value())
	if err != nil {
		return err
	}

	ctx.SecretService.SetToken(token)
	ctx.SecretService.SetPassword(m.inputs[1].Value())

	return nil
}

func (m AuthModel) registerUser(ctx tui.TUIContext) error {
	token, err := ctx.UserService.Register(m.inputs[0].Value(), m.inputs[1].Value())
	if err != nil {
		return err
	}

	ctx.SecretService.SetToken(token)
	ctx.SecretService.SetPassword(m.inputs[1].Value())

	return nil
}
