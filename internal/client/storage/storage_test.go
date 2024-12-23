package storage

import (
	gophKeeperErrors "beliaev-aa/GophKeeper/pkg/errors"
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/pkg/storage"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"time"
)

func TestNewSQLiteStorage_in_memory(t *testing.T) {
	type testCase struct {
		name    string
		dbPath  string
		wantErr bool
	}

	testCases := []testCase{
		{
			name:    "valid_in_memory_db",
			dbPath:  ":memory:",
			wantErr: false,
		},
		{
			name:    "invalid_file_path",
			dbPath:  "/non/existing/path",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			st, err := NewSQLiteStorage(tc.dbPath)
			if (err != nil) != tc.wantErr {
				t.Fatalf("NewSQLiteStorage() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if !tc.wantErr && st == nil {
				t.Error("expected non-nil storage when no error")
			}
		})
	}
}

func setupTestDB(t *testing.T) storage.Storage {
	t.Helper()

	st, err := NewSQLiteStorage(":memory:")
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	return st
}

func TestSQLiteStorage_Create(t *testing.T) {
	type testCase struct {
		name      string
		input     *models.StoredData
		wantErr   bool
		duplicate bool
	}

	st := setupTestDB(t)

	testCases := []testCase{
		{
			name: "valid_data_login_password",
			input: &models.StoredData{
				ID:          "id_1",
				Type:        models.LoginPasswordType,
				Description: "desc A",
				MetaInfo: []models.MetaData{
					{Key: "website", Value: "example.com"},
				},
				Content: models.LoginPassword{
					Login:    "user123",
					Password: "p@ssw0rd",
				},
				UpdatedAt: time.Now().Unix(),
			},
			wantErr:   false,
			duplicate: false,
		},
		{
			name: "duplicate_id",
			input: &models.StoredData{
				ID:          "id_1",
				Type:        models.TextDataType,
				Description: "desc B",
				MetaInfo: []models.MetaData{
					{Key: "note", Value: "duplicate check"},
				},
				Content: models.TextData{
					Text: "some text content",
				},
				UpdatedAt: time.Now().Unix(),
			},
			wantErr:   true,
			duplicate: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := st.Create(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Create() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestSQLiteStorage_Read(t *testing.T) {
	type testCase struct {
		name         string
		readID       string
		wantErr      bool
		wantNotFound bool
	}

	st := setupTestDB(t)

	knownData := &models.StoredData{
		ID:          "id_known",
		Type:        models.LoginPasswordType,
		Description: "known_desc",
		MetaInfo: []models.MetaData{
			{Key: "env", Value: "prod"},
		},
		Content: models.LoginPassword{
			Login:    "admin",
			Password: "secret123",
		},
		UpdatedAt: time.Now().Unix(),
	}
	if err := st.Create(knownData); err != nil {
		t.Fatalf("failed to create knownData: %v", err)
	}

	testCases := []testCase{
		{
			name:         "read_existing_id",
			readID:       "id_known",
			wantErr:      false,
			wantNotFound: false,
		},
		{
			name:         "read_unknown_id",
			readID:       "id_not_in_db",
			wantErr:      true,
			wantNotFound: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := st.Read(tc.readID)
			if tc.wantNotFound {
				if !errors.Is(err, gophKeeperErrors.ErrNotFound) {
					t.Fatalf("Read() error = %v, want ErrNotFound", err)
				}
				return
			}

			if (err != nil) != tc.wantErr {
				t.Fatalf("Read() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if got != nil && got.ID != tc.readID {
				t.Errorf("Read() got ID = %s, want ID = %s", got.ID, tc.readID)
			}
		})
	}
}

func TestSQLiteStorage_Update(t *testing.T) {
	type testCase struct {
		name         string
		updateID     string
		updateData   *models.StoredData
		wantErr      bool
		wantNotFound bool
	}

	st := setupTestDB(t)

	existingData := &models.StoredData{
		ID:          "id_for_update",
		Type:        models.LoginPasswordType,
		Description: "old_desc",
		MetaInfo: []models.MetaData{
			{Key: "old", Value: "data"},
		},
		Content: models.LoginPassword{
			Login:    "old_login",
			Password: "old_pass",
		},
		UpdatedAt: time.Now().Unix(),
	}
	if err := st.Create(existingData); err != nil {
		t.Fatalf("failed to create existingData: %v", err)
	}

	testCases := []testCase{
		{
			name:     "update_existing_record",
			updateID: "id_for_update",
			updateData: &models.StoredData{
				Type:        models.LoginPasswordType,
				Description: "new_desc",
				MetaInfo: []models.MetaData{
					{Key: "updated", Value: "true"},
				},
				Content: models.LoginPassword{
					Login:    "new_login",
					Password: "new_pass",
				},
				UpdatedAt: time.Now().Unix(),
			},
			wantErr:      false,
			wantNotFound: false,
		},
		{
			name:     "update_non_existing_record",
			updateID: "id_not_in_db",
			updateData: &models.StoredData{
				Type:        models.TextDataType,
				Description: "will_fail",
				MetaInfo:    []models.MetaData{},
				Content: models.TextData{
					Text: "dummy content",
				},
				UpdatedAt: time.Now().Unix(),
			},
			wantErr:      false,
			wantNotFound: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := st.Update(tc.updateID, tc.updateData)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Update() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if !tc.wantErr && !tc.wantNotFound {
				got, err := st.Read(tc.updateID)
				if err != nil {
					t.Fatalf("Read() after Update error = %v", err)
				}
				if got != nil && got.Description != tc.updateData.Description {
					t.Errorf("Update() got Description = %s, want = %s", got.Description, tc.updateData.Description)
				}
			}
		})
	}
}

func TestSQLiteStorage_Delete(t *testing.T) {
	type testCase struct {
		name         string
		deleteID     string
		wantErr      bool
		wantNotFound bool
	}

	st := setupTestDB(t)

	dataToDelete := &models.StoredData{
		ID:          "id_for_delete",
		Type:        models.CardDataType,
		Description: "desc del",
		MetaInfo: []models.MetaData{
			{Key: "del", Value: "val"},
		},
		Content: models.CardData{
			CardNumber:     "1111 2222 3333 4444",
			ExpiryDate:     "10/28",
			CVV:            "123",
			CardHolderName: "John Doe",
		},
		UpdatedAt: time.Now().Unix(),
	}
	if err := st.Create(dataToDelete); err != nil {
		t.Fatalf("failed to create dataToDelete: %v", err)
	}

	testCases := []testCase{
		{
			name:         "delete_existing_record",
			deleteID:     "id_for_delete",
			wantErr:      false,
			wantNotFound: false,
		},
		{
			name:         "delete_non_existing_record",
			deleteID:     "id_not_in_db",
			wantErr:      false,
			wantNotFound: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := st.Delete(tc.deleteID)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Delete() error = %v, wantErr = %v", err, tc.wantErr)
			}

			if !tc.wantErr && !tc.wantNotFound {
				got, err := st.Read(tc.deleteID)
				if err == nil && got != nil {
					t.Errorf("Delete() record still found, expected it to be deleted")
				}
			}
		})
	}
}

func TestSQLiteStorage_List(t *testing.T) {
	type testCase struct {
		name        string
		initialData []*models.StoredData
		wantCount   int
	}

	st := setupTestDB(t)

	testCases := []testCase{
		{
			name:        "empty_storage",
			initialData: []*models.StoredData{},
			wantCount:   0,
		},
		{
			name: "some_records",
			initialData: []*models.StoredData{
				{
					ID:          "id_list_1",
					Type:        models.BinaryDataType,
					Description: "d1",
					MetaInfo:    []models.MetaData{},
					Content: models.BinaryData{
						Data: []byte{0x01, 0x02, 0x03},
					},
					UpdatedAt: time.Now().Unix(),
				},
				{
					ID:          "id_list_2",
					Type:        models.TextDataType,
					Description: "d2",
					MetaInfo: []models.MetaData{
						{Key: "tag", Value: "sample"},
					},
					Content: models.TextData{
						Text: "some text data",
					},
					UpdatedAt: time.Now().Unix(),
				},
			},
			wantCount: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, d := range tc.initialData {
				if err := st.Create(d); err != nil {
					t.Fatalf("failed to create data: %v", err)
				}
			}

			gotList, err := st.List()
			if err != nil {
				t.Fatalf("List() error = %v", err)
			}
			if len(gotList) != tc.wantCount {
				t.Errorf("List() got %d records, want %d", len(gotList), tc.wantCount)
			}
		})
	}
}
