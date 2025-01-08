// Package storage предоставляет интерфейсы и реализации для управления хранилищем секретов.
// Этот пакет включает функционал для удаленного хранения данных с использованием gRPC и шифрование данных перед сохранением.
package storage

import (
	"beliaev-aa/GophKeeper/internal/client/crypto"
	"beliaev-aa/GophKeeper/internal/client/grpc"
	"beliaev-aa/GophKeeper/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// Storage описывает интерфейс для базовых операций с хранилищем секретов.
type Storage interface {
	Get(ctx context.Context, id uint64) (*models.Secret, error)
	GetAll(ctx context.Context) ([]*models.Secret, error)
	Create(ctx context.Context, secret *models.Secret) error
	Update(ctx context.Context, secret *models.Secret) error
	Delete(ctx context.Context, id uint64) error
}

// RemoteStorage реализует хранилище секретов, используя удаленный сервис через gRPC.
type RemoteStorage struct {
	client    *grpc.ClientGRPC
	deriveKey []byte
}

// NewRemoteStorage создает новый экземпляр RemoteStorage с предварительно вычисленным ключом шифрования.
func NewRemoteStorage(client *grpc.ClientGRPC) (*RemoteStorage, error) {
	deriveKey, err := crypto.DeriveKey(client.GetPassword(), strconv.FormatUint(client.GetClientID(), 10))
	if err != nil {
		return nil, err
	}
	return &RemoteStorage{
		client:    client,
		deriveKey: deriveKey,
	}, nil
}

// Get извлекает секрет по его идентификатору, расшифровывает его и возвращает.
func (store *RemoteStorage) Get(_ context.Context, id uint64) (*models.Secret, error) {
	secret, err := store.client.LoadSecret(context.Background(), id)
	if err != nil {
		return nil, err
	}

	err = store.decryptPayload(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// GetAll извлекает все секреты пользователя, расшифровывает их и возвращает.
func (store *RemoteStorage) GetAll(_ context.Context) ([]*models.Secret, error) {
	secrets, err := store.client.LoadSecrets(context.Background())
	if err != nil {
		return nil, err
	}

	for _, s := range secrets {
		err = store.decryptPayload(s)
		if err != nil {
			return nil, err
		}
	}

	return secrets, nil
}

// Create создает новый секрет в хранилище, предварительно зашифровав его.
func (store *RemoteStorage) Create(_ context.Context, secret *models.Secret) (err error) {
	err = store.encryptPayload(secret)
	if err != nil {
		return
	}

	err = store.client.SaveSecret(context.Background(), secret)
	return err
}

// Update обновляет существующий секрет, предварительно зашифровав его.
func (store *RemoteStorage) Update(_ context.Context, secret *models.Secret) (err error) {
	err = store.encryptPayload(secret)
	if err != nil {
		return
	}

	err = store.client.SaveSecret(context.Background(), secret)
	return err
}

// Delete удаляет секрет по его идентификатору.
func (store *RemoteStorage) Delete(_ context.Context, id uint64) (err error) {
	err = store.client.DeleteSecret(context.Background(), id)
	return err
}

// encryptPayload шифрует данные секрета перед сохранением.
func (store *RemoteStorage) encryptPayload(secret *models.Secret) (err error) {
	data, err := marshalSecret(secret)
	if err != nil {
		return fmt.Errorf("encryptPayload(): error serializing data: %w", err)
	}

	encryptedData, err := crypto.Encrypt(string(data), store.deriveKey)
	if err != nil {
		return fmt.Errorf("encryptPayload(): error encrypting Data: %w", err)
	} else {
		secret.Payload = []byte(encryptedData)
	}

	return err
}

// decryptPayload расшифровывает данные секрета после извлечения.
func (store *RemoteStorage) decryptPayload(secret *models.Secret) (err error) {
	decryptedData, err := crypto.Decrypt(string(secret.Payload), store.deriveKey)
	if err != nil {
		return fmt.Errorf("decryptPayload: failed to decrypt data: %w", err)

	}

	err = unmarshalSecret(secret, []byte(decryptedData))
	if err != nil {
		return fmt.Errorf("decryptPayload: failed to unmarshal data: %w", err)
	}

	return nil
}

// marshalSecret кодирует данные секрета в JSON.
func marshalSecret(secret *models.Secret) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	switch models.SecretType(secret.SecretType) {
	case models.CredSecret:
		data, err = json.Marshal(secret.Creds)
	case models.TextSecret:
		data, err = json.Marshal(secret.Text)
	case models.CardSecret:
		data, err = json.Marshal(secret.Card)
	case models.BlobSecret:
		data, err = json.Marshal(secret.Blob)
	}

	return data, err
}

// unmarshalSecret кодирует данные секрета из JSON.
func unmarshalSecret(secret *models.Secret, data []byte) error {
	var err error

	switch models.SecretType(secret.SecretType) {
	case models.CredSecret:
		err = json.Unmarshal(data, &secret.Creds)
	case models.TextSecret:
		err = json.Unmarshal(data, &secret.Text)
	case models.CardSecret:
		err = json.Unmarshal(data, &secret.Card)
	case models.BlobSecret:
		err = json.Unmarshal(data, &secret.Blob)
	}

	return err
}
