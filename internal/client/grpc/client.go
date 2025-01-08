// Package grpc предоставляет клиентские утилиты для работы с gRPC сервером.
// Включает методы для аутентификации, работы с секретами и уведомлениями.
package grpc

import (
	"beliaev-aa/GophKeeper/certs"
	"beliaev-aa/GophKeeper/internal/client/config"
	"beliaev-aa/GophKeeper/internal/client/grpc/interceptors"
	"beliaev-aa/GophKeeper/pkg/converter"
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/pkg/proto"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"math/rand/v2"
	"sync"
	"time"
)

// ClientGRPC управляет соединением с gRPC сервером и реализует методы для работы с серверными ресурсами.
type ClientGRPC struct {
	config        *config.Config
	usersClient   proto.UsersClient
	secretsClient proto.SecretsClient
	notifyClient  proto.NotificationClient
	accessToken   string
	password      string
	clientID      uint64
	previews      sync.Map
}

// NewClientGRPC создаёт новый экземпляр ClientGRPC с предварительной настройкой подключения к серверу.
func NewClientGRPC(cfg *config.Config) (*ClientGRPC, error) {
	var opts []grpc.DialOption

	newClient := ClientGRPC{
		config:   cfg,
		clientID: uint64(rand.IntN(math.MaxInt32)),
	}

	opts = append(
		opts,
		grpc.WithChainUnaryInterceptor(
			interceptors.Timeout(time.Second*5),
			interceptors.AddAuth(&newClient.accessToken, uint32(newClient.clientID)),
		),
	)

	opts = append(
		opts,
		grpc.WithStreamInterceptor(interceptors.AddAuthStream(&newClient.accessToken, newClient.clientID)),
	)

	tlsCredential, err := loadTLSConfig("ca-cert.pem", "client-cert.pem", "client-key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS config: %w", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(tlsCredential))

	c, err := grpc.NewClient(cfg.ServerAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	newClient.usersClient = proto.NewUsersClient(c)
	newClient.secretsClient = proto.NewSecretsClient(c)
	newClient.notifyClient = proto.NewNotificationClient(c)

	return &newClient, nil
}

// Login авторизует пользователя на сервере и получает токен доступа.
func (c *ClientGRPC) Login(ctx context.Context, login string, password string) (string, error) {
	req := &proto.LoginRequest{
		Login:    login,
		Password: password,
	}

	response, err := c.usersClient.Login(ctx, req)
	if err != nil {
		return "", parseError(err)
	}

	c.accessToken = response.AccessToken

	return response.AccessToken, nil
}

// Register регистрирует нового пользователя и получает токен доступа.
func (c *ClientGRPC) Register(ctx context.Context, login string, password string) (string, error) {
	req := &proto.RegisterRequest{
		Login:    login,
		Password: password,
	}

	response, err := c.usersClient.Register(ctx, req)
	if err != nil {
		return "", parseError(err)
	}

	c.accessToken = response.AccessToken

	return response.AccessToken, nil
}

// LoadSecrets загружает список секретов пользователя.
func (c *ClientGRPC) LoadSecrets(ctx context.Context) ([]*models.Secret, error) {
	request := emptypb.Empty{}

	response, err := c.secretsClient.GetUserSecrets(ctx, &request)
	if err != nil {
		return nil, parseError(err)
	}

	secrets := converter.ProtoToSecrets(response.Secrets)
	return secrets, nil
}

// LoadSecret загружает информацию о конкретном секрете.
func (c *ClientGRPC) LoadSecret(_ context.Context, ID uint64) (*models.Secret, error) {
	request := &proto.GetUserSecretRequest{
		Id: ID,
	}

	response, err := c.secretsClient.GetUserSecret(context.Background(), request)
	if err != nil {
		return nil, parseError(err)
	}

	secret := converter.ProtoToSecret(response.Secret)

	return secret, nil
}

// SaveSecret сохраняет или обновляет секрет пользователя на сервере.
func (c *ClientGRPC) SaveSecret(ctx context.Context, secret *models.Secret) error {
	sec := &proto.Secret{
		Title:      secret.Title,
		Metadata:   secret.Metadata,
		SecretType: converter.TypeToProto(secret.SecretType),
		Payload:    secret.Payload,
		CreatedAt:  timestamppb.New(secret.CreatedAt),
		UpdatedAt:  timestamppb.New(secret.UpdatedAt),
	}

	if secret.ID > 0 {
		sec.Id = secret.ID
	}

	request := &proto.SaveUserSecretRequest{Secret: sec}
	_, err := c.secretsClient.SaveUserSecret(ctx, request)

	return parseError(err)
}

// DeleteSecret удаляет секрет пользователя.
func (c *ClientGRPC) DeleteSecret(ctx context.Context, id uint64) error {
	request := &proto.DeleteUserSecretRequest{Id: id}
	_, err := c.secretsClient.DeleteUserSecret(ctx, request)

	return parseError(err)
}

// SetToken устанавливает текущий токен доступа клиента.
func (c *ClientGRPC) SetToken(token string) {
	c.accessToken = token
}

// GetToken возвращает текущий токен доступа клиента.
func (c *ClientGRPC) GetToken() string {
	return c.accessToken
}

// SetPassword устанавливает текущий пароль клиента.
func (c *ClientGRPC) SetPassword(password string) {
	c.password = password
}

// GetPassword возвращает текущий пароль клиента.
func (c *ClientGRPC) GetPassword() string {
	return c.password
}

// parseError анализирует ошибки от gRPC вызовов и конвертирует их в более понятный формат.
func parseError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.Unavailable:
		return errors.New("server unavailable")
	case codes.Unauthenticated:
		return errors.New("failed to authenticate")
	case codes.AlreadyExists:
		return errors.New("user already exists")
	default:
		return err
	}
}

// loadTLSConfig загружает TLS конфигурацию для подключения к серверу.
func loadTLSConfig(caCertFile, clientCertFile, clientKeyFile string) (credentials.TransportCredentials, error) {
	caPem, err := certs.Cert.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	clientCertPEM, err := certs.Cert.ReadFile(clientCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client cert: %w", err)
	}

	clientKeyPEM, err := certs.Cert.ReadFile(clientKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client key: %w", err)
	}

	clientCert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509 key pair: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, fmt.Errorf("failed to append CA cert to cert pool: %w", err)
	}

	tlcConfiguration := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(tlcConfiguration), nil
}

// Notifications подписывается на уведомления сервера и обновляет UI при получении новых данных.
func (c *ClientGRPC) Notifications(p *tea.Program) {
	var (
		stream proto.Notification_SubscribeClient
		err    error
	)

	for {
		if stream == nil {
			if stream, err = c.subscribe(); err != nil {
				time.Sleep(time.Second * 2)
				continue
			}
		}

		_, err = stream.Recv()
		if err != nil {
			stream = nil
			time.Sleep(time.Second * 2)
			continue
		}

		if p != nil {
			p.Send(ReloadSecretList{})
		}
	}
}

type ReloadSecretList struct{}

// Инициирует подписку на серверные уведомления, используя ID клиента.
func (c *ClientGRPC) subscribe() (proto.Notification_SubscribeClient, error) {
	return c.notifyClient.Subscribe(context.Background(), &proto.SubscribeRequest{
		Id: c.clientID,
	})
}
