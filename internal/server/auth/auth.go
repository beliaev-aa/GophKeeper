// Package auth предоставляет функции для аутентификации и авторизации, включая работу с паролями и JWT токенами.
package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// HashPassword хэширует пароль с использованием алгоритма bcrypt.
// Возвращает хэшированный пароль или ошибку.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword проверяет совпадение предоставленного пароля с его хэшем.
// Возвращает true, если пароли совпадают.
func CheckPassword(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}

// CreateToken создает JWT токен для пользователя.
// Принимает идентификатор пользователя, время истечения токена и секретный ключ.
// Возвращает строку с токеном или ошибку.
func CreateToken(userID int, expireDate time.Time, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"iss":     "gophkeeper",
		"exp":     expireDate.Unix(),
		"iat":     time.Now().Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyToken проверяет JWT токен и извлекает его утверждения (claims).
// Возвращает утверждения или ошибку, если токен невалиден или истек срок его действия.
func VerifyToken(tokenText string, secretKey []byte) (jwt.MapClaims, error) {
	token, err := ParseToken(tokenText, secretKey)
	if err != nil {
		return nil, err
	}
	claims, err := GetClaims(token)
	if err != nil {
		return nil, err
	}
	if IsExpired(claims) {
		return nil, fmt.Errorf("token is expired")
	}
	return claims, nil
}

// ParseToken разбирает строку JWT токена и возвращает объект токена.
// Возвращает токен или ошибку, если токен не удалось распознать.
func ParseToken(tokenText string, secretKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenText, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

// IsExpired проверяет, истек ли срок действия JWT токена.
// Возвращает true, если срок действия токена истек.
func IsExpired(claims jwt.MapClaims) bool {
	return float64(time.Now().Unix()) > claims["exp"].(float64)
}

// GetClaims извлекает claims из JWT токена.
// Возвращает утверждения или ошибку, если не удалось извлечь их из токена.
func GetClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from token")
	}
	return claims, nil
}
