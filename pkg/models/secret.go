// Package models содержит модели данных для приложения GophKeeper.
package models

import (
	"bytes"
	"fmt"
	"time"
)

type Secrets []*Secret

// Secret описывает структуру для хранения секретных данных пользователя.
type Secret struct {
	// ID - уникальный идентификатор секрета.
	ID uint64 `db:"id" json:"id"`
	// Metadata - метаданные, связанные с секретом.
	Metadata string `db:"metadata" json:"metadata"`
	// Payload - данные секрета в зашифрованном виде.
	Payload []byte `db:"payload" json:"payload"`
	// SecretType - тип секрета.
	SecretType string `db:"secret_type" json:"secret_type"`
	// Title - заголовок секрета.
	Title string `db:"title" json:"title"`
	// UserID - идентификатор пользователя, владельца секрета.
	UserID int `db:"user_id"`
	// CreatedAt - время создания секрета.
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	// UpdatedAt - время последнего обновления секрета.
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Следующие поля не включаются в базу данных, только во временные операции.
	// Creds - учетные данные, если SecretType = "credential".
	Creds *Credentials `db:"-"`
	// Text - текст, если SecretType = "text".
	Text *Text `db:"-"`
	// Blob - бинарные данные, если SecretType = "blob".
	Blob *Blob `db:"-"`
	// Card - данные карты, если SecretType = "card".
	Card *Card `db:"-"`
}

// NewSecret создаёт новый экземпляр Secret с указанным типом секрета.
func NewSecret(t SecretType) *Secret {
	return &Secret{SecretType: string(t)}
}

// SecretType определяет тип секрета, который может быть одним из предопределенных значений.
type SecretType string

const (
	// CredSecret - секрет, содержащий учетные данные.
	CredSecret SecretType = "credential"
	// TextSecret - секрет, содержащий текст.
	TextSecret SecretType = "text"
	// BlobSecret - секрет, содержащий бинарные данные.
	BlobSecret SecretType = "blob"
	// CardSecret - секрет, содержащий данные карты.
	CardSecret SecretType = "card"
	// UnknownSecret - неизвестный тип секрета.
	UnknownSecret SecretType = "unknown"
)

// Credentials описывает учетные данные пользователя.
type Credentials struct {
	// Login - логин пользователя.
	Login string `json:"login"`
	// Password - пароль пользователя.
	Password string `json:"password"`
}

// Text описывает текстовую информацию.
type Text struct {
	// Content - текстовое содержимое.
	Content string `json:"content"`
}

// Blob описывает бинарные данные файла.
type Blob struct {
	// FileName - имя файла.
	FileName string `json:"file_name"`
	// FileBytes - байты файла.
	FileBytes []byte `json:"file_bytes"`
}

// Card описывает данные банковской карты.
type Card struct {
	// Number - номер карты.
	Number string `json:"number"`
	// ExpYear - год истечения срока действия карты.
	ExpYear uint32 `json:"exp_year"`
	// ExpMonth - месяц истечения срока действия карты.
	ExpMonth uint32 `json:"exp_month"`
	// CVV - CVV-код карты.
	CVV uint32 `json:"cvv"`
}

// ToClipboard форматирует информацию секрета для использования в буфере обмена.
func (s Secret) ToClipboard() string {
	var b bytes.Buffer

	switch SecretType(s.SecretType) {
	case CredSecret:
		b.WriteString(fmt.Sprintf("login: %s\n", s.Creds.Login))
		b.WriteString(fmt.Sprintf("password: %s", s.Creds.Password))
	case CardSecret:
		b.WriteString(fmt.Sprintf("Card Number: %s\n", s.Card.Number))
		b.WriteString(fmt.Sprintf("Exp: %d/%d", s.Card.ExpMonth, s.Card.ExpYear))
		b.WriteString(fmt.Sprintf("CVV: %d", s.Card.CVV))
	case TextSecret:
		b.WriteString(fmt.Sprintf("Text: %s\n", s.Text.Content))
	case BlobSecret:
		b.WriteString("File data cannot be moved to clipboard\n")
	}

	return b.String()
}
