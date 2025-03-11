package security

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeriveKeyWithSaltProvided(t *testing.T) {
	service := NewEncryptionService()
	password := "testpassword"

	salt := make([]byte, defaultSaltLength)
	for i := range salt {
		salt[i] = byte(i)
	}

	key, returnedSalt, err := service.DeriveKey(password, salt)
	require.NoError(t, err)
	assert.Len(t, key, 32)
	assert.Equal(t, salt, returnedSalt)
}

func TestDeriveKeyWithoutSalt(t *testing.T) {
	service := NewEncryptionService()
	password := "testpassword"

	key, salt, err := service.DeriveKey(password, nil)

	require.NoError(t, err)
	assert.Len(t, key, 32)
	assert.Len(t, salt, 32)
}

func TestEncryptDecrypt(t *testing.T) {
	service := NewEncryptionService()
	password := "securepassword"
	plaintext := []byte("Hello, world!")

	ciphertext, err := service.Encrypt(plaintext, password)
	require.NoError(t, err)
	assert.NotEmpty(t, ciphertext)

	decrypted, err := service.Decrypt(ciphertext, password)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestDecryptWithWrongPassword(t *testing.T) {
	s := NewEncryptionService()
	password := "securepassword"
	wrongPassword := "wrongpassword"
	plaintext := []byte("Hello, world!")

	ciphertext, err := s.Encrypt(plaintext, password)
	require.NoError(t, err)

	decrypted, err := s.Decrypt(ciphertext, wrongPassword)
	assert.Error(t, err)
	assert.Nil(t, decrypted)
}

func TestEncryptUniqueCiphertext(t *testing.T) {
	s := NewEncryptionService()
	password := "securepassword"
	plaintext := []byte("Hello, world!")

	ciphertext1, err := s.Encrypt(plaintext, password)
	require.NoError(t, err)

	ciphertext2, err := s.Encrypt(plaintext, password)
	require.NoError(t, err)

	assert.NotEqual(t, ciphertext1, ciphertext2)
}
