package model

import (
	"fmt"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type SecretType string

const (
	SecretTypeCredential SecretType = "credential"
	SecretTypeText       SecretType = "text"
	SecretTypeBinary     SecretType = "binary"
	SecretTypeCard       SecretType = "card"
)

type Secret struct {
	ID        uint64     `db:"id"`
	UserID    uint64     `db:"user_id"`
	Title     string     `db:"title"`
	Type      SecretType `db:"type"`
	Content   []byte     `db:"content"`
	Metadata  string     `db:"metadata"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}

func ProtoToGoSecret(secret *pb.Secret) (*Secret, error) {
	secretType, err := ProtoToGoSecretType(secret.Type)
	if err != nil {
		zap.L().Error("Failed to map SecretType from Proto to Go enum", zap.Error(err))
		return nil, err
	}

	return &Secret{
		ID:        secret.Id,
		Title:     secret.Title,
		Type:      secretType,
		Content:   secret.Content,
		Metadata:  secret.Metadata,
		CreatedAt: secret.CreatedAt.AsTime(),
		UpdatedAt: secret.UpdatedAt.AsTime(),
	}, nil
}

func GoToProtoSecret(secret *Secret) (*pb.Secret, error) {
	secretType, err := GoToProtoSecretType(secret.Type)
	if err != nil {
		zap.L().Error("Failed to map SecretType from Go to Proto enum", zap.Error(err))
		return nil, err
	}

	return &pb.Secret{
		Id:        secret.ID,
		Title:     secret.Title,
		Type:      secretType,
		Content:   secret.Content,
		Metadata:  secret.Metadata,
		CreatedAt: timestamppb.New(secret.CreatedAt),
		UpdatedAt: timestamppb.New(secret.UpdatedAt),
	}, nil
}

func ProtoToGoSecretType(protoType pb.SecretType) (SecretType, error) {
	switch protoType {
	case pb.SecretType_SECRET_TYPE_CREDENTIAL:
		return SecretTypeCredential, nil
	case pb.SecretType_SECRET_TYPE_TEXT:
		return SecretTypeText, nil
	case pb.SecretType_SECRET_TYPE_BINARY:
		return SecretTypeBinary, nil
	case pb.SecretType_SECRET_TYPE_CARD:
		return SecretTypeCard, nil
	default:
		return "", fmt.Errorf("unknown SecretType: %v", protoType)
	}
}

func GoToProtoSecretType(goType SecretType) (pb.SecretType, error) {
	switch goType {
	case SecretTypeCredential:
		return pb.SecretType_SECRET_TYPE_CREDENTIAL, nil
	case SecretTypeText:
		return pb.SecretType_SECRET_TYPE_TEXT, nil
	case SecretTypeBinary:
		return pb.SecretType_SECRET_TYPE_BINARY, nil
	case SecretTypeCard:
		return pb.SecretType_SECRET_TYPE_CARD, nil
	default:
		return pb.SecretType(0), fmt.Errorf("unknown SecretType: %s", goType)
	}
}
