package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"testing"
)

func TestGoToProtoSecretType(t *testing.T) {
	tests := []struct {
		name string
		args SecretType
		want pb.SecretType
	}{
		{
			name: "SecretTypeCredential",
			args: SecretTypeCredential,
			want: pb.SecretType_SECRET_TYPE_CREDENTIAL,
		},
		{
			name: "SecretTypeText",
			args: SecretTypeText,
			want: pb.SecretType_SECRET_TYPE_TEXT,
		},
		{
			name: "SecretTypeBinary",
			args: SecretTypeBinary,
			want: pb.SecretType_SECRET_TYPE_BINARY,
		},
		{
			name: "SecretTypeCard",
			args: SecretTypeCard,
			want: pb.SecretType_SECRET_TYPE_CARD,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoToProtoSecretType(tt.args)

			assert.NoError(t, err)

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestProtoToGoSecretType(t *testing.T) {
	tests := []struct {
		name string
		args pb.SecretType
		want SecretType
	}{
		{
			name: "SecretTypeCredential",
			args: pb.SecretType_SECRET_TYPE_CREDENTIAL,
			want: SecretTypeCredential,
		},
		{
			name: "SecretTypeText",
			args: pb.SecretType_SECRET_TYPE_TEXT,
			want: SecretTypeText,
		},
		{
			name: "SecretTypeBinary",
			args: pb.SecretType_SECRET_TYPE_BINARY,
			want: SecretTypeBinary,
		},
		{
			name: "SecretTypeCard",
			args: pb.SecretType_SECRET_TYPE_CARD,
			want: SecretTypeCard,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProtoToGoSecretType(tt.args)

			assert.NoError(t, err)

			assert.Equal(t, got, tt.want)
		})
	}
}
