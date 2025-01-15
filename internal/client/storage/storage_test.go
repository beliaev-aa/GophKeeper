package storage

import (
	"beliaev-aa/GophKeeper/internal/client/crypto"
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/tests/mocks"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestRemoteStorage_Get_ErrorDeriveKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)
	mockClient.EXPECT().GetPassword().Return("").AnyTimes()

	_, err := NewRemoteStorage(mockClient)
	if err == nil {
		t.Error("expected error, got none")
	}
}

func TestRemoteStorage_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)

	password := "test-password"
	deriveKey, err := crypto.DeriveKey(password, "")
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	originalSecret := &models.Credentials{
		Login:    "user",
		Password: "pass",
	}
	encryptedData, err := crypto.Encrypt(`{"Login":"user","Password":"pass"}`, deriveKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	validSecret := &models.Secret{
		ID:         1,
		SecretType: string(models.CredSecret),
		Payload:    []byte(encryptedData),
		Creds:      originalSecret,
	}

	invalidPayloadSecret := &models.Secret{
		ID:         2,
		SecretType: string(models.CredSecret),
		Payload:    []byte("invalid-encrypted-data"),
	}

	mockClient.EXPECT().GetPassword().Return(password).AnyTimes()
	mockClient.EXPECT().LoadSecret(gomock.Any(), gomock.Eq(uint64(1))).Return(validSecret, nil)
	mockClient.EXPECT().LoadSecret(gomock.Any(), gomock.Eq(uint64(2))).Return(invalidPayloadSecret, nil)
	mockClient.EXPECT().LoadSecret(gomock.Any(), gomock.Eq(uint64(3))).Return(nil, fmt.Errorf("gRPC error"))

	rs, err := NewRemoteStorage(mockClient)
	if err != nil {
		t.Fatalf("Failed to create RemoteStorage: %v", err)
	}

	tests := []struct {
		name      string
		id        uint64
		expectErr bool
	}{
		{
			name:      "Get_Success",
			id:        1,
			expectErr: false,
		},
		{
			name:      "Get_Fail_InvalidPayload",
			id:        2,
			expectErr: true,
		},
		{
			name:      "Get_Fail_LoadSecretError",
			id:        3,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			secret, err := rs.Get(context.Background(), tc.id)
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err)
			}
			if tc.expectErr && secret != nil {
				t.Errorf("Expected no secret, got: %v", secret)
			}
			if !tc.expectErr && (secret == nil || secret.ID != tc.id) {
				t.Errorf("Expected secret with ID %d, got: %v", tc.id, secret)
			}
		})
	}
}

func TestRemoteStorage_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)

	password := "test-password"
	mockClient.EXPECT().GetPassword().Return(password).AnyTimes()

	deriveKey, err := crypto.DeriveKey(password, "")
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	encryptedData1, err := crypto.Encrypt(`{"Login":"user1","Password":"pass1"}`, deriveKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	encryptedData2, err := crypto.Encrypt(`{"Content":"some text"}`, deriveKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	encryptedData3, err := crypto.Encrypt(`{"FileName":"name.ext","FileBytes":"1"}`, deriveKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	encryptedData4, err := crypto.Encrypt(`{"Number":"num", "ExpYear":"2025", "ExpMonth": "01", "CVV": 666}`, deriveKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	validSecrets := []*models.Secret{
		{
			ID:         1,
			SecretType: string(models.CredSecret),
			Payload:    []byte(encryptedData1),
			Creds: &models.Credentials{
				Login:    "user1",
				Password: "pass1",
			},
		},
		{
			ID:         2,
			SecretType: string(models.TextSecret),
			Payload:    []byte(encryptedData2),
			Text: &models.Text{
				Content: "some text",
			},
		},
		{
			ID:         3,
			SecretType: string(models.BlobSecret),
			Payload:    []byte(encryptedData3),
			Blob: &models.Blob{
				FileName:  "name.ext",
				FileBytes: []byte("1"),
			},
		},
		{
			ID:         4,
			SecretType: string(models.CardSecret),
			Payload:    []byte(encryptedData4),
			Card: &models.Card{
				Number:   "num",
				ExpYear:  2025,
				ExpMonth: 01,
				CVV:      666,
			},
		},
	}

	invalidSecret := []*models.Secret{
		{
			ID:         3,
			SecretType: string(models.CredSecret),
			Payload:    []byte("invalid-encrypted-data"),
		},
	}

	rs, err := NewRemoteStorage(mockClient)
	if err != nil {
		t.Fatalf("Failed to create RemoteStorage: %v", err)
	}

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
		expected  []*models.Secret
	}{
		{
			name: "GetAll_Success",
			mockSetup: func() {
				mockClient.EXPECT().LoadSecrets(gomock.Any()).Return(validSecrets, nil).Times(1)
			},
			expectErr: false,
			expected:  validSecrets,
		},
		{
			name: "GetAll_Fail_LoadSecretsError",
			mockSetup: func() {
				mockClient.EXPECT().LoadSecrets(gomock.Any()).Return(nil, fmt.Errorf("failed to load secrets")).Times(1)
			},
			expectErr: true,
			expected:  nil,
		},
		{
			name: "GetAll_Fail_DecryptPayloadError",
			mockSetup: func() {
				mockClient.EXPECT().LoadSecrets(gomock.Any()).Return(invalidSecret, nil).Times(1)
			},
			expectErr: true,
			expected:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			secrets, err := rs.GetAll(context.Background())

			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err)
			}

			if !tc.expectErr && len(secrets) != len(tc.expected) {
				t.Errorf("Expected %d secrets, got %d", len(tc.expected), len(secrets))
			}
		})
	}
}

func TestRemoteStorage_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)

	password := "test-password"
	mockClient.EXPECT().GetPassword().Return(password).AnyTimes()
	mockClient.EXPECT().SaveSecret(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	_, err := crypto.DeriveKey(password, "")
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	testSecret := &models.Secret{
		ID:         1,
		SecretType: string(models.CredSecret),
		Creds: &models.Credentials{
			Login:    "user",
			Password: "pass",
		},
	}

	rs, err := NewRemoteStorage(mockClient)
	if err != nil {
		t.Fatalf("Failed to create RemoteStorage: %v", err)
	}

	err = rs.Create(context.Background(), testSecret)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRemoteStorage_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)

	password := "test-password"
	mockClient.EXPECT().GetPassword().Return(password).AnyTimes()
	mockClient.EXPECT().SaveSecret(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	_, err := crypto.DeriveKey(password, "")
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	testSecret := &models.Secret{
		ID:         1,
		SecretType: string(models.CredSecret),
		Creds: &models.Credentials{
			Login:    "user",
			Password: "pass",
		},
	}

	rs, err := NewRemoteStorage(mockClient)
	if err != nil {
		t.Fatalf("Failed to create RemoteStorage: %v", err)
	}

	err = rs.Update(context.Background(), testSecret)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRemoteStorage_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)

	mockClient.EXPECT().DeleteSecret(gomock.Any(), gomock.Eq(uint64(1))).Return(nil).AnyTimes()
	mockClient.EXPECT().GetPassword().Return("test-password").AnyTimes()

	rs, err := NewRemoteStorage(mockClient)
	if err != nil {
		t.Fatalf("Failed to create RemoteStorage: %v", err)
	}

	err = rs.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestEncryptPayload(t *testing.T) {
	password := "test-password"
	deriveKey, err := crypto.DeriveKey(password, "")
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	rs := &RemoteStorage{
		deriveKey: deriveKey,
	}

	tests := []struct {
		name      string
		secret    *models.Secret
		expectErr bool
	}{
		{
			name: "Encrypt_Success",
			secret: &models.Secret{
				SecretType: string(models.CredSecret),
				Creds: &models.Credentials{
					Login:    "user",
					Password: "pass",
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := rs.encryptPayload(tc.secret)
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err)
			}
			if !tc.expectErr && len(tc.secret.Payload) == 0 {
				t.Errorf("Expected encrypted payload, got empty payload")
			}
		})
	}
}

func TestDecryptPayload(t *testing.T) {
	password := "test-password"
	deriveKey, err := crypto.DeriveKey(password, "")
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	rs := &RemoteStorage{
		deriveKey: deriveKey,
	}

	secretData := &models.Credentials{
		Login:    "user",
		Password: "pass",
	}
	marshaledData, _ := json.Marshal(secretData)
	encryptedData, _ := crypto.Encrypt(string(marshaledData), deriveKey)

	tests := []struct {
		name      string
		secret    *models.Secret
		expectErr bool
	}{
		{
			name: "Decrypt_Success",
			secret: &models.Secret{
				SecretType: string(models.CredSecret),
				Payload:    []byte(encryptedData),
			},
			expectErr: false,
		},
		{
			name: "Decrypt_Fail_DecryptError",
			secret: &models.Secret{
				SecretType: string(models.CredSecret),
				Payload:    []byte("invalid-encrypted-data"),
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := rs.decryptPayload(tc.secret)
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err)
			}
			if !tc.expectErr && (tc.secret.Creds == nil || tc.secret.Creds.Login != secretData.Login || tc.secret.Creds.Password != secretData.Password) {
				t.Errorf("Decrypted data does not match expected data: %v", tc.secret.Creds)
			}
		})
	}
}

func TestString(t *testing.T) {
	rs := &RemoteStorage{}
	result := rs.String()
	if result != "remote storage" {
		t.Errorf("Expected remote storage, got: %v", result)
	}
}

func TestMarshalSecret(t *testing.T) {
	tests := []struct {
		name      string
		secret    *models.Secret
		expectErr bool
		expected  string
	}{
		{
			name: "Marshal_Credentials",
			secret: &models.Secret{
				SecretType: string(models.CredSecret),
				Creds: &models.Credentials{
					Login:    "user",
					Password: "pass",
				},
			},
			expectErr: false,
			expected:  `{"Login":"user","Password":"pass"}`,
		},
		{
			name: "Marshal_Text",
			secret: &models.Secret{
				SecretType: string(models.TextSecret),
				Text: &models.Text{
					Content: "some text",
				},
			},
			expectErr: false,
			expected:  `{"Content":"some text"}`,
		},
		{
			name: "Marshal_Card",
			secret: &models.Secret{
				SecretType: string(models.CardSecret),
				Card: &models.Card{
					Number:   "4111111111111111",
					ExpYear:  2030,
					ExpMonth: 12,
					CVV:      123,
				},
			},
			expectErr: false,
			expected:  `{"Number":"4111111111111111","ExpYear":2030,"ExpMonth":12,"CVV":123}`,
		},
		{
			name: "Marshal_Blob",
			secret: &models.Secret{
				SecretType: string(models.BlobSecret),
				Blob: &models.Blob{
					FileName:  "example.txt",
					FileBytes: []byte("example content"),
				},
			},
			expectErr: false,
			expected:  `{"FileName":"example.txt","FileBytes":"ZXhhbXBsZSBjb250ZW50"}`,
		},
		{
			name: "Marshal_UnsupportedType",
			secret: &models.Secret{
				SecretType: "unsupported_type",
			},
			expectErr: false,
			expected:  "",
		},
		{
			name: "Marshal_Credentials_Error",
			secret: &models.Secret{
				SecretType: string(models.CredSecret),
				Creds:      nil,
			},
			expectErr: false,
			expected:  "null",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := marshalSecret(tc.secret)
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err)
			}

			if !tc.expectErr {
				var actualMap map[string]interface{}
				var expectedMap map[string]interface{}

				_ = json.Unmarshal(data, &actualMap)
				_ = json.Unmarshal([]byte(tc.expected), &expectedMap)

				if len(actualMap) != len(expectedMap) {
					t.Errorf("Expected serialized data: %v, got: %v", tc.expected, string(data))
				}
			}
		})
	}
}
