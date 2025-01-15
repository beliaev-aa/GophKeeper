// Package repository содержит реализации интерфейсов доступа к данным.
// Он предоставляет методы для работы с данными пользователей в базе данных.
package repository

import (
	"beliaev-aa/GophKeeper/internal/server/models"
	gophKeeperErrors "beliaev-aa/GophKeeper/pkg/errors"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

// IUserRepository определяет интерфейс для репозитория пользователя,
// предоставляющего методы для работы с пользователями в базе данных.
type IUserRepository interface {
	Create(ctx context.Context, user models.User) (int, error)
	GetUserByID(ctx context.Context, ID int) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}

// UserRepository предоставляет методы для работы с пользователями в базе данных.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository создаёт новый экземпляр UserRepository.
// Функция принимает подключение к базе данных SQLX и возвращает указатель на UserRepository.
func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create регистрирует нового пользователя в системе.
// Принимает контекст выполнения и объект пользователя.
// Возвращает идентификатор нового пользователя или ошибку.
func (r *UserRepository) Create(ctx context.Context, user models.User) (int, error) {
	var newUserID int
	result := r.db.QueryRowContext(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id",
		user.Login,
		user.Password,
	)
	err := result.Scan(&newUserID)
	if err != nil {
		return 0, err
	}
	return newUserID, nil
}

// GetUserByID возвращает пользователя по его ID.
// Принимает контекст выполнения и идентификатор пользователя.
// В случае успеха возвращает объект пользователя или ошибку.
func (r *UserRepository) GetUserByID(ctx context.Context, ID int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowxContext(ctx, "SELECT id, login, created_at, password FROM users WHERE id = $1", ID).StructScan(&user)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, gophKeeperErrors.ErrNotFound
	}
	return &user, err
}

// GetUserByLogin возвращает пользователя по его логину.
// Принимает контекст выполнения и логин пользователя.
// В случае успеха возвращает объект пользователя или ошибку.
func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowxContext(ctx, "SELECT id, login, created_at, password FROM users WHERE login = $1", login).StructScan(&user)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, gophKeeperErrors.ErrNotFound
	}
	return &user, err
}
