package converter

import (
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/pkg/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestProtoToType(t *testing.T) {
	tests := []struct {
		name     string
		input    proto.SecretType
		expected models.SecretType
	}{
		{"Credential", proto.SecretType_SECRET_TYPE_CREDENTIAL, models.CredSecret},
		{"Text", proto.SecretType_SECRET_TYPE_TEXT, models.TextSecret},
		{"Blob", proto.SecretType_SECRET_TYPE_BLOB, models.BlobSecret},
		{"Card", proto.SecretType_SECRET_TYPE_CARD, models.CardSecret},
		{"Unknown", proto.SecretType_SECRET_TYPE_UNSPECIFIED, models.UnknownSecret},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProtoToType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTypeToProto(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected proto.SecretType
	}{
		{"Credential", string(models.CredSecret), proto.SecretType_SECRET_TYPE_CREDENTIAL},
		{"Text", string(models.TextSecret), proto.SecretType_SECRET_TYPE_TEXT},
		{"Blob", string(models.BlobSecret), proto.SecretType_SECRET_TYPE_BLOB},
		{"Card", string(models.CardSecret), proto.SecretType_SECRET_TYPE_CARD},
		{"Unknown", "unknown", proto.SecretType_SECRET_TYPE_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TypeToProto(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSecretToProto(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)
	secret := &models.Secret{
		ID:         1,
		Title:      "Test Secret",
		Metadata:   "metadata",
		Payload:    []byte("payload"),
		SecretType: string(models.CredSecret),
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	expected := &proto.Secret{
		Id:         1,
		Title:      "Test Secret",
		Metadata:   "metadata",
		Payload:    []byte("payload"),
		SecretType: proto.SecretType_SECRET_TYPE_CREDENTIAL,
		CreatedAt:  timestamppb.New(createdAt),
		UpdatedAt:  timestamppb.New(updatedAt),
	}

	result := SecretToProto(secret)
	assert.Equal(t, expected, result)
}
func TestProtoToSecret(t *testing.T) {
	createdAt := time.Now().In(time.UTC)
	updatedAt := createdAt.Add(1 * time.Hour)
	protoSecret := &proto.Secret{
		Id:         1,
		Title:      "Test Secret",
		Metadata:   "metadata",
		Payload:    []byte("payload"),
		SecretType: proto.SecretType_SECRET_TYPE_CREDENTIAL,
		CreatedAt:  timestamppb.New(createdAt),
		UpdatedAt:  timestamppb.New(updatedAt),
	}

	expected := &models.Secret{
		ID:         1,
		Title:      "Test Secret",
		Metadata:   "metadata",
		Payload:    []byte("payload"),
		SecretType: string(models.CredSecret),
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	result := ProtoToSecret(protoSecret)

	expected.CreatedAt = expected.CreatedAt.UTC()
	expected.UpdatedAt = expected.UpdatedAt.UTC()
	result.CreatedAt = result.CreatedAt.UTC()
	result.UpdatedAt = result.UpdatedAt.UTC()

	assert.Equal(t, expected, result)
}
func TestProtoToSecrets(t *testing.T) {
	createdAt := time.Now().In(time.UTC)
	updatedAt := createdAt.Add(1 * time.Hour)
	protoSecrets := []*proto.Secret{
		{
			Id:         1,
			Title:      "Test Secret",
			Metadata:   "metadata",
			Payload:    []byte("payload"),
			SecretType: proto.SecretType_SECRET_TYPE_CREDENTIAL,
			CreatedAt:  timestamppb.New(createdAt),
			UpdatedAt:  timestamppb.New(updatedAt),
		},
	}

	expected := []*models.Secret{
		{
			ID:         1,
			Title:      "Test Secret",
			Metadata:   "metadata",
			Payload:    []byte("payload"),
			SecretType: string(models.CredSecret),
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		},
	}

	result := ProtoToSecrets(protoSecrets)

	for i := range result {
		result[i].CreatedAt = result[i].CreatedAt.UTC()
		result[i].UpdatedAt = result[i].UpdatedAt.UTC()
		expected[i].CreatedAt = expected[i].CreatedAt.UTC()
		expected[i].UpdatedAt = expected[i].UpdatedAt.UTC()
	}

	assert.Equal(t, expected, result)
}

func TestSecretsToProto(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)
	secrets := []*models.Secret{
		{
			ID:         1,
			Title:      "Test Secret",
			Metadata:   "metadata",
			Payload:    []byte("payload"),
			SecretType: string(models.CredSecret),
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		},
	}

	expected := []*proto.Secret{
		{
			Id:         1,
			Title:      "Test Secret",
			Metadata:   "metadata",
			Payload:    []byte("payload"),
			SecretType: proto.SecretType_SECRET_TYPE_CREDENTIAL,
			CreatedAt:  timestamppb.New(createdAt),
			UpdatedAt:  timestamppb.New(updatedAt),
		},
	}

	result := SecretsToProto(secrets)
	assert.Equal(t, len(expected), len(result))
	for i := range result {
		assert.Equal(t, expected[i], result[i])
	}
}
