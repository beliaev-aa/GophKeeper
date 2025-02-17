// Package storage обеспечивает хранение и доступ к репозиториям данных хранящихся в PostgreSQL.
package storage

import (
	"beliaev-aa/GophKeeper/internal/server/storage/migrations"
	"beliaev-aa/GophKeeper/internal/server/storage/repository"
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"time"
)

// NewStorage — создаёт новое хранилище с подключением к PostgreSQL и инициализирует схему базы данных
func NewStorage(dsn string) (*Storage, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err = migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Storage{
		UserRepository:   repository.NewUserRepository(db),
		SecretRepository: repository.NewSecretRepository(db),
	}, nil
}

// Выполняет миграции в БД
func migrate(db *sqlx.DB) error {
	goose.SetBaseFS(migrations.Migrations)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := goose.RunContext(ctx, "up", db.DB, ".")
	if err != nil {
		return err
	}

	return nil
}
