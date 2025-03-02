package model

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
	LoginLoginText    = "Enter login: %s"
	LoginPasswordText = "Enter Password: %s"
)

type LoginModel struct {
	Login    string
	Password string
	InputPos int
}

func NewLoginModel(inputPos int) *LoginModel {
	return &LoginModel{
		InputPos: inputPos,
	}
}

func (m LoginModel) Update(app app.App, msg tea.Msg) (app.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "enter" {
			if m.InputPos == 0 {
				m.InputPos++
			} else {
				return m.authenticateUser(app)
			}
		} else {
			if m.InputPos == 0 {
				m.Login += key.String()
			} else {
				m.Password += key.String()
			}
		}
	}
	return m, nil
}

func (m LoginModel) View() string {
	if m.InputPos == 0 {
		return fmt.Sprintf(LoginLoginText, m.Login)
	}
	return fmt.Sprintf(LoginPasswordText, m.Password)
}

func (m LoginModel) authenticateUser(app app.App) (app.Model, tea.Cmd) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.LoginRequest{
		Login:    m.Login,
		Password: m.Password,
	}

	res, err := app.UserServiceClient.Login(ctx, req)
	if err != nil {
		zap.L().Error("Authentication error", zap.Error(err))
		return NewMessageModel(fmt.Sprintf("Authentication error %s", err)), nil
	}

	app.Token = res.Token

	return NewMessageModel("Authentication success"), nil
}
