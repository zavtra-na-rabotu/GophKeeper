package state

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"go.uber.org/zap"
	"time"
)

const (
	EnterLoginText    = "Enter login: %s"
	EnterPasswordText = "Enter Password: %s"
)

type LoginState struct {
	Login             string
	Password          string
	InputPos          int
	UserServiceClient pb.UserServiceClient
}

func (s LoginState) Init() tea.Cmd {
	return nil
}

func (s LoginState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "enter" {
			if s.InputPos == 0 {
				s.InputPos++
			} else {
				return s.authenticateUser()
				//return InitState{Choices: Choices}, nil
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

func (s LoginState) authenticateUser() (tea.Model, tea.Cmd) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.LoginRequest{
		Login:    s.Login,
		Password: s.Password,
	}

	res, err := s.UserServiceClient.Login(ctx, req)
	if err != nil {
		zap.L().Error("Authentication error", zap.Error(err))
		return LoginState{UserServiceClient: s.UserServiceClient}, nil
	}

	fmt.Println(res)

	fmt.Println("Успешный вход!")
	return InitState{Choices: Choices}, nil
}

func (s LoginState) registerUser() (tea.Model, tea.Cmd) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.RegisterRequest{
		Login:    s.Login,
		Password: s.Password,
	}

	res, err := s.UserServiceClient.Register(ctx, req)
	if err != nil {
		zap.L().Error("Authentication error", zap.Error(err))
		return LoginState{UserServiceClient: s.UserServiceClient}, nil
	}

	fmt.Println(res)

	fmt.Println("Успешная регистрация")
	return InitState{Choices: Choices}, nil
}
