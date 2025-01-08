package interceptors

import (
	"beliaev-aa/GophKeeper/pkg/consts"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strconv"
	"testing"
)

var (
	errAuth   = errors.New("no token")
	testToken = "AccessToken"
)

func checkCtx(ctx context.Context) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return errAuth
	}

	values := md.Get(consts.AccessTokenHeader)
	if len(values) == 0 {
		return errAuth
	}

	if values[0] != testToken {
		return errAuth
	}

	return nil
}

func TestAddAuth(t *testing.T) {
	invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return checkCtx(ctx)
	}

	tests := []struct {
		name        string
		token       string
		expectedErr error
	}{
		{
			name:        "empty_token",
			token:       "",
			expectedErr: errAuth,
		},
		{
			name:        "pass_token",
			token:       testToken,
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			interceptor := AddAuth(&tc.token, 11)
			err := interceptor(context.Background(), "SomeMethod", nil, nil, nil, invoker)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func checkStreamCtx(ctx context.Context, expectedToken string, expectedClientID uint64) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok || len(md) == 0 {
		return errors.New("no metadata in context")
	}

	token := md.Get(consts.AccessTokenHeader)
	if len(token) == 0 || token[0] != expectedToken {
		return errors.New("token does not match or missing")
	}

	clientID := md.Get(consts.ClientIDHeader)
	if len(clientID) == 0 || clientID[0] != strconv.Itoa(int(expectedClientID)) {
		return errors.New("clientID does not match or missing")
	}

	return nil
}

func TestAddAuthStream(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		clientID       uint64
		expectedResult error
		expectMetadata bool
	}{
		{
			name:           "no_token_provided",
			token:          "",
			clientID:       100,
			expectedResult: nil, // Expecting streamer to be invoked without errors
			expectMetadata: false,
		},
		{
			name:           "valid_token_and_clientID",
			token:          "validToken",
			clientID:       100,
			expectedResult: nil,
			expectMetadata: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			streamerMock := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				if tc.expectMetadata {
					err := checkStreamCtx(ctx, tc.token, tc.clientID)
					if err != nil && tc.expectedResult == nil {
						return nil, err
					}
				} else {
					md, ok := metadata.FromOutgoingContext(ctx)
					if ok && len(md) > 0 {
						return nil, errors.New("unexpected metadata in context")
					}
				}
				return nil, nil
			}

			interceptor := AddAuthStream(&tc.token, tc.clientID)
			_, err := interceptor(context.Background(), &grpc.StreamDesc{}, nil, "TestMethod", streamerMock)

			if tc.expectedResult != nil {
				assert.EqualError(t, err, tc.expectedResult.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
