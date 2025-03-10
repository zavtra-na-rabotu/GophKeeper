package model

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/style"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/model"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"strconv"
	"strings"
)

var (
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("8")).Foreground(lipgloss.Color("15"))
)

type MainModel struct {
	focusIndex int
	table      [][]string
	secrets    map[uint64]*pb.Secret
	error      string
	ctx        *tui.TUIContext
}

func NewMainModel(ctx *tui.TUIContext) MainModel {
	mainModel := MainModel{
		focusIndex: 0,
		ctx:        ctx,
		secrets:    make(map[uint64]*pb.Secret),
	}

	mainModel.resetTable()
	mainModel.getSecrets(ctx)

	return mainModel
}

func (m MainModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return NewInitModel(m.ctx), nil
		case "up":
			if m.focusIndex > 1 {
				m.focusIndex--
			}
		case "down":
			if m.focusIndex < len(m.table)-1 {
				m.focusIndex++
			}
		case "a":
			return NewAddModel(m.ctx), nil
		case "e":
			secretID, err := strconv.ParseUint(m.table[m.focusIndex][0], 10, 64)
			if err != nil {
				m.error = err.Error()
			}
			return m.editSecret(secretID), nil
		case "d":
			secretID, err := strconv.ParseUint(m.table[m.focusIndex][0], 10, 64)
			if err != nil {
				m.error = err.Error()
			}
			m.deleteSecret(m.ctx, secretID)
			m.getSecrets(m.ctx)
			if m.focusIndex > len(m.table)-1 {
				m.focusIndex = len(m.table) - 1
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

func (m *MainModel) getSecrets(ctx *tui.TUIContext) {
	secrets, err := ctx.SecretService.GetSecrets()
	if err != nil {
		m.error = err.Error()
		return
	}

	m.resetTable()
	m.secrets = make(map[uint64]*pb.Secret)

	for _, secret := range secrets {
		secretType, err := model.ProtoToGoSecretType(secret.Type)
		if err != nil {
			m.error = err.Error()
			return
		}

		m.secrets[secret.Id] = secret

		m.table = append(m.table, []string{
			strconv.FormatUint(secret.Id, 10),
			secret.Title,
			string(secretType),
			secret.CreatedAt.AsTime().Format("02 Jan 2006 15:04:05"),
			secret.UpdatedAt.AsTime().Format("02 Jan 2006 15:04:05"),
		})
	}
}

func (m MainModel) editSecret(secretID uint64) tea.Model {
	secret, found := m.secrets[secretID]
	if !found {
		m.error = "Secret not found"
		return m
	}

	content, err := m.ctx.SecretService.DecryptAndUnmarshal(secret.Content, secret.Type)
	if err != nil {
		m.error = err.Error()
	}

	switch secret.Type {
	case pb.SecretType_SECRET_TYPE_CREDENTIAL:
		credentials, ok := content.(*pb.Credential)
		if !ok {
			m.error = "Failed to parse Credential"
			return m
		}
		return NewCredentialSecretModel(m.ctx, secret, credentials)

	case pb.SecretType_SECRET_TYPE_TEXT:
		text, ok := content.(*pb.Text)
		if !ok {
			m.error = "Failed to parse Text"
			return m
		}
		return NewTextSecretModel(m.ctx, secret, text)

	case pb.SecretType_SECRET_TYPE_BINARY:
		binary, ok := content.(*pb.Binary)
		if !ok {
			m.error = "Failed to parse Binary"
			return m
		}
		return NewBinarySecretModel(m.ctx, secret, binary)

	case pb.SecretType_SECRET_TYPE_CARD:
		card, ok := content.(*pb.Card)
		if !ok {
			m.error = "Failed to parse Card"
			return m
		}
		return NewCardSecretModel(m.ctx, secret, card)
	}

	return m
}

func (m MainModel) deleteSecret(ctx *tui.TUIContext, secretID uint64) {
	err := ctx.SecretService.DeleteSecretById(secretID)
	if err != nil {
		m.error = err.Error()
	}
}

func (m *MainModel) resetTable() {
	m.table = [][]string{
		{"ID", "TITLE", "TYPE", "CREATED", "UPDATED"},
	}
}
