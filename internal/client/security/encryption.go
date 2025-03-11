package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
)

const (
	defaultSaltLength = 32
)

type EncryptionService struct {
}

// NewEncryptionService constructor
func NewEncryptionService() *EncryptionService {
	return &EncryptionService{}
}

// DeriveKey generates a cryptographic key from a password and a salt using Argon2.
// If no salt is provided, it generates a new random salt.
func (s *EncryptionService) DeriveKey(password string, salt []byte) ([]byte, []byte, error) {
	// If no salt - generate new random one
	if len(salt) == 0 {
		salt = make([]byte, defaultSaltLength)
		_, err := rand.Read(salt)
		if err != nil {
			return nil, nil, err
		}
	}

	// Derive key AES-256
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32), salt, nil
}

// Encrypt encrypts plaintext using AES-GCM with a key derived from the provided password.
func (s *EncryptionService) Encrypt(plaintext []byte, password string) ([]byte, error) {
	key, salt, err := s.DeriveKey(password, nil)
	if err != nil {
		zap.L().Error("Failed to derive key", zap.Error(err))
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		zap.L().Error("Failed to create AES cipher", zap.Error(err))
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		zap.L().Error("Failed to create GCM", zap.Error(err))
		return nil, err
	}

	// Generate random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		zap.L().Error("Failed to generate nonce", zap.Error(err))
		return nil, err
	}

	// Encrypt data
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

// Decrypt decrypts the given ciphertext using AES-GCM and a key derived from the provided password.
func (s *EncryptionService) Decrypt(ciphertext []byte, password string) ([]byte, error) {
	// Get salt and encrypted data
	endOfData := len(ciphertext) - defaultSaltLength
	salt := ciphertext[endOfData:]
	ciphertext = ciphertext[:endOfData]

	// Derive key using password and salt
	key, salt, err := s.DeriveKey(password, salt)
	if err != nil {
		zap.L().Error("Failed to derive key", zap.Error(err))
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		zap.L().Error("Failed to create AES cipher", zap.Error(err))
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		zap.L().Error("Failed to create GCM", zap.Error(err))
		return nil, err
	}

	// Extract nonce
	nonce := ciphertext[:aesGCM.NonceSize()]

	// Extract data
	data := ciphertext[aesGCM.NonceSize():]

	// Decrypt data
	plaintext, err := aesGCM.Open(nil, nonce, data, nil)
	if err != nil {
		zap.L().Error("Failed to decrypt data", zap.Error(err))
		return nil, err
	}

	return plaintext, nil
}
