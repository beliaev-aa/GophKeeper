// Package interceptors содержит gRPC интерцепторы, которые используются для добавления аутентификационных данных
// в метаданные запросов. Эти интерцепторы обеспечивают передачу токена доступа и идентификатора клиента
// через метаданные запросов в клиентских вызовах gRPC.
package interceptors

import (
	"beliaev-aa/GophKeeper/pkg/consts"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strconv"
)

// AddAuth возвращает UnaryClientInterceptor, который добавляет токен доступа и идентификатор клиента в метаданные запроса.
// Если токен пуст, вызов переходит к следующему обработчику без изменения контекста.
func AddAuth(token *string, clientID uint32) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if len(*token) == 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		md := metadata.New(map[string]string{
			consts.AccessTokenHeader: *token,
			consts.ClientIDHeader:    strconv.Itoa(int(clientID)),
		})

		mdCtx := metadata.NewOutgoingContext(ctx, md)
		return invoker(mdCtx, method, req, reply, cc, opts...)
	}
}

// AddAuthStream возвращает StreamClientInterceptor, который добавляет токен доступа и идентификатор клиента в метаданные запроса.
// Этот интерцептор используется для потоковых вызовов gRPC, где необходимо включить аутентификационные данные.
func AddAuthStream(token *string, clientID uint64) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if len(*token) == 0 {
			return streamer(ctx, desc, cc, method, opts...)
		}

		md := metadata.New(map[string]string{
			consts.AccessTokenHeader: *token,
			consts.ClientIDHeader:    strconv.Itoa(int(clientID)),
		})

		mdCtx := metadata.NewOutgoingContext(ctx, md)
		return streamer(mdCtx, desc, cc, method, opts...)
	}
}
