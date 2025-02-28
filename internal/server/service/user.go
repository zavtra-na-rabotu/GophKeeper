package service

import (
	"errors"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/db/repository"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/security"
	"go.uber.org/zap"
)

var ErrIncorrectLoginOrPassword = errors.New("incorrect login or password")

type UserService struct {
	userRepository *repository.UserRepository
	jwtGenerator   *security.JwtService
}

func NewUserService(userRepository *repository.UserRepository, jwtGenerator *security.JwtService) *UserService {
	return &UserService{
		userRepository: userRepository,
		jwtGenerator:   jwtGenerator,
	}
}

func (s *UserService) Login(request *pb.LoginRequest) (string, error) {
	user, err := s.userRepository.GetByLogin(request.Login)
	if err != nil {
		zap.L().Error("User not found", zap.Error(err))
		return "", err
	}

	if !security.CheckPassword(user.PasswordHash, request.Password) {
		zap.L().Error("Invalid password")
		return "", ErrIncorrectLoginOrPassword
	}

	return s.jwtGenerator.GenerateJwtToken(user.ID)
}

func (s *UserService) Register(request *pb.RegisterRequest) (string, error) {
	hash, err := security.HashPassword(request.Password)
	if err != nil {
		zap.L().Error("Failed to hash password", zap.Error(err))
		return "", err
	}

	userID, err := s.userRepository.Create(request.Login, hash)
	if err != nil {
		zap.L().Error("Failed to create user", zap.Error(err))
		return "", err
	}

	return s.jwtGenerator.GenerateJwtToken(userID)
}
