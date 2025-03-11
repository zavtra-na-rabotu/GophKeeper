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
	authLoginInputIndex     = 0
	authPasswordInputIndex  = 1
	authLoginButtonIndex    = 2
	authRegisterButtonIndex = 3
	authBackButtonIndex     = 4
	authLoginLimit          = 255
	authPasswordLimit       = 50
	authLoginButtonText     = "[ Login ]"
	authRegisterButtonText  = "[ Register ]"
	authLoginInputText      = "Login"
	authPasswordInputText   = "Password"
)

type AuthModel struct {
	focusIndex int
	error      string
	buttons    []*components.Button
	inputs     []textinput.Model
	ctx        *tui.TUIContext
}

func NewAuthModel(ctx *tui.TUIContext) AuthModel {
	m := AuthModel{
		focusIndex: 0,
		inputs: []textinput.Model{
			components.NewInput(components.InputSettings{Placeholder: authLoginInputText, Focus: true, CharLimit: authLoginLimit, Style: style.FocusedStyle}),
			components.NewInput(components.InputSettings{Placeholder: authPasswordInputText, CharLimit: authPasswordLimit}),
		},
		buttons: []*components.Button{
			{authLoginButtonIndex, authLoginButtonText},
			{authRegisterButtonIndex, authRegisterButtonText},
			{authBackButtonIndex, components.BackButtonText},
		},
		ctx: ctx,
	}

	m.inputs[authPasswordInputIndex].EchoMode = textinput.EchoPassword
	m.inputs[authPasswordInputIndex].EchoCharacter = '•'

	return m
}

func (m AuthModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Handle login button
			if m.focusIndex == authLoginButtonIndex {
				err := m.loginUser(m.ctx)
				if err != nil {
					m.error = err.Error()
				} else {
					return NewMainModel(m.ctx), nil
				}
			}

			// Handle register button
			if m.focusIndex == authRegisterButtonIndex {
				err := m.registerUser(m.ctx)
				if err != nil {
					m.error = err.Error()
				} else {
					return NewMainModel(m.ctx), nil
				}
			}

			// Handle back button
			if m.focusIndex == authBackButtonIndex {
				return NewInitModel(m.ctx), nil
			}

		case "up":
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = authBackButtonIndex
			}
			return m.updateInputStyles()

		case "down":
			m.focusIndex++
			if m.focusIndex > authBackButtonIndex {
				m.focusIndex = 0
			}
			return m.updateInputStyles()
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m AuthModel) updateInputStyles() (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
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

func (m AuthModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m AuthModel) View() string {
	var b strings.Builder

	b.WriteString(style.HintStyle.Render("Use (↑, ↓, 'Enter') to navigate") + "\n\n")

	// Render inputs
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}

	b.WriteRune('\n')

	// Render buttons
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

func (m AuthModel) loginUser(ctx *tui.TUIContext) error {
	token, err := ctx.UserService.Login(m.inputs[authLoginInputIndex].Value(), m.inputs[authPasswordInputIndex].Value())
	if err != nil {
		return err
	}

	ctx.SecretService.SetToken(token)
	ctx.SecretService.SetPassword(m.inputs[authPasswordInputIndex].Value())

	return nil
}

func (m AuthModel) registerUser(ctx *tui.TUIContext) error {
	token, err := ctx.UserService.Register(m.inputs[authLoginInputIndex].Value(), m.inputs[authPasswordInputIndex].Value())
	if err != nil {
		return err
	}

	ctx.SecretService.SetToken(token)
	ctx.SecretService.SetPassword(m.inputs[authPasswordInputIndex].Value())

	return nil
}
