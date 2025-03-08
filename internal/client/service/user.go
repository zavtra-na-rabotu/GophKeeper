package service

import (
	"context"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"time"
)

type UserService struct {
	UserServiceClient pb.UserServiceClient
}

func NewUserService(userServiceClient pb.UserServiceClient) *UserService {
	return &UserService{
		UserServiceClient: userServiceClient,
	}
}

func (s *UserService) Login(login string, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.LoginRequest{
		Login:    login,
		Password: password,
	}

	res, err := s.UserServiceClient.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Token, nil
}

func (s *UserService) Register(login string, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.RegisterRequest{
		Login:    login,
		Password: password,
	}

	res, err := s.UserServiceClient.Register(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Token, nil
}
