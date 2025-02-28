package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
)

type Model interface {
	Update(App, tea.Msg) (Model, tea.Cmd)
	View() string
}

type App struct {
	Model             Model
	UserServiceClient pb.UserServiceClient
}

func NewApp(model Model, userServiceClient pb.UserServiceClient) *App {
	return &App{
		Model:             model,
		UserServiceClient: userServiceClient,
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := a.Model.Update(a, msg)
	a.Model = newModel
	return a, cmd
}

func (a App) View() string {
	return a.Model.View()
}
