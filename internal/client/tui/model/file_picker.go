package model

import (
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
	ctx    *tui.TUIContext
}

func NewFilePickerModel(parent *BinarySecretModel, ctx *tui.TUIContext) FilePickerModel {
	defaultPath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	secretFilePicker := filepicker.New()
	secretFilePicker.CurrentDirectory = filepath.Join(defaultPath)
	secretFilePicker.AutoHeight = false
	secretFilePicker.Height = 15

	return FilePickerModel{
		parent: parent,
		picker: secretFilePicker,
		ctx:    ctx,
	}
}

func (m FilePickerModel) Init() tea.Cmd {
	return m.picker.Init()
}

func (m *FilePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m.parent, nil
		}
	}

	var cmd tea.Cmd

	m.picker, cmd = m.picker.Update(msg)

	if selected, path := m.picker.DidSelectFile(msg); selected {
		m.parent.filePath = path
		return m.parent, nil
	}

	return m, cmd
}

func (m FilePickerModel) View() string {
	var b strings.Builder

	b.WriteString(m.picker.CurrentDirectory + "\n")
	b.WriteString(m.picker.View())

	return b.String()
}
