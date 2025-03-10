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
	textSecretTitleInputIndex    = 0
	textSecretTextInputIndex     = 1
	textSecretMetadataInputIndex = 2
	textSecretSubmitButtonIndex  = 3
	textSecretBackButtonIndex    = 4
	textSecretTitleInputText     = "Title"
	textSecretTextInputText      = "Text"
	textSecretMetadataInputText  = "Metadata"
)

type TextSecretModel struct {
	focusIndex int
	error      string
	buttons    []*components.Button
	inputs     []textinput.Model
	ctx        *tui.TUIContext
}

func NewTextSecretModel(ctx *tui.TUIContext) *TextSecretModel {
	return &TextSecretModel{
		focusIndex: 0,
		inputs: []textinput.Model{
			components.NewInput(components.InputSettings{Placeholder: textSecretTitleInputText, Focus: true, Style: style.FocusedStyle}),
			components.NewInput(components.InputSettings{Placeholder: textSecretTextInputText}),
			components.NewInput(components.InputSettings{Placeholder: textSecretMetadataInputText}),
		},
		buttons: []*components.Button{
			{textSecretSubmitButtonIndex, components.SubmitButtonText},
			{textSecretBackButtonIndex, components.BackButtonText},
		},
		ctx: ctx,
	}
}

func (m *TextSecretModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *TextSecretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.focusIndex == textSecretSubmitButtonIndex {
				m.addSecret(m.ctx)
				return NewMainModel(m.ctx), nil
			}
			if m.focusIndex == textSecretBackButtonIndex {
				return NewAddModel(m.ctx), nil
			}
		case "up":
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = textSecretBackButtonIndex
			}
			return m.updateInputStyles()
		case "down":
			m.focusIndex++
			if m.focusIndex > textSecretBackButtonIndex {
				m.focusIndex = 0
			}
			return m.updateInputStyles()
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *TextSecretModel) updateInputStyles() (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = style.FocusedStyle
			m.inputs[i].TextStyle = style.FocusedStyle
			continue
		}
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = style.NoStyle
		m.inputs[i].TextStyle = style.NoStyle
	}

	return m, tea.Batch(cmds...)
}

func (m *TextSecretModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *TextSecretModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}

	b.WriteRune('\n')

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

func (m *TextSecretModel) addSecret(ctx *tui.TUIContext) {
	err := ctx.SecretService.CreateTextSecret(
		m.inputs[textSecretTitleInputIndex].Value(),
		m.inputs[textSecretTextInputIndex].Value(),
		m.inputs[textSecretMetadataInputIndex].Value(),
	)

	if err != nil {
		m.error = err.Error()
	}
}
