package tui

import (
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/service"
)

// This library is a nightmare... spent 80% of time making TUI :(
type TUIContext struct {
	UserService   *service.UserService
	SecretService *service.SecretService
}

func NewTUIContext(userService *service.UserService, secretService *service.SecretService) *TUIContext {
	return &TUIContext{
		UserService:   userService,
		SecretService: secretService,
	}
}
