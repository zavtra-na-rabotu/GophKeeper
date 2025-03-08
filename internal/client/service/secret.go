package service

import (
	"context"
	"errors"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/security"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var (
	ErrNoToken = errors.New("no token")
)

type SecretService struct {
	secretServiceClient pb.SecretServiceClient
	encryptionService   *security.EncryptionService
	password            string
	token               string
}

func NewSecretService(secretServiceClient pb.SecretServiceClient, encryptionService *security.EncryptionService) *SecretService {
	return &SecretService{
		secretServiceClient: secretServiceClient,
		encryptionService:   encryptionService,
	}
}

func (s *SecretService) SetPassword(password string) {
	s.password = password
}

func (s *SecretService) SetToken(token string) {
	s.token = token
}

func (s *SecretService) CreateTextSecret(secretTitle string, secretText string, secretMetadata string, password string) error {
	if len(s.token) == 0 {
		return ErrNoToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ctx = createMetadata(ctx, s.token)

	textProto := &pb.Text{
		Text: secretText,
	}

	serializedContent, err := proto.Marshal(textProto)
	if err != nil {
		return err
	}

	encryptedContent, err := s.encryptionService.Encrypt(serializedContent, password)
	if err != nil {
		return err
	}

	request := &pb.SaveSecretRequest{
		Secret: &pb.Secret{
			Title:     secretTitle,
			Type:      pb.SecretType_SECRET_TYPE_TEXT,
			Content:   encryptedContent,
			Metadata:  secretMetadata,
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
	}

	_, err = s.secretServiceClient.SaveSecret(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretService) GetSecrets() ([]*pb.Secret, error) {
	if len(s.token) == 0 {
		return nil, ErrNoToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ctx = createMetadata(ctx, s.token)

	res, err := s.secretServiceClient.GetSecrets(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return res.GetSecrets(), nil
}

func (s *SecretService) DeleteSecretById(secretID uint64) error {
	if len(s.token) == 0 {
		return ErrNoToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ctx = createMetadata(ctx, s.token)

	_, err := s.secretServiceClient.DeleteSecret(ctx, &pb.DeleteSecretByIdRequest{Id: secretID})
	if err != nil {
		return err
	}

	return nil
}

func createMetadata(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{"jwt": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return ctx
}
