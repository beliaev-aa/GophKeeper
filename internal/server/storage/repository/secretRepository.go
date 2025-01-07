// Package repository предоставляет доступ к данным секретов, хранящихся в базе данных.
package repository

import (
	gophKeeperErrors "beliaev-aa/GophKeeper/pkg/errors"
	"beliaev-aa/GophKeeper/pkg/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// SecretRepository обеспечивает методы для работы с данными секретов в базе данных.
type SecretRepository struct {
	db *sqlx.DB
}

// NewSecretRepository создаёт новый экземпляр SecretRepository.
// Принимает подключение к базе данных sqlx.DB и возвращает указатель на SecretRepository.
func NewSecretRepository(db *sqlx.DB) *SecretRepository {
	return &SecretRepository{
		db: db,
	}
}

// GetSecret извлекает секрет по его ID и ID пользователя.
// Возвращает указатель на модель Secret или ошибку, если секрет не найден или возникла другая ошибка.
func (r SecretRepository) GetSecret(ctx context.Context, secretID uint64, userID uint64) (*models.Secret, error) {
	var secret models.Secret

	query := `SELECT * FROM secrets WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRowxContext(ctx, query, secretID, userID).StructScan(&secret)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, gophKeeperErrors.ErrNotFound
	}

	return &secret, err
}

// GetUserSecrets извлекает все секреты пользователя по его ID.
// Возвращает срез секретов или ошибку.
func (r SecretRepository) GetUserSecrets(ctx context.Context, userID uint64) (models.Secrets, error) {
	var secrets models.Secrets

	query := "SELECT * FROM secrets WHERE user_id = $1 ORDER BY updated_at DESC"
	err := r.db.SelectContext(ctx, &secrets, query, userID)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

// Create добавляет новый секрет в базу данных.
// Принимает контекст и указатель на модель Secret.
// Возвращает ID нового секрета или ошибку.
func (r SecretRepository) Create(ctx context.Context, secret *models.Secret) (uint64, error) {
	var newSecretID uint64

	query := `INSERT INTO secrets (user_id, title, metadata, secret_type, payload)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	result := r.db.QueryRowxContext(ctx, query, secret.UserID, secret.Title, secret.Metadata, secret.SecretType, secret.Payload)
	err := result.Scan(&newSecretID)
	if err != nil {
		return 0, err
	}

	return newSecretID, nil
}

// Update обновляет данные секрета в базе данных.
// Принимает контекст и указатель на модель Secret.
// Возвращает ошибку, если обновление не удалось.
func (r SecretRepository) Update(ctx context.Context, secret *models.Secret) error {
	return runInTx(r.db, func(tx *sqlx.Tx) error {
		err := tx.QueryRowxContext(ctx, "SELECT 1 FROM secrets WHERE id = $1 FOR UPDATE", secret.ID).Scan(new(int))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("secret with ID %d not found: %w", secret.ID, err)
			}
			return err
		}

		query := `UPDATE secrets SET updated_at = $1, title = $2, metadata = $3, secret_type = $4, payload = $5 WHERE id = $6;`
		_, err = tx.ExecContext(ctx, query,
			secret.UpdatedAt,
			secret.Title,
			secret.Metadata,
			secret.SecretType,
			secret.Payload,
			secret.ID,
		)

		return err
	})
}

// Delete удаляет секрет из базы данных по его ID и ID пользователя.
// Принимает контекст, ID секрета и ID пользователя.
// Возвращает ошибку, если удаление не удалось.
func (r SecretRepository) Delete(ctx context.Context, secretID uint64, userID uint64) error {
	query := `DELETE FROM secrets WHERE id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, secretID, userID)

	return err
}

// runInTx выполняет функцию fn в рамках транзакции.
// Возвращает ошибку, если транзакция не удалась.
func runInTx(db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
