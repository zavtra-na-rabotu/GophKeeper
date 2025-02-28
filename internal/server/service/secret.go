package service

import (
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/db/repository"
)

type SecretService struct {
	secretRepository *repository.SecretRepository
}

func NewSecretService(secretRepository *repository.SecretRepository) *SecretService {
	return &SecretService{
		secretRepository: secretRepository,
	}
}
