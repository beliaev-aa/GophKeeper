// Package handlers содержит обработчики gRPC-запросов, связанных с управлением секретами пользователей.
package handlers

import (
	"beliaev-aa/GophKeeper/internal/server/service"
	"beliaev-aa/GophKeeper/pkg/consts"
	"beliaev-aa/GophKeeper/pkg/converter"
	"beliaev-aa/GophKeeper/pkg/proto"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

// SecretHandler реализует серверные функции для управления секретами пользователей.
type SecretHandler struct {
	proto.UnimplementedSecretsServer
	logger              *zap.Logger
	secretService       service.ISecretService
	notificationHandler *NotificationHandler
}

// NewSecretHandler создаёт новый экземпляр сервера для управления секретами.
// Возвращает инициализированный экземпляр SecretHandler.
func NewSecretHandler(logger *zap.Logger, secretService service.ISecretService) *SecretHandler {
	return &SecretHandler{
		logger:              logger,
		notificationHandler: NewNotificationHandler(logger),
		secretService:       secretService,
	}
}

// SaveUserSecret сохраняет или обновляет секрет пользователя.
// Возвращает пустой ответ или ошибку, если операция не удалась.
func (s *SecretHandler) SaveUserSecret(ctx context.Context, in *proto.SaveUserSecretRequest) (*emptypb.Empty, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	secret := converter.ProtoToSecret(in.Secret)
	secret.UserID = int(userID)
	if secret.ID > 0 {
		_, err = s.secretService.UpdateSecret(ctx, secret)
	} else {
		_, err = s.secretService.CreateSecret(ctx, secret)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var clientID uint64
	clientID, err = extractClientID(ctx)
	if err != nil {
		s.logger.Error("failed to extract client ID", zap.Error(err))
		return &emptypb.Empty{}, err
	}

	isUpdated := secret.ID > 0
	err = s.notificationHandler.notifyClients(userID, clientID, secret.ID, isUpdated)
	if err != nil {
		s.logger.Error("failed to notify clients", zap.Error(err))
	}

	return &emptypb.Empty{}, nil
}

// GetUserSecret извлекает конкретный секрет пользователя.
// Возвращает секрет или ошибку, если секрет не найден или запрос не может быть выполнен.
func (s *SecretHandler) GetUserSecret(ctx context.Context, in *proto.GetUserSecretRequest) (*proto.GetUserSecretResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	secret, err := s.secretService.GetSecret(ctx, in.Id, userID)
	if err != nil {
		if errors.Is(err, fmt.Errorf("secret not found (id=%d)", in.Id)) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.GetUserSecretResponse{Secret: converter.SecretToProto(secret)}, nil
}

// GetUserSecrets извлекает все секреты пользователя.
// Возвращает список секретов или ошибку при их отсутствии или других проблемах с запросом.
func (s *SecretHandler) GetUserSecrets(ctx context.Context, _ *emptypb.Empty) (*proto.GetUserSecretsResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	secrets, err := s.secretService.GetUserSecrets(ctx, userID)
	if err != nil && !errors.Is(err, errors.New("no secrets found")) {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.GetUserSecretsResponse{Secrets: converter.SecretsToProto(secrets)}, nil
}

// DeleteUserSecret удаляет секрет пользователя.
// Возвращает пустой ответ или ошибку, если секрет не найден или не может быть удалён.
func (s *SecretHandler) DeleteUserSecret(ctx context.Context, in *proto.DeleteUserSecretRequest) (*emptypb.Empty, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = s.secretService.DeleteSecret(ctx, in.Id, userID)
	if errors.Is(err, fmt.Errorf("secret not found (id=%d)", in.Id)) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// extractUserID извлекает идентификатор пользователя из контекста запроса.
// Возвращает идентификатор пользователя или ошибку, если он не может быть извлечен.
func extractUserID(ctx context.Context) (uint64, error) {
	uid := ctx.Value(consts.CtxUserIDKey)
	userID, ok := uid.(uint64)
	if !ok {
		return 0, errors.New("failed to extract user id from context")
	}
	return userID, nil
}

// extractClientID извлекает ID клиента из метаданных контекста запроса.
// Возвращает ID клиента или ошибку, если метаданные отсутствуют или неверны.
func extractClientID(ctx context.Context) (uint64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("failed to get metadata")
	}
	values := md.Get(consts.ClientIDHeader)
	if len(values) == 0 {
		return 0, errors.New("missing client id metadata")
	}
	id, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}
