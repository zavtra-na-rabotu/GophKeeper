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
	RegisterLoginText    = "Enter login: %s"
	RegisterPasswordText = "Enter Password: %s"
)

type RegisterModel struct {
	Login    string
	Password string
	InputPos int
}

func NewRegisterModel(inputPos int) *RegisterModel {
	return &RegisterModel{
		InputPos: inputPos,
	}
}

func (m RegisterModel) Update(app app.App, msg tea.Msg) (app.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "enter" {
			if m.InputPos == 0 {
				m.InputPos++
			} else {
				return m.registerUser(app)
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

func (m RegisterModel) View() string {
	if m.InputPos == 0 {
		return fmt.Sprintf(RegisterLoginText, m.Login)
	}
	return fmt.Sprintf(RegisterPasswordText, m.Password)
}

func (m RegisterModel) registerUser(app app.App) (app.Model, tea.Cmd) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.RegisterRequest{
		Login:    m.Login,
		Password: m.Password,
	}

	res, err := app.UserServiceClient.Register(ctx, req)
	if err != nil {
		zap.L().Error("Registration error", zap.Error(err))
		return NewMessageModel(fmt.Sprintf("Registration error %s", err)), nil
	}

	app.Token = res.Token

	return NewMessageModel(fmt.Sprintf("Registration success")), nil
}
