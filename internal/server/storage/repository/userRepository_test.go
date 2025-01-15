package repository

import (
	"beliaev-aa/GophKeeper/internal/server/models"
	gophKeeperErrors "beliaev-aa/GophKeeper/pkg/errors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"testing"
	"time"
)

func TestUserRepository(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		testFunc  func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name: "Create_Success",
			testFunc: func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users \(login, password\) VALUES \(\$1, \$2\) RETURNING id`).
					WithArgs("new_user", "hashed_password").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				user := models.User{
					Login:    "new_user",
					Password: "hashed_password",
				}
				id, err := repo.Create(ctx, user)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if id != 1 {
					t.Errorf("Expected ID 1, got %d", id)
				}
			},
			expectErr: false,
		},
		{
			name: "Create_Fail_DatabaseError",
			testFunc: func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users \(login, password\) VALUES \(\$1, \$2\) RETURNING id`).
					WithArgs("new_user", "hashed_password").
					WillReturnError(fmt.Errorf("database error"))

				user := models.User{
					Login:    "new_user",
					Password: "hashed_password",
				}
				_, err := repo.Create(ctx, user)
				if err == nil || err.Error() != "database error" {
					t.Errorf("Expected error 'database error', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "GetUserByID_Success",
			testFunc: func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, login, created_at, password FROM users WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "login", "created_at", "password"}).
						AddRow(1, "existing_user", time.Now(), "hashed_password"))

				user, err := repo.GetUserByID(ctx, 1)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if user.Login != "existing_user" {
					t.Errorf("Expected login 'existing_user', got %v", user.Login)
				}
			},
			expectErr: false,
		},
		{
			name: "GetUserByID_Fail_NotFound",
			testFunc: func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, login, created_at, password FROM users WHERE id = \$1`).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)

				_, err := repo.GetUserByID(ctx, 1)
				if err == nil || !errors.Is(err, gophKeeperErrors.ErrNotFound) {
					t.Errorf("Expected error 'ErrNotFound', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "GetUserByLogin_Success",
			testFunc: func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, login, created_at, password FROM users WHERE login = \$1`).
					WithArgs("existing_user").
					WillReturnRows(sqlmock.NewRows([]string{"id", "login", "created_at", "password"}).
						AddRow(1, "existing_user", time.Now(), "hashed_password"))

				user, err := repo.GetUserByLogin(ctx, "existing_user")
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if user.Login != "existing_user" {
					t.Errorf("Expected login 'existing_user', got %v", user.Login)
				}
			},
			expectErr: false,
		},
		{
			name: "GetUserByLogin_Fail_NotFound",
			testFunc: func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, login, created_at, password FROM users WHERE login = \$1`).
					WithArgs("nonexistent_user").
					WillReturnError(sql.ErrNoRows)

				_, err := repo.GetUserByLogin(ctx, "nonexistent_user")
				if err == nil || !errors.Is(err, gophKeeperErrors.ErrNotFound) {
					t.Errorf("Expected error 'ErrNotFound', got %v", err)
				}
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

			tc.testFunc(t, repo, mock)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %v", err)
			}
		})
	}
}
