// Package models содержит структуры данных для хранения и управления различными типами информации,
// включая пары логин/пароль, произвольные текстовые и бинарные данные, одноразовые пароли (OTP) и данные банковских карт.
package models

import (
	"beliaev-aa/GophKeeper/internal/client/crypto"
	"encoding/json"
	"errors"
	"fmt"
)

// DataType определяет типы хранимых данных.
type DataType string

const (
	// LoginPasswordType представляет данные типа пара логин/пароль.
	LoginPasswordType DataType = "login_password"
	// TextDataType представляет произвольные текстовые данные.
	TextDataType DataType = "text_data"
	// BinaryDataType представляет произвольные бинарные данные.
	BinaryDataType DataType = "binary_data"
	// CardDataType представляет данные банковских карт.
	CardDataType DataType = "card_data"
	// OTPType представляет одноразовые пароли (One-Time Passwords).
	OTPType DataType = "otp"
)

// MetaData представляет общую структуру для хранения метаинформации о данных,
// такой как принадлежность данных к веб-сайту, личности, банку и другие дополнительные сведения.
type MetaData struct {
	Key   string `json:"key"`   // Ключ метаинформации, например, "website" или "owner".
	Value string `json:"value"` // Значение метаинформации, например, "example.com" или "John Doe".
}

// StoredData представляет базовую структуру для хранения данных любого типа с поддержкой произвольной метаинформации.
type StoredData struct {
	ID          string      `json:"id"`          // Уникальный идентификатор записи.
	Type        DataType    `json:"type"`        // Тип хранимой информации, определяемый значением DataType.
	Description string      `json:"description"` // Описание или комментарий для данных.
	MetaInfo    []MetaData  `json:"meta_info"`   // Список произвольной метаинформации о данных.
	Content     interface{} `json:"content"`     // Содержимое данных, тип зависит от значения поля Type.
	UpdatedAt   int64       `json:"updated_at"`  // Временная метка последнего обновления данных (в формате Unix Timestamp).
}

// LoginPassword представляет данные типа пара логин/пароль.
type LoginPassword struct {
	Login    string `json:"login"`    // Логин пользователя.
	Password string `json:"password"` // Пароль пользователя.
}

// TextData представляет произвольные текстовые данные.
type TextData struct {
	Text string `json:"text"` // Хранимый текст.
}

// BinaryData представляет произвольные бинарные данные.
type BinaryData struct {
	Data []byte `json:"data"` // Хранимые бинарные данные в виде массива байтов.
}

// OneTimePassword представляет структуру для хранения одноразовых паролей (OTP).
type OneTimePassword struct {
	OTPCode string `json:"otp_code"` // Одноразовый пароль.
	Expiry  string `json:"expiry"`   // Дата истечения срока действия пароля в строковом формате.
}

// CardData представляет данные банковских карт.
type CardData struct {
	CardNumber     string `json:"card_number"`      // Номер банковской карты.
	ExpiryDate     string `json:"expiry_date"`      // Дата истечения срока действия карты (например, "12/24").
	CVV            string `json:"cvv"`              // Код безопасности карты (CVV).
	CardHolderName string `json:"card_holder_name"` // Имя владельца карты.
}

// EncryptContent шифрует поле Content в зависимости от типа данных.
func (s *StoredData) EncryptContent(key []byte) error {
	if s.Content == nil {
		return errors.New("content is nil")
	}

	var plaintext []byte
	var err error

	// Преобразуем Content в JSON для шифрования
	switch s.Type {
	case LoginPasswordType, TextDataType, BinaryDataType, CardDataType, OTPType:
		plaintext, err = json.Marshal(s.Content)
		if err != nil {
			return fmt.Errorf("failed to marshal content: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type: %s", s.Type)
	}

	// Шифруем данные
	encrypted, err := crypto.Encrypt(string(plaintext), key)
	if err != nil {
		return fmt.Errorf("failed to encrypt content: %w", err)
	}

	// Сохраняем зашифрованную строку в Content
	s.Content = encrypted
	return nil
}

// DecryptContent расшифровывает поле Content в зависимости от типа данных.
func (s *StoredData) DecryptContent(key []byte) error {
	encrypted, ok := s.Content.(string)
	if !ok {
		return errors.New("content is not a string")
	}

	// Расшифровываем данные
	decrypted, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		return fmt.Errorf("failed to decrypt content: %w", err)
	}

	// Преобразуем расшифрованные данные обратно в соответствующую структуру
	switch s.Type {
	case LoginPasswordType:
		var content LoginPassword
		if err = json.Unmarshal([]byte(decrypted), &content); err != nil {
			return fmt.Errorf("failed to unmarshal content to LoginPassword: %w", err)
		}
		s.Content = content
	case TextDataType:
		var content TextData
		if err = json.Unmarshal([]byte(decrypted), &content); err != nil {
			return fmt.Errorf("failed to unmarshal content to TextData: %w", err)
		}
		s.Content = content
	case BinaryDataType:
		var content BinaryData
		if err = json.Unmarshal([]byte(decrypted), &content); err != nil {
			return fmt.Errorf("failed to unmarshal content to BinaryData: %w", err)
		}
		s.Content = content
	case CardDataType:
		var content CardData
		if err = json.Unmarshal([]byte(decrypted), &content); err != nil {
			return fmt.Errorf("failed to unmarshal content to CardData: %w", err)
		}
		s.Content = content
	case OTPType:
		var content OneTimePassword
		if err = json.Unmarshal([]byte(decrypted), &content); err != nil {
			return fmt.Errorf("failed to unmarshal content to OneTimePassword: %w", err)
		}
		s.Content = content
	default:
		return fmt.Errorf("unsupported data type: %s", s.Type)
	}

	return nil
}
