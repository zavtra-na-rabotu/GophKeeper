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

	userID := ctx.Value(interceptor.UserIDContextKey)

	secret.UserID = userID.(uint64)

	_, err = s.secretRepository.Save(ctx, secret)
	if err != nil {
		zap.L().Error("Failed to save secret", zap.Error(err))
		return fmt.Errorf("failed to save secret: %w", err)
	}

	return nil
}

func (s *SecretService) GetAll(ctx context.Context) ([]*pb.Secret, error) {
	userID := ctx.Value(interceptor.UserIDContextKey).(uint64)

	secrets, err := s.secretRepository.GetAllByUserID(ctx, userID)
	if err != nil {
		zap.L().Error("Failed to get all user secrets", zap.Error(err))
		return nil, fmt.Errorf("failed to get all user secrets: %w", err)
	}

	var protoSecrets []*pb.Secret
	for _, secret := range secrets {
		protoSecret, err := model.GoToProtoSecret(secret)
		if err != nil {
			zap.L().Error("Failed to convert secret to proto", zap.Error(err))
			return nil, fmt.Errorf("failed to convert secret to proto: %w", err)
		}
		protoSecrets = append(protoSecrets, protoSecret)
	}

	return protoSecrets, nil
}
