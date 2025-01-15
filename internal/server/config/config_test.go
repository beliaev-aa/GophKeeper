package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		setupEnv       func()
		expectedError  string
		expectedConfig *Config
	}{
		{
			name: "No_ENV_Variables",
			setupEnv: func() {
				os.Unsetenv("GOPHKEEPER_ADDRESS")
				os.Unsetenv("GOPHKEEPER_POSTGRES_DSN")
				os.Unsetenv("GOPHKEEPER_SECRET_KEY")
			},
			expectedError: "server address is not set: set GOPHKEEPER_ADDRESS environment variable",
		},
		{
			name: "No_Address",
			setupEnv: func() {
				os.Setenv("GOPHKEEPER_ADDRESS", "")
				os.Setenv("GOPHKEEPER_POSTGRES_DSN", "some-dsn")
				os.Setenv("GOPHKEEPER_SECRET_KEY", "some-secret")
			},
			expectedError: "server address is not set: set GOPHKEEPER_ADDRESS environment variable",
		},
		{
			name: "No_PostgresDSN",
			setupEnv: func() {
				os.Setenv("GOPHKEEPER_ADDRESS", "127.0.0.1:5000")
				os.Setenv("GOPHKEEPER_POSTGRES_DSN", "")
				os.Setenv("GOPHKEEPER_SECRET_KEY", "some-secret")
			},
			expectedError: "PostgreSQL DSN is not set: set GOPHKEEPER_POSTGRES_DSN environment variable",
		},
		{
			name: "No_SecretKey",
			setupEnv: func() {
				os.Setenv("GOPHKEEPER_ADDRESS", "127.0.0.1:5000")
				os.Setenv("GOPHKEEPER_POSTGRES_DSN", "some-dsn")
				os.Setenv("GOPHKEEPER_SECRET_KEY", "")
			},
			expectedError: "secret key for signing JWT is not set: set GOPHKEEPER_SECRET_KEY environment variable",
		},
		{
			name: "All_Variables_Set_Correctly",
			setupEnv: func() {
				os.Setenv("GOPHKEEPER_ADDRESS", "127.0.0.1:5000")
				os.Setenv("GOPHKEEPER_POSTGRES_DSN", "some-dsn")
				os.Setenv("GOPHKEEPER_SECRET_KEY", "some-secret")
			},
			expectedConfig: &Config{
				Address:     "127.0.0.1:5000",
				PostgresDSN: "some-dsn",
				SecretKey:   "some-secret",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupEnv()
			viper.Reset()

			config, err := LoadConfig()

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Nil(t, config)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, config)
			}
		})
	}
}
