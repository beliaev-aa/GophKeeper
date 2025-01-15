package handlers

import (
	"beliaev-aa/GophKeeper/pkg/proto"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"testing"
)

type mockStream struct {
	mock.Mock
}

func (m *mockStream) Send(resp *proto.SubscribeResponse) error {
	args := m.Called(resp)
	return args.Error(0)
}

func (m *mockStream) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func (m *mockStream) SetHeader(_ metadata.MD) error {
	return nil
}

func (m *mockStream) SendHeader(_ metadata.MD) error {
	return nil
}

func (m *mockStream) SetTrailer(_ metadata.MD) {}

func (m *mockStream) SendMsg(_ any) error {
	return nil
}

func (m *mockStream) RecvMsg(_ any) error {
	return nil
}

func TestNotificationHandler_Subscribe(t *testing.T) {
	logger := zap.NewNop()
	handler := NewNotificationHandler(logger)

	tests := []struct {
		name      string
		setup     func(stream *mockStream)
		input     *proto.SubscribeRequest
		expectErr string
	}{
		{
			name: "Subscribe_Success",
			setup: func(stream *mockStream) {
				ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
					"user_id": "123",
				}))
				stream.On("Context").Return(ctx)
			},
			input: &proto.SubscribeRequest{
				Id: 1,
			},
			expectErr: "rpc error: code = Internal desc = failed to extract user id from context",
		},
		{
			name: "Subscribe_MissingUserID",
			setup: func(stream *mockStream) {
				ctx := context.Background()
				stream.On("Context").Return(ctx)
			},
			input: &proto.SubscribeRequest{
				Id: 1,
			},
			expectErr: "rpc error: code = Internal desc = failed to extract user id from context",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stream := new(mockStream)
			tc.setup(stream)

			err := handler.Subscribe(tc.input, stream)
			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNotificationHandler_NotifyClients(t *testing.T) {
	logger := zap.NewNop()
	handler := NewNotificationHandler(logger)

	tests := []struct {
		name      string
		setup     func(handler *NotificationHandler)
		userID    uint64
		clientID  uint64
		ID        uint64
		updated   bool
		expectErr string
	}{
		{
			name: "NotifyClients_Success",
			setup: func(handler *NotificationHandler) {
				stream := new(mockStream)
				stream.On("Send", &proto.SubscribeResponse{Id: 1, Updated: true}).Return(nil)

				handler.subscribers.Store(uint64(123), []subscriber{
					{
						stream:   stream,
						id:       2,
						finished: make(chan bool),
					},
				})
			},
			userID:    123,
			clientID:  2,
			ID:        1,
			updated:   true,
			expectErr: "",
		},
		{
			name: "NotifyClients_NoSubscribers",
			setup: func(handler *NotificationHandler) {
				// Оставляем пустую карту подписчиков.
			},
			userID:    123,
			clientID:  2,
			ID:        1,
			updated:   true,
			expectErr: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(handler)
			}

			err := handler.notifyClients(tc.userID, tc.clientID, tc.ID, tc.updated)
			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
