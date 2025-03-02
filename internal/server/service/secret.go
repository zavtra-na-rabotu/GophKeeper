package service

import (
	"context"
	"fmt"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/model"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/db/repository"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/interceptor"
	"go.uber.org/zap"
)

type SecretService struct {
	secretRepository *repository.SecretRepository
}

func NewSecretService(secretRepository *repository.SecretRepository) *SecretService {
	return &SecretService{
		secretRepository: secretRepository,
	}
}

func (s *SecretService) Save(ctx context.Context, request *pb.SaveSecretRequest) error {
	secret, err := model.ProtoToGoSecret(request.Secret)
	if err != nil {
		zap.L().Error("Secret request mapping failed", zap.Error(err))
		return fmt.Errorf("secret mapping failed: %w", err)
	}

	// Extract userID from context
	userID := ctx.Value(interceptor.UserIDContextKey)

	secret.UserID = userID.(uint64)

	_, err = s.secretRepository.Save(ctx, secret)
	if err != nil {
		zap.L().Error("Secret repository save failed", zap.Error(err))
		return fmt.Errorf("secret save failed: %w", err)
	}

	return nil
}
