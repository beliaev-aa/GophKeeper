package interceptors

import (
	"beliaev-aa/GophKeeper/internal/server/auth"
	"beliaev-aa/GophKeeper/pkg/consts"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

type AuthTestCase struct {
	name      string
	setup     func() context.Context
	method    string
	handler   func(ctx context.Context, req any) (any, error)
	expectErr string
	expectRes any
}

func TestAuthentication(t *testing.T) {
	secretKey := "test"
	authInterceptor := Authentication([]byte(secretKey))

	handler := func(ctx context.Context, req any) (any, error) {
		var isAuth bool
		uid := ctx.Value(consts.CtxUserIDKey)

		if uid != nil {
			isAuth = true
		}

		return isAuth, nil
	}

	tests := []AuthTestCase{
		{
			name: "auth_methods_skip",
			setup: func() context.Context {
				return context.Background()
			},
			method:    "..../RegisterV1",
			handler:   handler,
			expectErr: "",
			expectRes: false,
		},
		{
			name: "valid_auth",
			setup: func() context.Context {
				userID := uint64(111)
				token, err := auth.CreateToken(int(userID), time.Now().Add(time.Hour), []byte(secretKey))
				require.NoError(t, err)

				md := metadata.New(map[string]string{
					consts.AccessTokenHeader: token,
				})

				return metadata.NewIncomingContext(context.Background(), md)
			},
			method:    "SomeMethod",
			handler:   handler,
			expectErr: "",
			expectRes: true,
		},
		{
			name: "failed_auth",
			setup: func() context.Context {
				return context.Background()
			},
			method:    "SomeMethod",
			handler:   handler,
			expectErr: "rpc error: code = Unauthenticated desc = unable to extract metadata",
			expectRes: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.setup()

			res, err := authInterceptor(ctx, nil, &grpc.UnaryServerInfo{
				FullMethod: tc.method,
			}, tc.handler)

			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectRes, res)
			}
		})
	}
}

func TestAuthContext(t *testing.T) {
	secretKey := []byte("test")

	tests := []AuthTestCase{
		{
			name: "missing_access_token",
			setup: func() context.Context {
				md := metadata.New(nil)
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectErr: "rpc error: code = Unauthenticated desc = missing access token",
		},
		{
			name: "failed_to_verify",
			setup: func() context.Context {
				md := metadata.New(map[string]string{
					consts.AccessTokenHeader: "invalid",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectErr: "rpc error: code = Unauthenticated desc = failed to verify token: token contains an invalid number of segments",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.setup()

			_, err := authContext(secretKey, ctx)

			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
