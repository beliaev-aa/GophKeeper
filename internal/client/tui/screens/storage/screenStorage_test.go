package storage

import (
	"beliaev-aa/GophKeeper/internal/client/grpc"
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/tests/mocks"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_BrowseStorageScreen_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()

	screen := NewStorageBrowseScreenScreen(mockStorage)
	cmd := screen.Init()

	if cmd != nil {
		t.Errorf("Init should not return a command")
	}

	if len(screen.table.Rows()) != 0 {
		t.Errorf("Expected empty table rows, got %d", len(screen.table.Rows()))
	}
}

func Test_BrowseStorageScreen_View(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
	mockStorage.EXPECT().String().AnyTimes()
	screen := NewStorageBrowseScreenScreen(mockStorage)

	view := screen.View()
	if !strings.Contains(view, "Operating storage") {
		t.Errorf("View does not contain expected text")
	}
}

func Test_BrowseStorageScreen_updateRows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)

	now := time.Now()
	earlier := now.Add(-time.Hour)
	secrets := []*models.Secret{
		{ID: 2, Title: "Second", SecretType: "Text", CreatedAt: earlier, UpdatedAt: earlier},
		{ID: 1, Title: "First", SecretType: "Card", CreatedAt: now, UpdatedAt: now},
	}

	mockStorage.EXPECT().GetAll(gomock.Any()).Return(secrets, nil).AnyTimes()

	screen := NewStorageBrowseScreenScreen(mockStorage)
	screen.updateRows()

	sortedSecrets := make([]*models.Secret, len(secrets))
	copy(sortedSecrets, secrets)
	sort.Slice(sortedSecrets, func(i, j int) bool {
		return sortedSecrets[i].UpdatedAt.After(sortedSecrets[j].UpdatedAt)
	})

	if len(screen.table.Rows()) != len(sortedSecrets) {
		t.Fatalf("Expected %d rows in table, got %d", len(sortedSecrets), len(screen.table.Rows()))
	}

	for i, sec := range sortedSecrets {
		row := screen.table.Rows()[i]
		if row[0] != strconv.Itoa(int(sec.ID)) || row[1] != sec.Title ||
			row[2] != sec.SecretType || row[3] != sec.CreatedAt.Format("02 Jan 06 15:04") ||
			row[4] != sec.UpdatedAt.Format("02 Jan 06 15:04") {
			t.Errorf("Row %d did not match expected values", i)
		}
	}
}

func Test_BrowseStorageScreen_HelpBindings(t *testing.T) {
	screen := BrowseStorageScreen{}
	bindings := screen.HelpBindings()

	expectedBindings := []struct {
		keys []string
		help string
	}{
		{[]string{"a"}, "add secret"},
		{[]string{"e"}, "edit secret"},
		{[]string{"d"}, "delete secret"},
		{[]string{"c"}, "copy/save secret"},
	}

	if len(bindings) != len(expectedBindings) {
		t.Fatalf("Expected %d bindings, got %d", len(expectedBindings), len(bindings))
	}

	for i, b := range bindings {
		exp := expectedBindings[i]
		if !sliceEqual(b.Keys(), exp.keys) || b.Help().Desc != exp.help {
			t.Errorf("Binding %d did not match expected. Got keys: '%v', help: '%s'; want keys: '%v', help: '%s'",
				i, b.Keys(), b.Help(), exp.keys, exp.help)
		}
	}
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Test_BrowseStorageScreen_handleEdit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)

	type testCase struct {
		name           string
		mockSetup      func()
		mockGet        func()
		mockGetAll     func()
		expectedErrMsg string
	}

	testCases := []testCase{
		{
			name: "Successful_Edit",
			mockSetup: func() {
				mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&models.Secret{ID: 1, SecretType: string(models.TextSecret)}, nil).Times(1)
				mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
			},
			expectedErrMsg: "",
		},
		{
			name: "getSelectedSecret_Failure",
			mockSetup: func() {
				mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("failed to load secret")).Times(1)
				mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
			},
			expectedErrMsg: "failed to load secret",
		},
		{
			name: "getScreenForSecret_Failure",
			mockSetup: func() {
				mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&models.Secret{ID: 1, SecretType: "UnknownSecretType"}, nil).Times(1)
				mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
			},
			expectedErrMsg: "failed to get screen",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			screen := NewStorageBrowseScreenScreen(mockStorage)
			screen.table.SetRows([]table.Row{{"1"}})

			cmd := screen.handleEdit()
			if cmd == nil {
				t.Fatal("Expected a command to be returned")
			}

			if tc.expectedErrMsg != "" {
				err := cmd().(error)
				if err == nil || !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("Expected error containing '%s', got '%v'", tc.expectedErrMsg, err)
				}
			} else {
				if _, ok := cmd().(error); ok {
					t.Fatal("Did not expect an error for successful case")
				}
			}
		})
	}
}

func Test_BrowseStorageScreen_handleCopy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)

	type testCase struct {
		name           string
		mockSetup      func()
		expectedErrMsg string
		expectedMsg    string
		msgType        interface{}
	}

	testCases := []testCase{
		{
			name: "Successful_Copy_TextSecret",
			mockSetup: func() {
				mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&models.Secret{
					ID:         1,
					SecretType: string(models.TextSecret),
					Title:      "Test Secret",
					Text:       &models.Text{Content: "Sample text"},
				}, nil).Times(1)
				mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
			},
			expectedErrMsg: "",
			expectedMsg:    "secret copied successfully",
			msgType:        tui.InfoMsg(""),
		},
		{
			name: "Blob_Secret_Copy",
			mockSetup: func() {
				mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&models.Secret{
					ID:         2,
					SecretType: string(models.BlobSecret),
				}, nil).Times(1)
				mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
			},
			expectedErrMsg: "",
			expectedMsg:    "choose path to save",
			msgType:        tui.PromptMsg{},
		},
		{
			name: "getSelectedSecret_Failure",
			mockSetup: func() {
				mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("failed to load secret")).Times(1)
				mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
			},
			expectedErrMsg: "failed to load secret",
			expectedMsg:    "",
			msgType:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			screen := NewStorageBrowseScreenScreen(mockStorage)
			screen.table.SetRows([]table.Row{{"1"}})

			cmd := screen.handleCopy()
			if cmd == nil {
				t.Fatal("Expected a command to be returned")
			}

			msg := cmd()
			if tc.expectedErrMsg != "" {
				err, ok := msg.(error)
				if !ok {
					t.Fatalf("Expected an error, got %T", msg)
				}
				if !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("Expected error containing '%s', got '%v'", tc.expectedErrMsg, err)
				}
			} else {
				if msgType := tc.msgType; msgType != nil {
					switch msgType.(type) {
					case tui.InfoMsg:
						infoMsg, ok := msg.(tui.InfoMsg)
						if !ok {
							t.Fatalf("Expected a tui.InfoMsg, got %T", msg)
						}
						if !strings.Contains(string(infoMsg), tc.expectedMsg) {
							t.Errorf("Expected message containing '%s', got '%s'", tc.expectedMsg, infoMsg)
						}
					case tui.PromptMsg:
						promptMsg, ok := msg.(tui.PromptMsg)
						if !ok {
							t.Fatalf("Expected a tui.PromptMsg, got %T", msg)
						}
						if !strings.Contains(promptMsg.Prompt, tc.expectedMsg) {
							t.Errorf("Expected prompt message containing '%s', got '%s'", tc.expectedMsg, promptMsg.Prompt)
						}
					}
				}
			}
		})
	}
}

func Test_BrowseStorageScreen_handleDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)

	type testCase struct {
		name           string
		mockSetup      func()
		expectedErrMsg string
		expectedCmdMsg string
	}

	testCases := []testCase{
		{
			name: "Successful_Delete",
			mockSetup: func() {
				mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&models.Secret{
					ID: 1,
				}, nil).Times(1)
				mockStorage.EXPECT().Delete(gomock.Any(), uint64(1)).Return(nil).Times(1)
				mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
			},
			expectedErrMsg: "",
			expectedCmdMsg: "secret deleted",
		},
		{
			name: "getSelectedSecret_Failure",
			mockSetup: func() {
				mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("failed to load secret")).Times(1)
				mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
			},
			expectedErrMsg: "failed to load secret",
			expectedCmdMsg: "",
		},
		{
			name: "Delete_Failure",
			mockSetup: func() {
				mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&models.Secret{
					ID: 2,
				}, nil).Times(1)
				mockStorage.EXPECT().Delete(gomock.Any(), uint64(2)).Return(fmt.Errorf("delete failed")).Times(1)
				mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
			},
			expectedErrMsg: "failed to delete secret",
			expectedCmdMsg: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			screen := NewStorageBrowseScreenScreen(mockStorage)
			screen.table.SetRows([]table.Row{{"1"}})

			cmd := screen.handleDelete()
			if cmd == nil {
				t.Fatal("Expected a command to be returned")
			}

			if tc.expectedErrMsg != "" {
				err, ok := cmd().(error)
				if !ok {
					t.Fatalf("Expected an error, got %T", cmd())
				}
				if !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.expectedErrMsg, err.Error())
				}
				return
			}

			if tc.expectedCmdMsg != "" {
				msg, ok := cmd().(tui.InfoMsg)
				if !ok {
					t.Fatalf("Expected a tui.InfoMsg, got %T", cmd())
				}
				if !strings.Contains(string(msg), tc.expectedCmdMsg) {
					t.Errorf("Expected message containing '%s', got '%s'", tc.expectedCmdMsg, msg)
				}
			}
		})
	}
}

func Test_BrowseStorageScreen_Make(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()

	msg := tui.NavigationMsg{Storage: mockStorage}
	screen := &BrowseStorageScreen{}

	resultScreen, err := screen.Make(msg, 0, 0)

	assert.NoError(t, err, "Make should not return an error")
	assert.IsType(t, &BrowseStorageScreen{}, resultScreen, "Make should return a *BrowseStorageScreen")
}

func Test_BrowseStorageScreen_getScreenForSecret(t *testing.T) {
	screen := &BrowseStorageScreen{}

	tests := []struct {
		name           string
		secret         *models.Secret
		expectedScreen tui.Screen
		expectedError  string
	}{
		{
			name:           "CredentialSecret",
			secret:         &models.Secret{SecretType: string(models.CredSecret)},
			expectedScreen: tui.CredentialEditScreen,
		},
		{
			name:           "TextSecret",
			secret:         &models.Secret{SecretType: string(models.TextSecret)},
			expectedScreen: tui.TextEditScreen,
		},
		{
			name:           "BlobSecret",
			secret:         &models.Secret{SecretType: string(models.BlobSecret)},
			expectedScreen: tui.BlobEditScreen,
		},
		{
			name:           "CardSecret",
			secret:         &models.Secret{SecretType: string(models.CardSecret)},
			expectedScreen: tui.CardEditScreen,
		},
		{
			name:          "UnknownSecret",
			secret:        &models.Secret{SecretType: "unknown"},
			expectedError: "unknown secret type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resultScreen, err := screen.getScreenForSecret(tc.secret)

			if tc.expectedError != "" {
				assert.Error(t, err, "Expected an error")
				assert.Contains(t, err.Error(), tc.expectedError, "Error message should match")
			} else {
				assert.NoError(t, err, "Should not return an error")
				assert.Equal(t, tc.expectedScreen, resultScreen, "Should return the correct screen")
			}
		})
	}
}

func Test_BrowseStorageScreen_colsWidth(t *testing.T) {
	screen := &BrowseStorageScreen{}
	screen.table = prepareTable()

	expectedWidth := tableBorderSize
	for _, col := range screen.table.Columns() {
		expectedWidth += col.Width
	}

	actualWidth := screen.colsWidth()
	assert.Equal(t, expectedWidth, actualWidth, "colsWidth should return the total column width")
}

func Test_BrowseStorageScreen_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()

	type testCase struct {
		name           string
		message        tea.Msg
		mockSetup      func()
		expectedErrMsg string
		expectedCmdMsg string
		expectedScreen tui.Screen
	}

	testCases := []testCase{
		{
			name:    "Reload_Secret_List",
			message: grpc.ReloadSecretList{},
			mockSetup: func() {
				mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
			},
		},
		{
			name: "Save_Path_File_Error",
			message: savePathMsg{
				path:   "/invalid/path",
				secret: &models.Secret{Blob: &models.Blob{FileBytes: []byte("test content")}},
			},
			mockSetup:      func() {},
			expectedErrMsg: "no such file or directory",
		},
		{
			name:    "Key_A",
			message: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")},
			mockSetup: func() {
				mockStorage.EXPECT().GetAll(gomock.Any()).Return([]*models.Secret{}, nil).AnyTimes()
			},
			expectedScreen: tui.SecretTypeScreen,
		},
		{
			name:    "Key_E_Enter",
			message: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")},
			mockSetup: func() {
				mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&models.Secret{
					ID:         1,
					SecretType: string(models.TextSecret),
				}, nil).Times(1)
			},
			expectedScreen: tui.TextEditScreen,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			screen := NewStorageBrowseScreenScreen(mockStorage)
			screen.table.SetRows([]table.Row{{"1"}})

			cmd := screen.Update(tc.message)
			if cmd == nil && tc.expectedCmdMsg != "" {
				t.Fatal("Expected a command to be returned")
			}

			if tc.expectedErrMsg != "" {
				err, ok := cmd().(error)
				if !ok {
					t.Fatalf("Expected an error, got %T", cmd())
				}
				if !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.expectedErrMsg, err.Error())
				}
			} else if tc.expectedCmdMsg != "" {
				switch v := cmd().(type) {
				case tui.InfoMsg:
					if !strings.Contains(string(v), tc.expectedCmdMsg) {
						t.Errorf("Expected message containing '%s', got '%s'", tc.expectedCmdMsg, v)
					}
				case tui.NavigationMsg:
					if v.Screen != tc.expectedScreen {
						t.Errorf("Expected screen %v, got %v", tc.expectedScreen, v.Screen)
					}
				default:
					t.Fatalf("Unexpected command type: %T", v)
				}
			}
		})
	}
}
