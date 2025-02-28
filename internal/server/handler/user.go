package handler

import (
	"context"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := h.userService.Login(request)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "login failed %v", err)
	}
	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (h *UserHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	token, err := h.userService.Register(request)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register %v", err)
	}
	return &pb.RegisterResponse{
		Token: token,
	}, nil
}
