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
	binarySecretTitleInputIndex    = 0
	binarySecretMetadataInputIndex = 1
	binarySecretFileButtonIndex    = 2
	binarySecretSubmitButtonIndex  = 3
	binarySecretBackButtonIndex    = 4
	binarySecretTitleInputText     = "Title"
	binarySecretMetadataInputText  = "Metadata"
	binarySecretFileButtonText     = "Select File"
)

type BinarySecretModel struct {
	focusIndex int
	error      string
	buttons    []*components.Button
	inputs     []textinput.Model
	filePath   string
}

func NewBinarySecretModel() *BinarySecretModel {
	return &BinarySecretModel{
		focusIndex: 0,
		inputs: []textinput.Model{
			components.NewInput(components.InputSettings{Placeholder: binarySecretTitleInputText, Focus: true, Style: style.FocusedStyle}),
			components.NewInput(components.InputSettings{Placeholder: binarySecretMetadataInputText}),
		},
		buttons: []*components.Button{
			{binarySecretFileButtonIndex, binarySecretFileButtonText},
			{binarySecretSubmitButtonIndex, components.SubmitButtonText},
			{binarySecretBackButtonIndex, components.BackButtonText},
		},
	}
}

func (m *BinarySecretModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *BinarySecretModel) Update(ctx tui.TUIContext, msg tea.Msg) (tui.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.focusIndex == binarySecretSubmitButtonIndex {
				m.addSecret(ctx)
				return NewMainModel(ctx), nil
			}
			if m.focusIndex == binarySecretBackButtonIndex {
				return NewAddModel(), nil
			}
			if m.focusIndex == binarySecretFileButtonIndex {
				return NewFilePickerModel(m), nil
			}
		case "up":
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = binarySecretBackButtonIndex
			}
			return m.updateInputStyles()
		case "down":
			m.focusIndex++
			if m.focusIndex > binarySecretBackButtonIndex {
				m.focusIndex = 0
			}
			return m.updateInputStyles()
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *BinarySecretModel) updateInputStyles() (tui.Model, tea.Cmd) {
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

func (m *BinarySecretModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *BinarySecretModel) View() string {
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

	b.WriteRune('\n')

	b.WriteString("Selected file: " + m.filePath + "\n")

	if len(m.error) > 0 {
		b.WriteString(style.ErrorStyle.Render(m.error))
	}

	return b.String()
}

func (m *BinarySecretModel) addSecret(ctx tui.TUIContext) {
	err := ctx.SecretService.CreateBinarySecret(
		m.inputs[binarySecretTitleInputIndex].Value(),
		m.filePath,
		m.inputs[binarySecretMetadataInputIndex].Value(),
	)

	if err != nil {
		m.error = err.Error()
	}
}

//package model
//
//import (
//	"github.com/charmbracelet/bubbles/filepicker"
//	"github.com/charmbracelet/bubbles/textinput"
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
//	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/components"
//	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/style"
//	"strings"
//)
//
//const (
//	binarySecretTitleInputIndex    = 0
//	binarySecretFilePickerIndex    = 1
//	binarySecretMetadataInputIndex = 2
//	binarySecretSubmitButtonIndex  = 3
//	binarySecretBackButtonIndex    = 4
//	binarySecretTitleInputText     = "Title"
//	binarySecretFilePickerText     = "Select File"
//	binarySecretMetadataInputText  = "Metadata"
//)
//
//type BinarySecretModel struct {
//	focusIndex int
//	error      string
//	buttons    []*components.Button
//	inputs     []textinput.Model
//	filePicker filepicker.Model
//}
//
//func NewBinarySecretModel() *BinarySecretModel {
//	return &BinarySecretModel{
//		focusIndex: 0,
//		inputs: []textinput.Model{
//			components.NewInput(components.InputSettings{Placeholder: binarySecretTitleInputText, Focus: true, Style: style.FocusedStyle}),
//			components.NewInput(components.InputSettings{Placeholder: binarySecretMetadataInputText}),
//		},
//		filePicker: filepicker.New(),
//		buttons: []*components.Button{
//			{binarySecretSubmitButtonIndex, components.SubmitButtonText},
//			{binarySecretBackButtonIndex, components.BackButtonText},
//		},
//	}
//}
//
//func (m *BinarySecretModel) Init() tea.Cmd {
//	return textinput.Blink
//}
//
//func (m *BinarySecretModel) Update(ctx tui.TUIContext, msg tea.Msg) (tui.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		switch msg.String() {
//		case "enter":
//			if m.focusIndex == binarySecretSubmitButtonIndex {
//				m.addSecret(ctx)
//				return NewMainModel(ctx), nil
//			}
//			if m.focusIndex == binarySecretBackButtonIndex {
//				return NewAddModel(), nil
//			}
//		case "up":
//			m.focusIndex--
//			if m.focusIndex < 0 {
//				m.focusIndex = binarySecretBackButtonIndex
//			}
//			return m.updateInputStyles()
//		case "down":
//			m.focusIndex++
//			if m.focusIndex > binarySecretBackButtonIndex {
//				m.focusIndex = 0
//			}
//			return m.updateInputStyles()
//		}
//	}
//
//	cmd := m.updateInputs(msg)
//
//	return m, cmd
//}
//
//func (m *BinarySecretModel) updateInputStyles() (tui.Model, tea.Cmd) {
//	cmds := make([]tea.Cmd, len(m.inputs))
//
//	for i := range m.inputs {
//		if i == m.focusIndex {
//			cmds[i] = m.inputs[i].Focus()
//			m.inputs[i].PromptStyle = style.FocusedStyle
//			m.inputs[i].TextStyle = style.FocusedStyle
//			continue
//		}
//		m.inputs[i].Blur()
//		m.inputs[i].PromptStyle = style.NoStyle
//		m.inputs[i].TextStyle = style.NoStyle
//	}
//
//	return m, tea.Batch(cmds...)
//}
//
//func (m *BinarySecretModel) updateInputs(msg tea.Msg) tea.Cmd {
//	cmds := make([]tea.Cmd, len(m.inputs))
//
//	for i := range m.inputs {
//		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
//	}
//
//	m.filePicker, _ = m.filePicker.Update(msg)
//
//	return tea.Batch(cmds...)
//}
//
//func (m *BinarySecretModel) View() string {
//	var b strings.Builder
//
//	for i := range m.inputs {
//		b.WriteString(m.inputs[i].View())
//		b.WriteRune('\n')
//	}
//
//	b.WriteString(m.filePicker.View())
//	b.WriteRune('\n')
//
//	for _, btn := range m.buttons {
//		btnStyle := style.BlurredStyle
//		if m.focusIndex == btn.Index {
//			btnStyle = style.FocusedStyle
//		}
//		b.WriteString(btnStyle.Render(btn.Text) + "\n")
//	}
//
//	if len(m.error) > 0 {
//		b.WriteString(style.ErrorStyle.Render(m.error))
//	}
//
//	return b.String()
//}
//
//func (m *BinarySecretModel) addSecret(ctx tui.TUIContext) {
//	err := ctx.SecretService.CreateBinarySecret(
//		m.inputs[binarySecretTitleInputIndex].Value(),
//		m.filePicker.Path,
//		m.inputs[binarySecretMetadataInputIndex].Value(),
//	)
//
//	if err != nil {
//		m.error = err.Error()
//	}
//}
