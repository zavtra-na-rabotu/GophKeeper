package handler

import (
	"context"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SecretHandler struct {
	pb.UnimplementedSecretServiceServer
	secretService *service.SecretService
}

func NewSecretHandler(secretService *service.SecretService) *SecretHandler {
	return &SecretHandler{
		secretService: secretService,
	}
}

// SaveSecret handles the gRPC request to save a secret
func (h *SecretHandler) SaveSecret(ctx context.Context, request *pb.SaveSecretRequest) (*emptypb.Empty, error) {
	err := h.secretService.Save(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// GetSecrets handles the gRPC request to retrieve all stored secrets
func (h *SecretHandler) GetSecrets(ctx context.Context, _ *emptypb.Empty) (*pb.GetSecretsResponse, error) {
	secrets, err := h.secretService.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetSecretsResponse{Secrets: secrets}, nil
}

// DeleteSecret handles the gRPC request to delete a secret by its ID
func (h *SecretHandler) DeleteSecret(ctx context.Context, request *pb.DeleteSecretByIdRequest) (*emptypb.Empty, error) {
	err := h.secretService.DeleteSecret(ctx, request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
