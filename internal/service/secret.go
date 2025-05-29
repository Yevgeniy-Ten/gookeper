package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"gokeeper/internal/models"
	"gokeeper/internal/repository/postgres"
)

type SecretService struct {
	repo *postgres.Postgres
	key  []byte
}

func NewSecretService(repo *postgres.Postgres, key []byte) *SecretService {
	return &SecretService{repo: repo, key: key}
}

func (s *SecretService) Create(ctx context.Context, userID int64, dto *models.CreateSecretDTO) (*models.Secret, error) {
	enc, err := encryptAES(s.key, dto.Data)
	if err != nil {
		return nil, err
	}
	copyDTO := *dto
	copyDTO.Data = enc
	return s.repo.CreateSecret(ctx, userID, &copyDTO)
}

func (s *SecretService) Get(ctx context.Context, userID, secretID int64) (*models.Secret, error) {
	sec, err := s.repo.GetSecret(ctx, userID, secretID)
	if err != nil || sec == nil {
		return nil, err
	}
	dec, err := decryptAES(s.key, sec.Data)
	if err != nil {
		return nil, err
	}
	sec.Data = dec
	return sec, nil
}

func (s *SecretService) List(ctx context.Context, userID int64) ([]*models.Secret, error) {
	secrets, err := s.repo.ListSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

func (s *SecretService) Delete(ctx context.Context, userID, secretID int64) error {
	return s.repo.DeleteSecret(ctx, userID, secretID)
}

func encryptAES(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return append(nonce, gcm.Seal(nil, nonce, plaintext, nil)...), nil
}

func decryptAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce := ciphertext[:gcm.NonceSize()]
	ct := ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ct, nil)
}
