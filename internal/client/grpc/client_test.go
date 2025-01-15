package grpc

import (
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/pkg/proto"
	"beliaev-aa/GophKeeper/tests/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClientGRPC_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersClient := mocks.NewMockUsersClient(ctrl)
	client := &ClientGRPC{
		UsersClient: mockUsersClient,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr error
	}{
		{
			name: "Login_Success",
			setupMock: func() {
				req := &proto.LoginRequest{
					Login:    "test",
					Password: "1234",
				}
				resp := &proto.LoginResponse{
					AccessToken: "access_token",
				}
				mockUsersClient.EXPECT().Login(gomock.Any(), req).Return(resp, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Login_Failed_Unavailable",
			setupMock: func() {
				req := &proto.LoginRequest{
					Login:    "test",
					Password: "1234",
				}
				mockUsersClient.EXPECT().Login(gomock.Any(), req).Return(nil, status.Error(codes.Unavailable, "server unavailable"))
			},
			expectedErr: errors.New("server unavailable"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			_, err := client.Login(context.Background(), "test", "1234")
			if (err != nil && tc.expectedErr == nil) || (err == nil && tc.expectedErr != nil) || (err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestClientGRPC_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersClient := mocks.NewMockUsersClient(ctrl)
	client := &ClientGRPC{
		UsersClient: mockUsersClient,
	}

	tests := []struct {
		name          string
		login         string
		password      string
		setupMock     func(login, password string)
		expectedErr   error
		expectedToken string
	}{
		{
			name:     "Register_Success",
			login:    "new_user",
			password: "password123",
			setupMock: func(login, password string) {
				req := &proto.RegisterRequest{
					Login: login,

					Password: password,
				}
				resp := &proto.RegisterResponse{
					AccessToken: "new_access_token",
				}
				mockUsersClient.EXPECT().Register(gomock.Any(), gomock.Eq(req)).Return(resp, nil)
			},
			expectedErr:   nil,
			expectedToken: "new_access_token",
		},
		{
			name:     "Register_Failed_AlreadyExists",
			login:    "existing_user",
			password: "password123",
			setupMock: func(login, password string) {
				req := &proto.RegisterRequest{
					Login:    login,
					Password: password,
				}
				mockUsersClient.EXPECT().Register(gomock.Any(), gomock.Eq(req)).Return(nil, status.Error(codes.AlreadyExists, "user already exists"))
			},
			expectedErr:   errors.New("user already exists"),
			expectedToken: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock(tc.login, tc.password)
			token, err := client.Register(context.Background(), tc.login, tc.password)
			if (err != nil && tc.expectedErr == nil) || (err == nil && tc.expectedErr != nil) || (err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
			if token != tc.expectedToken {
				t.Errorf("Expected token: %s, got: %s", tc.expectedToken, token)
			}
		})
	}
}

func TestClientGRPC_LoadSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSecretsClient := mocks.NewMockSecretsClient(ctrl)
	client := &ClientGRPC{
		SecretsClient: mockSecretsClient,
	}

	now := time.Now()
	testSecrets := []*models.Secret{
		{ID: 1, Title: "Test Secret", Metadata: "Test Metadata", SecretType: "text", Payload: []byte("Test Payload"), CreatedAt: now, UpdatedAt: now},
	}
	protoSecrets := []*proto.Secret{
		{Id: 1, Title: "Test Secret", Metadata: "Test Metadata", SecretType: proto.SecretType_SECRET_TYPE_TEXT, Payload: []byte("Test Payload"), CreatedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)},
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr string
		expected    []*models.Secret
	}{
		{
			name: "Load_Secrets_Success",
			setupMock: func() {
				resp := &proto.GetUserSecretsResponse{Secrets: protoSecrets}
				mockSecretsClient.EXPECT().GetUserSecrets(gomock.Any(), &emptypb.Empty{}).Return(resp, nil)
			},
			expectedErr: "",
			expected:    testSecrets,
		},
		{
			name: "Load_Secrets_Failed",
			setupMock: func() {
				mockSecretsClient.EXPECT().GetUserSecrets(gomock.Any(), &emptypb.Empty{}).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedErr: "internal error",
			expected:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			secrets, err := client.LoadSecrets(context.Background())
			if !compareSecrets(secrets, tc.expected) || !compareErrors(err, tc.expectedErr) {
				t.Errorf("LoadSecrets() got secrets = %v, err = %v, want secrets = %v, err = %v", secrets, err, tc.expected, tc.expectedErr)
			}
		})
	}
}

func compareSecrets(got, want []*models.Secret) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i].ID != want[i].ID ||
			got[i].Title != want[i].Title ||
			got[i].Metadata != want[i].Metadata ||
			got[i].SecretType != want[i].SecretType ||
			!bytes.Equal(got[i].Payload, want[i].Payload) ||
			!got[i].CreatedAt.Equal(want[i].CreatedAt) ||
			!got[i].UpdatedAt.Equal(want[i].UpdatedAt) {
			return false
		}
	}
	return true
}

func compareErrors(err error, expectedErr string) bool {
	if err == nil && expectedErr == "" {
		return true
	}
	if err == nil || expectedErr == "" {
		return false
	}

	return strings.Contains(err.Error(), expectedErr)
}

func TestClientGRPC_LoadSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSecretsClient := mocks.NewMockSecretsClient(ctrl)
	client := &ClientGRPC{
		SecretsClient: mockSecretsClient,
	}

	fixedTime := time.Date(2025, time.January, 15, 10, 49, 5, 278598, time.UTC)

	testSecret := &models.Secret{
		ID:         1,
		Title:      "Test Secret",
		Metadata:   "Test Metadata",
		SecretType: "text",
		Payload:    []byte("Test Payload"),
		CreatedAt:  fixedTime,
		UpdatedAt:  fixedTime,
	}

	protoSecret := &proto.Secret{
		Id:         1,
		Title:      "Test Secret",
		Metadata:   "Test Metadata",
		SecretType: proto.SecretType_SECRET_TYPE_TEXT,
		Payload:    []byte("Test Payload"),
		CreatedAt:  timestamppb.New(fixedTime),
		UpdatedAt:  timestamppb.New(fixedTime),
	}

	tests := []struct {
		name        string
		id          uint64
		setupMock   func(id uint64)
		expectedErr string
		expected    *models.Secret
	}{
		{
			name: "Load_Secret_Success",
			id:   1,
			setupMock: func(id uint64) {
				req := &proto.GetUserSecretRequest{Id: id}
				resp := &proto.GetUserSecretResponse{
					Secret: protoSecret,
				}
				mockSecretsClient.EXPECT().GetUserSecret(gomock.Any(), req).Return(resp, nil)
			},
			expectedErr: "",
			expected:    testSecret,
		},
		{
			name: "Load_Secret_Failed",
			id:   1,
			setupMock: func(id uint64) {
				req := &proto.GetUserSecretRequest{Id: id}
				mockSecretsClient.EXPECT().GetUserSecret(gomock.Any(), req).Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedErr: "not found",
			expected:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock(tc.id)
			secret, err := client.LoadSecret(context.Background(), tc.id)
			if !reflect.DeepEqual(secret, tc.expected) || !compareErrors(err, tc.expectedErr) {
				t.Errorf("LoadSecret() got secret = %v, err = %v, want secret = %v, err = %v", secret, err, tc.expected, tc.expectedErr)
			}
		})
	}
}

func TestClientGRPC_SaveSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSecretsClient := mocks.NewMockSecretsClient(ctrl)
	client := &ClientGRPC{
		SecretsClient: mockSecretsClient,
	}

	fixedTime := time.Date(2025, time.January, 15, 10, 49, 5, 278598, time.UTC)

	testSecret := &models.Secret{
		ID:         1,
		Title:      "Test Secret",
		Metadata:   "Test Metadata",
		SecretType: "text",
		Payload:    []byte("Test Payload"),
		CreatedAt:  fixedTime,
		UpdatedAt:  fixedTime,
	}

	protoSecret := &proto.Secret{
		Id:         1,
		Title:      "Test Secret",
		Metadata:   "Test Metadata",
		SecretType: proto.SecretType_SECRET_TYPE_TEXT,
		Payload:    []byte("Test Payload"),
		CreatedAt:  timestamppb.New(fixedTime),
		UpdatedAt:  timestamppb.New(fixedTime),
	}

	tests := []struct {
		name        string
		secret      *models.Secret
		setupMock   func(secret *models.Secret)
		expectedErr string
	}{
		{
			name:   "Save_Secret_Success",
			secret: testSecret,
			setupMock: func(secret *models.Secret) {
				req := &proto.SaveUserSecretRequest{Secret: protoSecret}
				mockSecretsClient.EXPECT().SaveUserSecret(gomock.Any(), req).Return(&emptypb.Empty{}, nil)
			},
			expectedErr: "",
		},
		{
			name:   "Save_Secret_Failed",
			secret: testSecret,
			setupMock: func(secret *models.Secret) {
				req := &proto.SaveUserSecretRequest{Secret: protoSecret}
				mockSecretsClient.EXPECT().SaveUserSecret(gomock.Any(), req).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedErr: "internal error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock(tc.secret)
			err := client.SaveSecret(context.Background(), tc.secret)
			if !compareErrors(err, tc.expectedErr) {
				t.Errorf("SaveSecret() got err = %v, want err = %v", err, tc.expectedErr)
			}
		})
	}
}

func TestClientGRPC_DeleteSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSecretsClient := mocks.NewMockSecretsClient(ctrl)
	client := &ClientGRPC{
		SecretsClient: mockSecretsClient,
	}

	tests := []struct {
		name        string
		id          uint64
		setupMock   func(id uint64)
		expectedErr string
	}{
		{
			name: "Delete_Secret_Success",
			id:   1,
			setupMock: func(id uint64) {
				req := &proto.DeleteUserSecretRequest{Id: id}
				mockSecretsClient.EXPECT().DeleteUserSecret(gomock.Any(), req).Return(&emptypb.Empty{}, nil)
			},
			expectedErr: "",
		},
		{
			name: "Delete_Secret_NotFound",
			id:   1,
			setupMock: func(id uint64) {
				req := &proto.DeleteUserSecretRequest{Id: id}
				mockSecretsClient.EXPECT().DeleteUserSecret(gomock.Any(), req).Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedErr: "not found",
		},
		{
			name: "Delete_Secret_InternalError",
			id:   1,
			setupMock: func(id uint64) {
				req := &proto.DeleteUserSecretRequest{Id: id}
				mockSecretsClient.EXPECT().DeleteUserSecret(gomock.Any(), req).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedErr: "internal error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock(tc.id)
			err := client.DeleteSecret(context.Background(), tc.id)
			if !compareErrors(err, tc.expectedErr) {
				t.Errorf("DeleteSecret() got err = %v, want err = %v", err, tc.expectedErr)
			}
		})
	}
}

func TestTokenAndPasswordSetGet(t *testing.T) {
	client := &ClientGRPC{}

	token := "test-token"
	client.SetToken(token)
	if got := client.GetToken(); got != token {
		t.Errorf("GetToken() = %v, want %v", got, token)
	}

	password := "test-password"
	client.SetPassword(password)
	if got := client.GetPassword(); got != password {
		t.Errorf("GetPassword() = %v, want %v", got, password)
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
		name        string
		inputError  error
		expectedErr error
	}{
		{
			name:        "No_error",
			inputError:  nil,
			expectedErr: nil,
		},
		{
			name:        "Non_gRPC_error",
			inputError:  errors.New("regular error"),
			expectedErr: errors.New("regular error"),
		},
		{
			name:        "gRPC_Unavailable_error",
			inputError:  status.Error(codes.Unavailable, "service unavailable"),
			expectedErr: errors.New("server unavailable"),
		},
		{
			name:        "gRPC_Unauthenticated_error",
			inputError:  status.Error(codes.Unauthenticated, "authentication failed"),
			expectedErr: errors.New("failed to authenticate"),
		},
		{
			name:        "gRPC_AlreadyExists_error",
			inputError:  status.Error(codes.AlreadyExists, "item already exists"),
			expectedErr: errors.New("user already exists"),
		},
		{
			name:        "gRPC_Internal_error",
			inputError:  status.Error(codes.Internal, "internal server error"),
			expectedErr: status.Error(codes.Internal, "internal server error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotErr := parseError(tc.inputError)
			if !errorEqual(gotErr, tc.expectedErr) {
				t.Errorf("parseError(%v) = %v, want %v", tc.inputError, gotErr, tc.expectedErr)
			}
		})
	}
}

func errorEqual(a, b error) bool {
	if errors.Is(a, b) {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Error() == b.Error()
}
