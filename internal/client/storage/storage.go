// Package storage предоставляет реализацию хранилища данных с использованием SQLite.
// Хранилище позволяет выполнять CRUD-операции над данными, описанными в пакете models.
package storage

import (
	gophKeeperErrors "beliaev-aa/GophKeeper/pkg/errors"
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/pkg/storage"
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStorage предоставляет реализацию интерфейса хранилища данных с использованием SQLite.
type SQLiteStorage struct {
	db *sql.DB // Соединение с базой данных.
}

// NewSQLiteStorage создаёт новый экземпляр SQLiteStorage.
//
// Параметры:
//   - dbPath: путь к файлу базы данных SQLite.
//
// Возвращает:
//   - указатель на созданный SQLiteStorage;
//   - ошибку, если инициализация не удалась.
func NewSQLiteStorage(dbPath string) (storage.Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	store := &SQLiteStorage{db: db}
	if err = store.initSchema(); err != nil {
		return nil, err
	}

	return store, nil
}

// initSchema инициализирует схему таблицы для хранения данных в базе данных.
//
// Возвращает:
//   - ошибку, если создание таблицы не удалось.
func (s *SQLiteStorage) initSchema() error {
	query := `
    CREATE TABLE IF NOT EXISTS stored_data (
        id TEXT PRIMARY KEY,
        type TEXT NOT NULL,
        description TEXT,
        meta_info TEXT,
        content TEXT,
        updated_at INTEGER
    );
    `
	_, err := s.db.Exec(query)
	return err
}

// Create добавляет новую запись в хранилище.
//
// Параметры:
//   - data: указатель на структуру models. StoredData, содержащую данные для сохранения.
//
// Возвращает:
//   - ошибку, если операция не удалась.
func (s *SQLiteStorage) Create(data *models.StoredData) error {
	metaInfo, err := json.Marshal(data.MetaInfo)
	if err != nil {
		return err
	}

	content, err := json.Marshal(data.Content)
	if err != nil {
		return err
	}

	query := `
    INSERT INTO stored_data (id, type, description, meta_info, content, updated_at)
    VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err = s.db.Exec(query, data.ID, data.Type, data.Description, string(metaInfo), string(content), data.UpdatedAt)
	return err
}

// Read возвращает запись из хранилища по её уникальному идентификатору.
//
// Параметры:
//   - id: уникальный идентификатор записи.
//
// Возвращает:
//   - указатель на models. StoredData с данными записи;
//   - nil, если запись не найдена;
//   - ошибку, если операция не удалась.
func (s *SQLiteStorage) Read(id string) (*models.StoredData, error) {
	query := `SELECT id, type, description, meta_info, content, updated_at FROM stored_data WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var (
		storedData models.StoredData
		metaInfo   string
		content    string
	)

	if err := row.Scan(&storedData.ID, &storedData.Type, &storedData.Description, &metaInfo, &content, &storedData.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gophKeeperErrors.ErrNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(metaInfo), &storedData.MetaInfo); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(content), &storedData.Content); err != nil {
		return nil, err
	}

	return &storedData, nil
}

// Update обновляет существующую запись в хранилище.
//
// Параметры:
//   - id: уникальный идентификатор записи;
//   - data: указатель на структуру models. StoredData с обновлёнными данными.
//
// Возвращает:
//   - ошибку, если операция не удалась.
func (s *SQLiteStorage) Update(id string, data *models.StoredData) error {
	metaInfo, err := json.Marshal(data.MetaInfo)
	if err != nil {
		return err
	}

	content, err := json.Marshal(data.Content)
	if err != nil {
		return err
	}

	query := `
    UPDATE stored_data
    SET type = ?, description = ?, meta_info = ?, content = ?, updated_at = ?
    WHERE id = ?
    `
	_, err = s.db.Exec(query, data.Type, data.Description, string(metaInfo), string(content), data.UpdatedAt, id)
	return err
}

// Delete удаляет запись из хранилища по её уникальному идентификатору.
//
// Параметры:
//   - id: уникальный идентификатор записи.
//
// Возвращает:
//   - ошибку, если операция не удалась.
func (s *SQLiteStorage) Delete(id string) error {
	query := `DELETE FROM stored_data WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

// List возвращает список всех записей из хранилища.
//
// Возвращает:
//   - массив указателей на структуры models. StoredData с данными записей;
//   - ошибку, если операция не удалась.
func (s *SQLiteStorage) List() ([]*models.StoredData, error) {
	query := `SELECT id, type, description, meta_info, content, updated_at FROM stored_data`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []*models.StoredData
	for rows.Next() {
		var (
			storedData models.StoredData
			metaInfo   string
			content    string
		)

		if err = rows.Scan(&storedData.ID, &storedData.Type, &storedData.Description, &metaInfo, &content, &storedData.UpdatedAt); err != nil {
			return nil, err
		}

		if err = json.Unmarshal([]byte(metaInfo), &storedData.MetaInfo); err != nil {
			return nil, err
		}

		if err = json.Unmarshal([]byte(content), &storedData.Content); err != nil {
			return nil, err
		}

		dataList = append(dataList, &storedData)
	}

	return dataList, nil
}
