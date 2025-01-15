// Package grpc предоставляет функционал для создания и управления сервером gRPC,
// включая настройку TLS, регистрацию обработчиков и обработку сигналов завершения.
package grpc

import (
	"beliaev-aa/GophKeeper/certs"
	"beliaev-aa/GophKeeper/internal/server/config"
	"beliaev-aa/GophKeeper/internal/server/grpc/handlers"
	"beliaev-aa/GophKeeper/internal/server/grpc/interceptors"
	"beliaev-aa/GophKeeper/internal/server/service"
	"beliaev-aa/GophKeeper/internal/server/storage"
	"beliaev-aa/GophKeeper/pkg/proto"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server представляет сервер gRPC, содержащий конфигурацию, логгер и сам gRPC сервер.
type Server struct {
	config     *config.Config
	grpcServer *grpc.Server
	logger     *zap.Logger
}

// NewServer создает и инициализирует новый экземпляр сервера gRPC с заданными параметрами.
// Принимает конфигурацию сервера, хранилище данных и логгер.
func NewServer(config *config.Config, storage *storage.Storage, logger *zap.Logger) *Server {
	grpcServer := setupGRPCServer(config, storage, logger)
	return &Server{
		config:     config,
		grpcServer: grpcServer,
		logger:     logger,
	}
}

// setupGRPCServer настраивает и возвращает gRPC сервер с конфигурацией TLS и interceptors.
func setupGRPCServer(cfg *config.Config, storage *storage.Storage, logger *zap.Logger) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors.Authentication([]byte(cfg.SecretKey))),
		grpc.StreamInterceptor(interceptors.StreamAuthentication([]byte(cfg.SecretKey))),
	}

	tlsCredentials, err := loadTLSConfig("ca-cert.pem", "server-cert.pem", "server-key.pem")
	if err != nil {
		logger.Fatal("failed to load TLS credentials", zap.Error(err))
	}
	opts = append(opts, grpc.Creds(tlsCredentials))

	server := grpc.NewServer(opts...)

	proto.RegisterUsersServer(server, handlers.NewUserHandler(cfg, service.NewUserService(storage.UserRepository)))
	proto.RegisterSecretsServer(server, handlers.NewSecretHandler(logger, service.NewSecretService(storage.SecretRepository)))
	proto.RegisterNotificationServer(server, handlers.NewNotificationHandler(logger))

	return server
}

// Start запускает сервер gRPC и ожидает сигналы ОС для graceful завершения работы сервера.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Error(err))
	}

	go func() {
		if err = s.grpcServer.Serve(listener); err != nil {
			s.logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-quit
	s.logger.Info("Interrupt signal received", zap.String("signal", sig.String()))
	s.shutdown()

	return nil
}

// Останавливает сервер gRPC, осуществляя его graceful завершение.
func (s *Server) shutdown() {
	s.logger.Info("Shutting down server...")
	stopped := make(chan struct{})
	stopCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		s.logger.Info("Server shutdown successful")
	case <-stopCtx.Done():
		s.logger.Info("Shutdown timeout exceeded")
	}
}

// loadTLSConfig загружает TLS конфигурацию для сервера из указанных файлов сертификата и ключа.
func loadTLSConfig(caCertFile, serverCertFile, serverKeyFile string) (credentials.TransportCredentials, error) {
	caPem, err := certs.Cert.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	serverCertPEM, err := certs.Cert.ReadFile(serverCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read server cert: %w", err)
	}

	serverKeyPEM, err := certs.Cert.ReadFile(serverKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read server key: %w", err)
	}

	serverCert, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509 key pair: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, fmt.Errorf("failed to append CA cert to cert pool: %w", err)
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(tlsCfg), nil
}
