// Package service содержит бизнес-логику приложения GophKeeper.
// Он включает сервисы для управления пользователями и их аутентификации.
package service

import (
	"beliaev-aa/GophKeeper/internal/server/models"
	"beliaev-aa/GophKeeper/internal/server/storage/repository"
	gophKeeperErrors "beliaev-aa/GophKeeper/pkg/errors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// ErrBadCredentials определяет ошибку, возникающую при неверных учетных данных для аутентификации.
var ErrBadCredentials = errors.New("bad auth credentials")

// UserService предоставляет методы для регистрации и аутентификации пользователей.
type UserService struct {
	userRepository repository.IUserRepository // userRepository представляет репозиторий для работы с пользователями.
}

// NewUserService создает новый экземпляр UserService с использованием заданного репозитория пользователей.
func NewUserService(userRepository repository.IUserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

// RegisterUser регистрирует нового пользователя в системе.
// Принимает контекст, логин и пароль. Возвращает зарегистрированного пользователя или ошибку.
func (s *UserService) RegisterUser(ctx context.Context, login string, password string) (*models.User, error) {
	var newUser models.User

	user, err := s.userRepository.GetUserByLogin(ctx, login)
	if err != nil && !errors.Is(err, gophKeeperErrors.ErrNotFound) {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user != nil {
		return nil, fmt.Errorf("user already exists (%s)", login)
	}

	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	newUser = models.User{Login: login, Password: hashedPassword}

	newUserID, err := s.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	newUser.ID = newUserID
	return &newUser, nil
}

// LoginUser аутентифицирует пользователя по логину и паролю.
// Возвращает пользователя или ошибку, если аутентификация не удалась.
func (s *UserService) LoginUser(ctx context.Context, login string, password string) (*models.User, error) {
	user, err := s.userRepository.GetUserByLogin(ctx, login)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBadCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate user: %w", err)
	}
	if !s.comparePassword(user.Password, password) {
		return nil, ErrBadCredentials
	}
	return user, nil
}

// hashPassword хэширует пароль с использованием bcrypt.
func (s *UserService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// comparePassword сравнивает хэшированный пароль и введенный пароль.
// Возвращает true, если пароли совпадают.
func (s *UserService) comparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
