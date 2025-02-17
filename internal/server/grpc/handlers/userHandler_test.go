package handlers

import (
	"beliaev-aa/GophKeeper/internal/server/config"
	"beliaev-aa/GophKeeper/internal/server/models"
	"beliaev-aa/GophKeeper/pkg/proto"
	"beliaev-aa/GophKeeper/tests/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIUserService(ctrl)
	cfg := &config.Config{SecretKey: "test-secret-key"}
	handler := NewUserHandler(cfg, mockService)

	tests := []struct {
		name      string
		setupMock func()
		input     *proto.RegisterRequest
		expectErr string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockService.EXPECT().RegisterUser(gomock.Any(), "new_user", "password123").Return(&models.User{ID: 1}, nil).Times(1)
			},
			input:     &proto.RegisterRequest{Login: "new_user", Password: "password123"},
			expectErr: "",
		},
		{
			name: "User_Already_Exists",
			setupMock: func() {
				mockService.EXPECT().RegisterUser(gomock.Any(), "existing_user", "password123").Return(nil, errors.New("user already exists (existing_user)")).Times(1)
			},
			input:     &proto.RegisterRequest{Login: "existing_user", Password: "password123"},
			expectErr: "rpc error: code = Internal desc = user already exists (existing_user)",
		},
		{
			name: "Internal_Error",
			setupMock: func() {
				mockService.EXPECT().RegisterUser(gomock.Any(), "new_user", "password123").Return(nil, errors.New("internal error")).Times(1)
			},
			input:     &proto.RegisterRequest{Login: "new_user", Password: "password123"},
			expectErr: "rpc error: code = Internal desc = internal error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMock != nil {
				tc.setupMock()
			}

			_, err := handler.Register(context.Background(), tc.input)

			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIUserService(ctrl)
	cfg := &config.Config{SecretKey: "test-secret-key"}
	handler := NewUserHandler(cfg, mockService)

	tests := []struct {
		name      string
		setupMock func()
		input     *proto.LoginRequest
		expectErr string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockService.EXPECT().LoginUser(gomock.Any(), "valid_user", "password123").Return(&models.User{ID: 1}, nil).Times(1)
			},
			input:     &proto.LoginRequest{Login: "valid_user", Password: "password123"},
			expectErr: "",
		},
		{
			name: "Invalid_Credentials",
			setupMock: func() {
				mockService.EXPECT().LoginUser(gomock.Any(), "invalid_user", "password123").Return(nil, errors.New("invalid credentials")).Times(1)
			},
			input:     &proto.LoginRequest{Login: "invalid_user", Password: "password123"},
			expectErr: "rpc error: code = Unauthenticated desc = invalid credentials",
		},
		{
			name: "Internal_Error",
			setupMock: func() {
				mockService.EXPECT().LoginUser(gomock.Any(), "valid_user", "password123").Return(nil, errors.New("internal error")).Times(1)
			},
			input:     &proto.LoginRequest{Login: "valid_user", Password: "password123"},
			expectErr: "rpc error: code = Unauthenticated desc = internal error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMock != nil {
				tc.setupMock()
			}

			_, err := handler.Login(context.Background(), tc.input)

			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
