package model

import (
	"fmt"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
	"os"
	"path/filepath"
	"strings"
)

type FilePickerModel struct {
	parent *BinarySecretModel
	picker filepicker.Model
}

func NewFilePickerModel(parent *BinarySecretModel) *FilePickerModel {
	defaultPath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	secretFilePicker := filepicker.New()
	secretFilePicker.CurrentDirectory = filepath.Join(defaultPath)
	secretFilePicker.AutoHeight = false
	secretFilePicker.Height = 20

	secretFilePicker.Init()

	return &FilePickerModel{
		parent: parent,
		picker: secretFilePicker,
	}
}

func (m FilePickerModel) Init() tea.Cmd {
	return m.picker.Init()
}

func (m *FilePickerModel) Update(ctx tui.TUIContext, msg tea.Msg) (tui.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m.parent, nil
		}
	}

	m.picker, cmd = m.picker.Update(msg)

	if m.picker.Path != "" {
		m.parent.filePath = m.picker.Path
		return m.parent, nil
	}

	return m, cmd
}

func (m FilePickerModel) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%20s%s:\n", "", m.picker.CurrentDirectory))
	b.WriteString(m.picker.View())

	return b.String()
}
