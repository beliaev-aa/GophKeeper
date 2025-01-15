package grpc

import (
	"beliaev-aa/GophKeeper/internal/server/config"
	"beliaev-aa/GophKeeper/internal/server/storage"
	"beliaev-aa/GophKeeper/tests/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
	"testing"
)

func TestNewServer(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		Address:   ":50051",
		SecretKey: "test-secret",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := &storage.Storage{
		UserRepository:   mocks.NewMockIUserRepository(ctrl),
		SecretRepository: mocks.NewMockISecretRepository(ctrl),
	}

	srv := NewServer(cfg, mockStorage, logger)

	assert.NotNil(t, srv)
	assert.Equal(t, cfg, srv.config)
	assert.NotNil(t, srv.grpcServer)
	assert.Equal(t, logger, srv.logger)
}

func TestSetupGRPCServer(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		SecretKey: "test-secret",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := &storage.Storage{
		UserRepository:   mocks.NewMockIUserRepository(ctrl),
		SecretRepository: mocks.NewMockISecretRepository(ctrl),
	}

	server := setupGRPCServer(cfg, mockStorage, logger)

	assert.NotNil(t, server)
}

func TestLoadTLSConfig(t *testing.T) {
	t.Run("Valid certificates", func(t *testing.T) {
		tlsConfig, err := loadTLSConfig("ca-cert.pem", "server-cert.pem", "server-key.pem")
		if err != nil {
			t.Skip("TLS setup requires valid certificates")
		}
		assert.NotNil(t, tlsConfig)
		assert.Implements(t, (*credentials.TransportCredentials)(nil), tlsConfig)
	})

	t.Run("Invalid certificates", func(t *testing.T) {
		_, err := loadTLSConfig("invalid-ca.pem", "invalid-cert.pem", "invalid-key.pem")
		assert.Error(t, err)
	})
}
