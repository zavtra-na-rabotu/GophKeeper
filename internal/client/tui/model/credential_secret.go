package model

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/components"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/style"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"strings"
)

const (
	credentialTitleInputIndex    = 0
	credentialLoginInputIndex    = 1
	credentialPasswordInputIndex = 2
	credentialMetadataInputIndex = 3
	credentialSubmitButtonIndex  = 4
	credentialBackButtonIndex    = 5
	credentialTitleInputText     = "Title"
	credentialLoginInputText     = "Login"
	credentialPasswordInputText  = "Password"
	credentialMetadataInputText  = "Metadata"
)

type CredentialSecretModel struct {
	focusIndex int
	error      string
	buttons    []*components.Button
	inputs     []textinput.Model
	ctx        *tui.TUIContext
}

func NewCredentialSecretModel(ctx *tui.TUIContext, secret *pb.Secret, content *pb.Credential) *CredentialSecretModel {
	model := &CredentialSecretModel{
		focusIndex: 0,
		inputs: []textinput.Model{
			components.NewInput(components.InputSettings{Placeholder: credentialTitleInputText, Focus: true, Style: style.FocusedStyle}),
			components.NewInput(components.InputSettings{Placeholder: credentialLoginInputText}),
			components.NewInput(components.InputSettings{Placeholder: credentialPasswordInputText}),
			components.NewInput(components.InputSettings{Placeholder: credentialMetadataInputText}),
		},
		buttons: []*components.Button{
			{credentialSubmitButtonIndex, components.SubmitButtonText},
			{credentialBackButtonIndex, components.BackButtonText},
		},
		ctx: ctx,
	}
	
	if secret != nil {
		model.inputs[credentialTitleInputIndex].SetValue(secret.GetTitle())
		model.inputs[credentialMetadataInputIndex].SetValue(secret.GetMetadata())
	}

	if content != nil {
		model.inputs[credentialLoginInputIndex].SetValue(content.GetLogin())
		model.inputs[credentialPasswordInputIndex].SetValue(content.GetPassword())
	}

	return model
}

func (m *CredentialSecretModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *CredentialSecretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.focusIndex == credentialSubmitButtonIndex {
				m.addSecret(m.ctx)
				return NewMainModel(m.ctx), nil
			}
			if m.focusIndex == credentialBackButtonIndex {
				return NewAddModel(m.ctx), nil
			}

		case "up":
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = credentialBackButtonIndex
			}
			return m.updateInputStyles()

		case "down":
			m.focusIndex++
			if m.focusIndex > credentialBackButtonIndex {
				m.focusIndex = 0
			}
			return m.updateInputStyles()
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *CredentialSecretModel) updateInputStyles() (tea.Model, tea.Cmd) {
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

func (m *CredentialSecretModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *CredentialSecretModel) View() string {
	var b strings.Builder

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

func (m *CredentialSecretModel) addSecret(ctx *tui.TUIContext) {
	err := ctx.SecretService.CreateCredentialSecret(
		m.inputs[credentialTitleInputIndex].Value(),
		m.inputs[credentialLoginInputIndex].Value(),
		m.inputs[credentialPasswordInputIndex].Value(),
		m.inputs[credentialMetadataInputIndex].Value(),
	)

	if err != nil {
		m.error = err.Error()
	}
}
