package auth

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestPasswords(t *testing.T) {
	password := "test"
	longPassword := strings.Repeat("A", 100)

	t.Run("hash_password", func(t *testing.T) {
		pwHash, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, pwHash, "Hash should not be empty")
	})

	t.Run("hash_long_password", func(t *testing.T) {
		_, err := HashPassword(longPassword)
		assert.Error(t, err, "Should return an error for long password")
	})

	t.Run("check_correct_password", func(t *testing.T) {
		pwHash, _ := HashPassword(password)
		checkResult := CheckPassword(pwHash, password)
		assert.True(t, checkResult, "Check password should return true for correct password")
	})

	t.Run("check_incorrect_password", func(t *testing.T) {
		pwHash, _ := HashPassword(password)
		checkResult := CheckPassword(pwHash, "wrongpassword")
		assert.False(t, checkResult, "Check password should return false for incorrect password")
	})
}

func TestTokens(t *testing.T) {
	secret := []byte("test_secret_key")
	testUserID := 1337
	expireDate := time.Now().Add(time.Hour)

	t.Run("create_token", func(t *testing.T) {
		token, err := CreateToken(testUserID, expireDate, secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("successful_token_verification", func(t *testing.T) {
		tokenString, _ := CreateToken(testUserID, expireDate, secret)
		claims, err := VerifyToken(tokenString, secret)
		assert.NoError(t, err)
		assert.Equal(t, float64(testUserID), claims["user_id"].(float64), "Claims user_id should match")
	})

	t.Run("verify_invalid_token", func(t *testing.T) {
		_, err := VerifyToken("invalidtoken", secret)
		assert.Error(t, err, "Verification should fail for invalid token")
	})

	t.Run("verify_token_with_incorrect_alg", func(t *testing.T) {
		tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
		_, err := ParseToken(tokenString, secret)
		assert.Error(t, err, "Parse should fail with incorrect algorithm")
	})
}
