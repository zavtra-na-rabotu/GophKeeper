package utils

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type InputSettings struct {
	Placeholder string
	CharLimit   int
	Focus       bool
	Style       lipgloss.Style
}

func NewInput(settings InputSettings) textinput.Model {
	t := textinput.New()
	t.CharLimit = settings.CharLimit
	t.Placeholder = settings.Placeholder

	if settings.Focus {
		t.Focus()
		t.PromptStyle = settings.Style
		t.TextStyle = settings.Style
	}

	return t
}
