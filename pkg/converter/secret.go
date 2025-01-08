// Package converter содержит методы необходимые для кодирования данных используемых в обмене между клиентом и сервером.
package converter

import (
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProtoToType конвертирует тип SecretType из protobuf в тип SecretType модели данных.
// Возвращает соответствующий тип SecretType из модели данных.
func ProtoToType(pbType proto.SecretType) models.SecretType {
	switch pbType {
	case proto.SecretType_SECRET_TYPE_CREDENTIAL:
		return models.CredSecret
	case proto.SecretType_SECRET_TYPE_TEXT:
		return models.TextSecret
	case proto.SecretType_SECRET_TYPE_BLOB:
		return models.BlobSecret
	case proto.SecretType_SECRET_TYPE_CARD:
		return models.CardSecret
	default:
		return models.UnknownSecret
	}
}

// TypeToProto конвертирует тип SecretType из модели данных в тип SecretType protobuf.
// Возвращает соответствующий тип SecretType из protobuf.
func TypeToProto(sType string) proto.SecretType {
	switch sType {
	case string(models.CredSecret):
		return proto.SecretType_SECRET_TYPE_CREDENTIAL
	case string(models.TextSecret):
		return proto.SecretType_SECRET_TYPE_TEXT
	case string(models.BlobSecret):
		return proto.SecretType_SECRET_TYPE_BLOB
	case string(models.CardSecret):
		return proto.SecretType_SECRET_TYPE_CARD
	default:
		return proto.SecretType_SECRET_TYPE_UNSPECIFIED
	}
}

// SecretToProto конвертирует объект Secret из модели данных в объект Secret protobuf.
// Возвращает новый объект Secret protobuf.
func SecretToProto(secret *models.Secret) *proto.Secret {
	return &proto.Secret{
		Id:         secret.ID,
		Title:      secret.Title,
		Metadata:   secret.Metadata,
		Payload:    secret.Payload,
		SecretType: TypeToProto(secret.SecretType),
		CreatedAt:  timestamppb.New(secret.CreatedAt),
		UpdatedAt:  timestamppb.New(secret.UpdatedAt),
	}
}

// ProtoToSecret конвертирует объект Secret из protobuf в объект Secret модели данных.
// Возвращает новый объект Secret модели данных.
func ProtoToSecret(pbSecret *proto.Secret) *models.Secret {
	return &models.Secret{
		ID:         pbSecret.Id,
		Title:      pbSecret.Title,
		Metadata:   pbSecret.Metadata,
		SecretType: string(ProtoToType(pbSecret.SecretType)),
		Payload:    pbSecret.Payload,
		CreatedAt:  pbSecret.CreatedAt.AsTime(),
		UpdatedAt:  pbSecret.UpdatedAt.AsTime(),
	}
}

// ProtoToSecrets конвертирует список объектов Secret из protobuf в список объектов Secret модели данных.
// Возвращает новый список объектов Secret модели данных.
func ProtoToSecrets(pbSecrets []*proto.Secret) []*models.Secret {
	var secrets []*models.Secret
	for _, s := range pbSecrets {
		secrets = append(secrets, ProtoToSecret(s))
	}
	return secrets
}

// SecretsToProto конвертирует список объектов Secret из модели данных в список объектов Secret protobuf.
// Возвращает новый список объектов Secret protobuf.
func SecretsToProto(secrets []*models.Secret) []*proto.Secret {
	var pbSecrets []*proto.Secret
	for _, s := range secrets {
		pbSecrets = append(pbSecrets, SecretToProto(s))
	}
	return pbSecrets
}
