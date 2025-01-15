// Package config содержит конфигурации для сервера.
package config

import (
	"errors"
	"github.com/spf13/viper"
	"strings"
)

// Config представляет основную конфигурацию клиентского приложения.
type Config struct {
	Address     string // Address определяет адрес сервера.
	PostgresDSN string // PostgresDSN содержит строку подключения к PostgreSQL.
	SecretKey   string // SecretKey используется для подписи JWT.
}

// LoadConfig инициализирует и возвращает новый экземпляр конфигурации.
// Ошибка возвращается, если обязательные конфигурационные параметры не заданы.
func LoadConfig() (*Config, error) {
	viper.SetEnvPrefix("GOPHKEEPER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	address := viper.GetString("address")
	if address == "" {
		return nil, errors.New("server address is not set: set GOPHKEEPER_ADDRESS environment variable")
	}

	postgresDSN := viper.GetString("postgres-dsn")
	if postgresDSN == "" {
		return nil, errors.New("PostgreSQL DSN is not set: set GOPHKEEPER_POSTGRES_DSN environment variable")
	}

	secretKey := viper.GetString("secret-key")
	if secretKey == "" {
		return nil, errors.New("secret key for signing JWT is not set: set GOPHKEEPER_SECRET_KEY environment variable")
	}

	return &Config{
		Address:     address,
		PostgresDSN: postgresDSN,
		SecretKey:   secretKey,
	}, nil
}
