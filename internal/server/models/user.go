// Package models содержит модели данных приложения.
package models

import "time"

// User описывает структуру данных пользователя.
// Она включает в себя поля для хранения времени создания и обновления пользователя,
// а также уникальный идентификатор, логин и пароль.
type User struct {
	// ID это уникальный идентификатор пользователя.
	ID int `json:"id" db:"id"`
	// Login содержит логин пользователя, используемый для входа в систему.
	Login string `json:"login" db:"login"`
	// Password содержит пароль пользователя. Это поле не включается в JSON представление.
	Password string `json:"-" db:"password"`
	// CreatedAt содержит временную метку создания аккаунта пользователя.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// UpdatedAt содержит временную метку последнего обновления данных аккаунта пользователя.
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
