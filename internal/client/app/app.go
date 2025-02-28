package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
)

type State interface {
	Update(App, tea.Msg) (State, tea.Cmd)
	View() string
}

type App struct {
	State             State
	UserServiceClient pb.UserServiceClient
}

func NewApp(state State, userServiceClient pb.UserServiceClient) *App {
	return &App{
		State:             state,
		UserServiceClient: userServiceClient,
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newState, cmd := a.State.Update(a, msg)
	a.State = newState
	return a, cmd
}

func (a App) View() string {
	return a.State.View()
}
