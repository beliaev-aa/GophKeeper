// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.29.2
// source: secrets.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SecretType int32

const (
	SecretType_SECRET_TYPE_UNSPECIFIED SecretType = 0
	SecretType_SECRET_TYPE_CREDENTIAL  SecretType = 1
	SecretType_SECRET_TYPE_TEXT        SecretType = 2
	SecretType_SECRET_TYPE_BLOB        SecretType = 3
	SecretType_SECRET_TYPE_CARD        SecretType = 4
)

// Enum value maps for SecretType.
var (
	SecretType_name = map[int32]string{
		0: "SECRET_TYPE_UNSPECIFIED",
		1: "SECRET_TYPE_CREDENTIAL",
		2: "SECRET_TYPE_TEXT",
		3: "SECRET_TYPE_BLOB",
		4: "SECRET_TYPE_CARD",
	}
	SecretType_value = map[string]int32{
		"SECRET_TYPE_UNSPECIFIED": 0,
		"SECRET_TYPE_CREDENTIAL":  1,
		"SECRET_TYPE_TEXT":        2,
		"SECRET_TYPE_BLOB":        3,
		"SECRET_TYPE_CARD":        4,
	}
)

func (x SecretType) Enum() *SecretType {
	p := new(SecretType)
	*p = x
	return p
}

func (x SecretType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SecretType) Descriptor() protoreflect.EnumDescriptor {
	return file_secrets_proto_enumTypes[0].Descriptor()
}

func (SecretType) Type() protoreflect.EnumType {
	return &file_secrets_proto_enumTypes[0]
}

func (x SecretType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SecretType.Descriptor instead.
func (SecretType) EnumDescriptor() ([]byte, []int) {
	return file_secrets_proto_rawDescGZIP(), []int{0}
}

type Secret struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title      string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Metadata   string                 `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Payload    []byte                 `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
	SecretType SecretType             `protobuf:"varint,5,opt,name=secret_type,json=secretType,proto3,enum=proto.SecretType" json:"secret_type,omitempty"`
	CreatedAt  *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt  *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Secret) Reset() {
	*x = Secret{}
	mi := &file_secrets_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Secret) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Secret) ProtoMessage() {}

func (x *Secret) ProtoReflect() protoreflect.Message {
	mi := &file_secrets_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Secret.ProtoReflect.Descriptor instead.
func (*Secret) Descriptor() ([]byte, []int) {
	return file_secrets_proto_rawDescGZIP(), []int{0}
}

func (x *Secret) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Secret) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Secret) GetMetadata() string {
	if x != nil {
		return x.Metadata
	}
	return ""
}

func (x *Secret) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *Secret) GetSecretType() SecretType {
	if x != nil {
		return x.SecretType
	}
	return SecretType_SECRET_TYPE_UNSPECIFIED
}

func (x *Secret) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Secret) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type GetUserSecretsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Secrets []*Secret `protobuf:"bytes,1,rep,name=secrets,proto3" json:"secrets,omitempty"`
}

func (x *GetUserSecretsResponse) Reset() {
	*x = GetUserSecretsResponse{}
	mi := &file_secrets_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserSecretsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserSecretsResponse) ProtoMessage() {}

func (x *GetUserSecretsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_secrets_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserSecretsResponse.ProtoReflect.Descriptor instead.
func (*GetUserSecretsResponse) Descriptor() ([]byte, []int) {
	return file_secrets_proto_rawDescGZIP(), []int{1}
}

func (x *GetUserSecretsResponse) GetSecrets() []*Secret {
	if x != nil {
		return x.Secrets
	}
	return nil
}

type GetUserSecretRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetUserSecretRequest) Reset() {
	*x = GetUserSecretRequest{}
	mi := &file_secrets_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserSecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserSecretRequest) ProtoMessage() {}

func (x *GetUserSecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_secrets_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserSecretRequest.ProtoReflect.Descriptor instead.
func (*GetUserSecretRequest) Descriptor() ([]byte, []int) {
	return file_secrets_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserSecretRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetUserSecretResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Secret *Secret `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"`
}

func (x *GetUserSecretResponse) Reset() {
	*x = GetUserSecretResponse{}
	mi := &file_secrets_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserSecretResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserSecretResponse) ProtoMessage() {}

func (x *GetUserSecretResponse) ProtoReflect() protoreflect.Message {
	mi := &file_secrets_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserSecretResponse.ProtoReflect.Descriptor instead.
func (*GetUserSecretResponse) Descriptor() ([]byte, []int) {
	return file_secrets_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserSecretResponse) GetSecret() *Secret {
	if x != nil {
		return x.Secret
	}
	return nil
}

type SaveUserSecretRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Secret *Secret `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"`
}

func (x *SaveUserSecretRequest) Reset() {
	*x = SaveUserSecretRequest{}
	mi := &file_secrets_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SaveUserSecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveUserSecretRequest) ProtoMessage() {}

func (x *SaveUserSecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_secrets_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveUserSecretRequest.ProtoReflect.Descriptor instead.
func (*SaveUserSecretRequest) Descriptor() ([]byte, []int) {
	return file_secrets_proto_rawDescGZIP(), []int{4}
}

func (x *SaveUserSecretRequest) GetSecret() *Secret {
	if x != nil {
		return x.Secret
	}
	return nil
}

type DeleteUserSecretRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteUserSecretRequest) Reset() {
	*x = DeleteUserSecretRequest{}
	mi := &file_secrets_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteUserSecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserSecretRequest) ProtoMessage() {}

func (x *DeleteUserSecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_secrets_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserSecretRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserSecretRequest) Descriptor() ([]byte, []int) {
	return file_secrets_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteUserSecretRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_secrets_proto protoreflect.FileDescriptor

var file_secrets_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8e, 0x02, 0x0a, 0x06, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x32, 0x0a, 0x0b, 0x73,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x0a, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x41, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x27, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52,
	0x07, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x22, 0x26, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x3e, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x06, 0x73, 0x65, 0x63,
	0x72, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x22, 0x3e, 0x0a, 0x15, 0x53, 0x61, 0x76, 0x65, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x06, 0x73, 0x65, 0x63,
	0x72, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x22, 0x29, 0x0a, 0x17, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65,
	0x63, 0x72, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x2a, 0x87, 0x01, 0x0a, 0x0a,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x17, 0x53, 0x45,
	0x43, 0x52, 0x45, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x43, 0x52, 0x45,
	0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x52, 0x45, 0x44, 0x45, 0x4e, 0x54, 0x49, 0x41,
	0x4c, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x45, 0x43, 0x52, 0x45, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x54, 0x45, 0x58, 0x54, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x45, 0x43,
	0x52, 0x45, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x42, 0x4c, 0x4f, 0x42, 0x10, 0x03, 0x12,
	0x14, 0x0a, 0x10, 0x53, 0x45, 0x43, 0x52, 0x45, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43,
	0x41, 0x52, 0x44, 0x10, 0x04, 0x32, 0xb2, 0x02, 0x0a, 0x07, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x73, 0x12, 0x47, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x1b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x0e, 0x53, 0x61, 0x76, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x61, 0x76, 0x65, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4a,
	0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0b, 0x5a, 0x09, 0x70, 0x6b,
	0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_secrets_proto_rawDescOnce sync.Once
	file_secrets_proto_rawDescData = file_secrets_proto_rawDesc
)

func file_secrets_proto_rawDescGZIP() []byte {
	file_secrets_proto_rawDescOnce.Do(func() {
		file_secrets_proto_rawDescData = protoimpl.X.CompressGZIP(file_secrets_proto_rawDescData)
	})
	return file_secrets_proto_rawDescData
}

var file_secrets_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_secrets_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_secrets_proto_goTypes = []any{
	(SecretType)(0),                 // 0: proto.SecretType
	(*Secret)(nil),                  // 1: proto.Secret
	(*GetUserSecretsResponse)(nil),  // 2: proto.GetUserSecretsResponse
	(*GetUserSecretRequest)(nil),    // 3: proto.GetUserSecretRequest
	(*GetUserSecretResponse)(nil),   // 4: proto.GetUserSecretResponse
	(*SaveUserSecretRequest)(nil),   // 5: proto.SaveUserSecretRequest
	(*DeleteUserSecretRequest)(nil), // 6: proto.DeleteUserSecretRequest
	(*timestamppb.Timestamp)(nil),   // 7: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),           // 8: google.protobuf.Empty
}
var file_secrets_proto_depIdxs = []int32{
	0,  // 0: proto.Secret.secret_type:type_name -> proto.SecretType
	7,  // 1: proto.Secret.created_at:type_name -> google.protobuf.Timestamp
	7,  // 2: proto.Secret.updated_at:type_name -> google.protobuf.Timestamp
	1,  // 3: proto.GetUserSecretsResponse.secrets:type_name -> proto.Secret
	1,  // 4: proto.GetUserSecretResponse.secret:type_name -> proto.Secret
	1,  // 5: proto.SaveUserSecretRequest.secret:type_name -> proto.Secret
	8,  // 6: proto.Secrets.GetUserSecrets:input_type -> google.protobuf.Empty
	3,  // 7: proto.Secrets.GetUserSecret:input_type -> proto.GetUserSecretRequest
	5,  // 8: proto.Secrets.SaveUserSecret:input_type -> proto.SaveUserSecretRequest
	6,  // 9: proto.Secrets.DeleteUserSecret:input_type -> proto.DeleteUserSecretRequest
	2,  // 10: proto.Secrets.GetUserSecrets:output_type -> proto.GetUserSecretsResponse
	4,  // 11: proto.Secrets.GetUserSecret:output_type -> proto.GetUserSecretResponse
	8,  // 12: proto.Secrets.SaveUserSecret:output_type -> google.protobuf.Empty
	8,  // 13: proto.Secrets.DeleteUserSecret:output_type -> google.protobuf.Empty
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_secrets_proto_init() }
func file_secrets_proto_init() {
	if File_secrets_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_secrets_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_secrets_proto_goTypes,
		DependencyIndexes: file_secrets_proto_depIdxs,
		EnumInfos:         file_secrets_proto_enumTypes,
		MessageInfos:      file_secrets_proto_msgTypes,
	}.Build()
	File_secrets_proto = out.File
	file_secrets_proto_rawDesc = nil
	file_secrets_proto_goTypes = nil
	file_secrets_proto_depIdxs = nil
}
