package state

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/app"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"go.uber.org/zap"
	"time"
)

const (
	EnterLoginText    = "Enter login: %s"
	EnterPasswordText = "Enter Password: %s"
)

type LoginState struct {
	Login    string
	Password string
	InputPos int
}

func (s LoginState) Update(app app.App, msg tea.Msg) (app.State, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "enter" {
			if s.InputPos == 0 {
				s.InputPos++
			} else {
				return s.authenticateUser(app)
			}
		} else {
			if s.InputPos == 0 {
				s.Login += key.String()
			} else {
				s.Password += key.String()
			}
		}
	}
	return s, nil
}

func (s LoginState) View() string {
	if s.InputPos == 0 {
		return fmt.Sprintf(EnterLoginText, s.Login)
	}
	return fmt.Sprintf(EnterPasswordText, s.Password)
}

func (s LoginState) authenticateUser(app app.App) (app.State, tea.Cmd) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.LoginRequest{
		Login:    s.Login,
		Password: s.Password,
	}

	res, err := app.UserServiceClient.Login(ctx, req)
	if err != nil {
		zap.L().Error("Authentication error", zap.Error(err))
		return LoginState{}, nil
	}

	fmt.Println(res)

	fmt.Println("Успешный вход!")
	return InitState{Choices: Choices}, nil
}

func (s LoginState) registerUser(app app.App) (app.State, tea.Cmd) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.RegisterRequest{
		Login:    s.Login,
		Password: s.Password,
	}

	res, err := app.UserServiceClient.Register(ctx, req)
	if err != nil {
		zap.L().Error("Authentication error", zap.Error(err))
		return LoginState{}, nil
	}

	fmt.Println(res)

	fmt.Println("Успешная регистрация")
	return InitState{Choices: Choices}, nil
}
