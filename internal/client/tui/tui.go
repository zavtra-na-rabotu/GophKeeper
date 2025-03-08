package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/service"
)

// This library is a nightmare... spent 80% of time making TUI :(
type Model interface {
	Init() tea.Cmd
	Update(TUIContext, tea.Msg) (Model, tea.Cmd)
	View() string
}

type TUIContext struct {
	Model         Model
	UserService   *service.UserService
	SecretService *service.SecretService
}

func NewTUIContext(model Model, userService *service.UserService, secretService *service.SecretService) *TUIContext {
	return &TUIContext{
		Model:         model,
		UserService:   userService,
		SecretService: secretService,
	}
}

func (a TUIContext) Init() tea.Cmd {
	return nil
}

func (a TUIContext) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := a.Model.Update(a, msg)
	a.Model = newModel
	return a, cmd
}

func (a TUIContext) View() string {
	return a.Model.View()
}
