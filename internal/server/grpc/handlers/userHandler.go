// Package handlers содержит обработчики запросов gRPC для сервиса пользователей.
package handlers

import (
	"beliaev-aa/GophKeeper/internal/server/auth"
	"beliaev-aa/GophKeeper/internal/server/config"
	"beliaev-aa/GophKeeper/internal/server/service"
	"beliaev-aa/GophKeeper/pkg/proto"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// UserHandler реализует интерфейс UnimplementedUsersServer для обработки запросов пользователей.
type UserHandler struct {
	proto.UnimplementedUsersServer
	config      *config.Config
	userService *service.UserService
}

// NewUserHandler создает новый экземпляр UserHandler.
// Принимает конфигурацию сервера и сервис пользователя, возвращая инициализированный сервер пользователей.
func NewUserHandler(config *config.Config, userService *service.UserService) *UserHandler {
	return &UserHandler{
		config:      config,
		userService: userService,
	}
}

// Register регистрирует нового пользователя в системе и возвращает токен доступа.
// Принимает контекст и запрос регистрации, возвращая ответ регистрации или ошибку.
func (s *UserHandler) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user, err := s.userService.RegisterUser(ctx, in.Login, in.Password)
	if errors.Is(err, fmt.Errorf("user already exists (%s)", in.Login)) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	token, err := s.authUser(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to auth: %v", err)
	}
	return &proto.RegisterResponse{AccessToken: token}, nil
}

// Login аутентифицирует пользователя и возвращает токен доступа.
// Принимает контекст и запрос на вход, возвращая ответ входа или ошибку.
func (s *UserHandler) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	user, err := s.userService.LoginUser(ctx, in.Login, in.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	token, err := s.authUser(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to auth: %v", err)
	}
	return &proto.LoginResponse{AccessToken: token}, nil
}

// authUser генерирует токен доступа для идентифицированного пользователя.
// Возвращает строку с токеном или ошибку.
func (s *UserHandler) authUser(userID int) (string, error) {
	return auth.CreateToken(userID, time.Now().Add(time.Hour), []byte(s.config.SecretKey))
}
