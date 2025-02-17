// Code generated by MockGen. DO NOT EDIT.
// Source: internal/server/storage/repository/secretRepository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "beliaev-aa/GophKeeper/pkg/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockISecretRepository is a mock of ISecretRepository interface.
type MockISecretRepository struct {
	ctrl     *gomock.Controller
	recorder *MockISecretRepositoryMockRecorder
}

// MockISecretRepositoryMockRecorder is the mock recorder for MockISecretRepository.
type MockISecretRepositoryMockRecorder struct {
	mock *MockISecretRepository
}

// NewMockISecretRepository creates a new mock instance.
func NewMockISecretRepository(ctrl *gomock.Controller) *MockISecretRepository {
	mock := &MockISecretRepository{ctrl: ctrl}
	mock.recorder = &MockISecretRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISecretRepository) EXPECT() *MockISecretRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockISecretRepository) Create(ctx context.Context, secret *models.Secret) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, secret)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockISecretRepositoryMockRecorder) Create(ctx, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockISecretRepository)(nil).Create), ctx, secret)
}

// Delete mocks base method.
func (m *MockISecretRepository) Delete(ctx context.Context, secretID, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, secretID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockISecretRepositoryMockRecorder) Delete(ctx, secretID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockISecretRepository)(nil).Delete), ctx, secretID, userID)
}

// GetSecret mocks base method.
func (m *MockISecretRepository) GetSecret(ctx context.Context, secretID, userID uint64) (*models.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", ctx, secretID, userID)
	ret0, _ := ret[0].(*models.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockISecretRepositoryMockRecorder) GetSecret(ctx, secretID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockISecretRepository)(nil).GetSecret), ctx, secretID, userID)
}

// GetUserSecrets mocks base method.
func (m *MockISecretRepository) GetUserSecrets(ctx context.Context, userID uint64) (models.Secrets, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSecrets", ctx, userID)
	ret0, _ := ret[0].(models.Secrets)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSecrets indicates an expected call of GetUserSecrets.
func (mr *MockISecretRepositoryMockRecorder) GetUserSecrets(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSecrets", reflect.TypeOf((*MockISecretRepository)(nil).GetUserSecrets), ctx, userID)
}

// Update mocks base method.
func (m *MockISecretRepository) Update(ctx context.Context, secret *models.Secret) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, secret)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockISecretRepositoryMockRecorder) Update(ctx, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockISecretRepository)(nil).Update), ctx, secret)
}
