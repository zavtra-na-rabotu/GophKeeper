package handler

import (
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/service"
)

type SecretHandler struct {
	secretService *service.SecretService
}

func NewSecretHandler(secretService *service.SecretService) *SecretHandler {
	return &SecretHandler{
		secretService: secretService,
	}
}
