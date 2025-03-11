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
	cardSecretTitleInputIndex       = 0
	cardSecretNumberInputIndex      = 1
	cardSecretExpiryMonthInputIndex = 2
	cardSecretExpiryYearInputIndex  = 3
	cardSecretCSCInputIndex         = 4
	cardSecretNameInputIndex        = 5
	cardSecretMetadataInputIndex    = 6
	cardSecretSubmitButtonIndex     = 7
	cardSecretBackButtonIndex       = 8
	cardSecretTitleInputText        = "Title"
	cardSecretNumberInputText       = "Card Number"
	cardSecretExpiryMonthInputText  = "Expiry Month (MM)"
	cardSecretExpiryYearInputText   = "Expiry Year (YY)"
	cardSecretCSCInputText          = "CSC"
	cardSecretNameInputText         = "Cardholder Name"
	cardSecretMetadataInputText     = "Metadata"
)

type CardSecretModel struct {
	focusIndex int
	error      string
	buttons    []*components.Button
	inputs     []textinput.Model
	ctx        *tui.TUIContext
	secretID   uint64
}

func NewCardSecretModel(ctx *tui.TUIContext, secret *pb.Secret, content *pb.Card) *CardSecretModel {
	model := &CardSecretModel{
		focusIndex: 0,
		inputs: []textinput.Model{
			components.NewInput(components.InputSettings{Placeholder: cardSecretTitleInputText, Focus: true, Style: style.FocusedStyle}),
			components.NewInput(components.InputSettings{Placeholder: cardSecretNumberInputText}),
			components.NewInput(components.InputSettings{Placeholder: cardSecretExpiryMonthInputText, CharLimit: 2}),
			components.NewInput(components.InputSettings{Placeholder: cardSecretExpiryYearInputText, CharLimit: 2}),
			components.NewInput(components.InputSettings{Placeholder: cardSecretCSCInputText, CharLimit: 4}),
			components.NewInput(components.InputSettings{Placeholder: cardSecretNameInputText}),
			components.NewInput(components.InputSettings{Placeholder: cardSecretMetadataInputText}),
		},
		buttons: []*components.Button{
			{cardSecretSubmitButtonIndex, components.SubmitButtonText},
			{cardSecretBackButtonIndex, components.BackButtonText},
		},
		ctx: ctx,
	}

	if secret != nil {
		model.secretID = secret.Id
		model.inputs[cardSecretTitleInputIndex].SetValue(secret.GetTitle())
		model.inputs[cardSecretMetadataInputIndex].SetValue(secret.GetMetadata())
	}

	if content != nil {
		model.inputs[cardSecretNumberInputIndex].SetValue(content.GetNumber())
		model.inputs[cardSecretExpiryMonthInputIndex].SetValue(content.GetExpiryMonth())
		model.inputs[cardSecretExpiryYearInputIndex].SetValue(content.GetExpiryYear())
		model.inputs[cardSecretCSCInputIndex].SetValue(content.GetCsc())
		model.inputs[cardSecretNameInputIndex].SetValue(content.GetName())
	}

	return model
}

func (m *CardSecretModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *CardSecretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.focusIndex == cardSecretSubmitButtonIndex {
				m.addSecret(m.ctx)
				return NewMainModel(m.ctx), nil
			}
			if m.focusIndex == cardSecretBackButtonIndex {
				return NewAddModel(m.ctx), nil
			}
		case "up":
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = cardSecretBackButtonIndex
			}
			return m.updateInputStyles()
		case "down":
			m.focusIndex++
			if m.focusIndex > cardSecretBackButtonIndex {
				m.focusIndex = 0
			}
			return m.updateInputStyles()
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *CardSecretModel) updateInputStyles() (tea.Model, tea.Cmd) {
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

func (m *CardSecretModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *CardSecretModel) View() string {
	var b strings.Builder

	b.WriteString(style.HintStyle.Render("Use (↑, ↓, 'Enter') to navigate") + "\n\n")

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

func (m *CardSecretModel) addSecret(ctx *tui.TUIContext) {
	err := ctx.SecretService.CreateCardSecret(
		m.secretID,
		m.inputs[cardSecretTitleInputIndex].Value(),
		m.inputs[cardSecretNumberInputIndex].Value(),
		m.inputs[cardSecretExpiryMonthInputIndex].Value(),
		m.inputs[cardSecretExpiryYearInputIndex].Value(),
		m.inputs[cardSecretCSCInputIndex].Value(),
		m.inputs[cardSecretNameInputIndex].Value(),
		m.inputs[cardSecretMetadataInputIndex].Value(),
	)

	if err != nil {
		m.error = err.Error()
	}
}
