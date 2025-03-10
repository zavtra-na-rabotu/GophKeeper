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
	"os"
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

func (s *SecretService) CreateCredentialSecret(secretTitle string, secretLogin string, secretPassword string, secretMetadata string) error {
	if len(s.token) == 0 {
		return ErrNoToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ctx = s.createMetadata(ctx, s.token)

	credentialsProto := &pb.Credential{
		Login:    secretLogin,
		Password: secretPassword,
	}

	encryptedContent, err := s.marshalAndEncryptMessage(credentialsProto)
	if err != nil {
		return err
	}

	request := s.createSaveSecretRequest(secretTitle, pb.SecretType_SECRET_TYPE_CREDENTIAL, encryptedContent, secretMetadata)

	_, err = s.secretServiceClient.SaveSecret(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretService) CreateTextSecret(secretTitle string, secretText string, secretMetadata string) error {
	if len(s.token) == 0 {
		return ErrNoToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ctx = s.createMetadata(ctx, s.token)

	textProto := &pb.Text{
		Text: secretText,
	}

	encryptedContent, err := s.marshalAndEncryptMessage(textProto)
	if err != nil {
		return err
	}

	request := s.createSaveSecretRequest(secretTitle, pb.SecretType_SECRET_TYPE_TEXT, encryptedContent, secretMetadata)

	_, err = s.secretServiceClient.SaveSecret(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretService) CreateBinarySecret(secretTitle string, secretBinaryPath string, secretMetadata string) error {
	if len(s.token) == 0 {
		return ErrNoToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ctx = s.createMetadata(ctx, s.token)

	data, err := os.ReadFile(secretBinaryPath)
	if err != nil {
		return err
	}

	binaryProto := &pb.Binary{
		Binary: data,
	}

	encryptedContent, err := s.marshalAndEncryptMessage(binaryProto)
	if err != nil {
		return err
	}

	request := s.createSaveSecretRequest(secretTitle, pb.SecretType_SECRET_TYPE_BINARY, encryptedContent, secretMetadata)

	_, err = s.secretServiceClient.SaveSecret(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretService) CreateCardSecret(secretTitle string, cardNumber string, cardExpiryMonth string, cardExpiryYear string, cardCsc string, cardName string, secretMetadata string) error {
	if len(s.token) == 0 {
		return ErrNoToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ctx = s.createMetadata(ctx, s.token)

	binaryProto := &pb.Card{
		Number:      cardNumber,
		ExpiryMonth: cardExpiryMonth,
		ExpiryYear:  cardExpiryYear,
		Csc:         cardCsc,
		Name:        cardName,
	}

	encryptedContent, err := s.marshalAndEncryptMessage(binaryProto)
	if err != nil {
		return err
	}

	request := s.createSaveSecretRequest(secretTitle, pb.SecretType_SECRET_TYPE_CARD, encryptedContent, secretMetadata)

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

	ctx = s.createMetadata(ctx, s.token)

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

	ctx = s.createMetadata(ctx, s.token)

	_, err := s.secretServiceClient.DeleteSecret(ctx, &pb.DeleteSecretByIdRequest{Id: secretID})
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretService) DecryptAndUnmarshal(content []byte, secretType pb.SecretType) (proto.Message, error) {
	decryptedContent, err := s.encryptionService.Decrypt(content, s.password)
	if err != nil {
		return nil, err
	}

	var message proto.Message
	switch secretType {
	case pb.SecretType_SECRET_TYPE_CREDENTIAL:
		message = &pb.Credential{}
		err = proto.Unmarshal(decryptedContent, message)
	case pb.SecretType_SECRET_TYPE_TEXT:
		message = &pb.Text{}
		err = proto.Unmarshal(decryptedContent, message)
	case pb.SecretType_SECRET_TYPE_BINARY:
		message = &pb.Binary{}
		err = proto.Unmarshal(decryptedContent, message)
	case pb.SecretType_SECRET_TYPE_CARD:
		message = &pb.Card{}
		err = proto.Unmarshal(decryptedContent, message)
	}
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *SecretService) marshalAndEncryptMessage(message proto.Message) ([]byte, error) {
	marshaledContent, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	encryptedContent, err := s.encryptionService.Encrypt(marshaledContent, s.password)
	if err != nil {
		return nil, err
	}

	return encryptedContent, nil
}

func (s *SecretService) createSaveSecretRequest(secretTitle string, secretType pb.SecretType, secretContent []byte, secretMetadata string) *pb.SaveSecretRequest {
	return &pb.SaveSecretRequest{Secret: &pb.Secret{
		Title:     secretTitle,
		Type:      secretType,
		Content:   secretContent,
		Metadata:  secretMetadata,
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	},
	}
}

func (s *SecretService) createMetadata(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{"jwt": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return ctx
}
