// Package service предоставляет бизнес-логику для управления секретами пользователей.
package service

import (
	"beliaev-aa/GophKeeper/internal/server/storage/repository"
	"beliaev-aa/GophKeeper/pkg/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ISecretService интерфейс для сервиса управления секретами в хранилище.
type ISecretService interface {
	GetSecret(ctx context.Context, secretID uint64, userID uint64) (*models.Secret, error)
	GetUserSecrets(ctx context.Context, userID uint64) (models.Secrets, error)
	CreateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error)
	UpdateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error)
	DeleteSecret(ctx context.Context, secretID uint64, userID uint64) error
}

// SecretService предоставляет методы для управления секретами в хранилище.
type SecretService struct {
	secretRepository *repository.SecretRepository // secretRepository является репозиторием для доступа к секретам в базе данных.
}

// NewSecretService создает новый экземпляр SecretService.
// Принимает в качестве аргумента репозиторий секретов и возвращает ссылку на сервис.
func NewSecretService(secretRepository *repository.SecretRepository) ISecretService {
	return &SecretService{secretRepository: secretRepository}
}

// GetSecret извлекает секрет по его ID и ID пользователя.
// В случае отсутствия секрета возвращает ошибку.
func (s *SecretService) GetSecret(ctx context.Context, secretID uint64, userID uint64) (*models.Secret, error) {
	secret, err := s.secretRepository.GetSecret(ctx, secretID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("secret not found (id=%d)", secretID)
	}
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// GetUserSecrets возвращает список всех секретов пользователя.
// Если секреты не найдены, возвращает ошибку.
func (s *SecretService) GetUserSecrets(ctx context.Context, userID uint64) (models.Secrets, error) {
	secrets, err := s.secretRepository.GetUserSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(secrets) == 0 {
		return nil, errors.New("no secrets found")
	}
	return secrets, nil
}

// CreateSecret создает новый секрет.
// Возвращает созданный секрет или ошибку при неудаче.
func (s *SecretService) CreateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error) {
	var err error
	secret.ID, err = s.secretRepository.Create(ctx, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}
	return secret, nil
}

// UpdateSecret обновляет существующий секрет.
// Возвращает обновленный секрет или ошибку, если секрет не найден или не удалось сохранить изменения.
func (s *SecretService) UpdateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error) {
	err := s.secretRepository.Update(ctx, secret)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("secret not found (id=%d)", secret.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to store secret: %w", err)
	}
	return secret, nil
}

// DeleteSecret удаляет секрет по его ID и ID пользователя.
// Возвращает ошибку, если удаление не произошло.
func (s *SecretService) DeleteSecret(ctx context.Context, secretID uint64, userID uint64) error {
	err := s.secretRepository.Delete(ctx, secretID, userID)
	return err
}
