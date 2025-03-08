package model

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/style"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/model"
	"strconv"
	"strings"
)

var (
	// Стили
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("8")).Foreground(lipgloss.Color("15"))
)

type MainModel struct {
	focusIndex int
	table      [][]string
	error      string
}

func NewMainModel(ctx tui.TUIContext) MainModel {
	mainModel := MainModel{
		focusIndex: 0,
		table: [][]string{
			{"ID", "TITLE", "TYPE", "CREATED", "UPDATED"},
		},
	}

	mainModel.getSecrets(ctx)

	return mainModel
}

func (m MainModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m MainModel) Update(app tui.TUIContext, msg tea.Msg) (tui.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return NewInitModel(Choices, 0), nil
		case "up":
			if m.focusIndex > 1 {
				m.focusIndex--
			}
		case "down":
			if m.focusIndex < len(m.table)-1 {
				m.focusIndex++
			}
		}
	}
	return m, nil
}

func (m MainModel) View() string {
	var b strings.Builder

	for i, row := range m.table {
		line := fmt.Sprintf("%-4s %-10s %-10s %-20s %-20s", row[0], row[1], row[2], row[3], row[4])
		if i == 0 {
			b.WriteString(headerStyle.Render(line) + "\n")
		} else if i == m.focusIndex {
			b.WriteString(selectedStyle.Render(line) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}

	if len(m.error) > 0 {
		b.WriteString(style.ErrorStyle.Render(m.error))
	}

	return b.String()
}

func (m *MainModel) getSecrets(ctx tui.TUIContext) {
	secrets, err := ctx.SecretService.GetSecrets()
	if err != nil {
		m.error = err.Error()
		return
	}

	for _, secret := range secrets {
		secretType, err := model.ProtoToGoSecretType(secret.Type)
		if err != nil {
			m.error = err.Error()
			return
		}

		m.table = append(m.table, []string{
			strconv.FormatUint(secret.Id, 10),
			secret.Title,
			string(secretType),
			secret.CreatedAt.AsTime().Format("02 Jan 2006 15:04:05"),
			secret.UpdatedAt.AsTime().Format("02 Jan 2006 15:04:05"),
		})
	}
}
