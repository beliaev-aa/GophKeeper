// Package interceptors предоставляет функционал для перехвата и аутентификации запросов к gRPC-серверу.
package interceptors

import (
	"beliaev-aa/GophKeeper/internal/server/auth"
	"beliaev-aa/GophKeeper/pkg/consts"
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

// authContext производит проверку токена доступа из метаданных контекста и добавляет ID пользователя в контекст.
// Возвращает обновленный контекст и ошибку, если аутентификация не пройдена.
func authContext(secretKey []byte, ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to extract metadata")
	}

	values := md.Get(consts.AccessTokenHeader)
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing access token")
	}

	tokenText := values[0]
	tokenMap, err := auth.VerifyToken(tokenText, secretKey)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to verify token: %s", err.Error())
	}

	uid, ok := tokenMap["user_id"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user id in claims")
	}

	userID, ok := uid.(float64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid user id in claims")
	}

	ctx = context.WithValue(ctx, consts.CtxUserIDKey, uint64(userID))
	return ctx, nil
}

// Authentication создает и возвращает interceptor для серверных вызовов gRPC.
// Автоматически применяется ко всем вызовам, кроме методов регистрации и входа в систему.
func Authentication(secretKey []byte) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if strings.Contains(info.FullMethod, "Register") || strings.Contains(info.FullMethod, "Login") {
			return handler(ctx, req)
		}

		ctx, err := authContext(secretKey, ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// StreamAuthentication создаёт interceptor для потоковых серверных вызовов gRPC.
// Автоматически применяется ко всем потоковым вызовам для аутентификации пользователей
// с помощью переданного секретного ключа.
func StreamAuthentication(secretKey []byte) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, err := authContext(secretKey, ss.Context())
		if err != nil {
			return err
		}

		wrappedStream := middleware.WrapServerStream(ss)
		wrappedStream.WrappedContext = ctx

		return handler(srv, wrappedStream)
	}
}
