package model

import "time"

type SecretType string

const (
	SecretTypeCredential SecretType = "credential"
	SecretTypeText       SecretType = "text"
	SecretTypeBinary     SecretType = "binary"
	SecretTypeCard       SecretType = "card"
)

type Secret struct {
	ID        int
	UserID    int
	Title     string
	Type      SecretType
	Content   []byte
	Metadata  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
