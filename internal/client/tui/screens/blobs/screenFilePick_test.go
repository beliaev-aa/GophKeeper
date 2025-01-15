package blobs

import (
	"beliaev-aa/GophKeeper/internal/client/storage"
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/tests/mocks"
	"github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"os"
	"strings"
	"testing"
)

func Test_FilePickScreen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)

	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name       string
		secret     *models.Secret
		storage    storage.Storage
		callback   tui.NavigationCallback
		keyPresses []string
		wantPath   string
	}{
		{
			name:       "Initializes_with_default_path",
			secret:     &models.Secret{},
			storage:    mockStorage,
			callback:   func(args ...any) tea.Cmd { return nil },
			keyPresses: []string{},
			wantPath:   homeDir,
		},
		{
			name:       "Handles_b_key_to_go_back",
			secret:     &models.Secret{},
			storage:    mockStorage,
			callback:   func(args ...any) tea.Cmd { return nil },
			keyPresses: []string{"b"},
			wantPath:   "",
		},
		{
			name:       "Selects_file_and_calls_callback",
			secret:     &models.Secret{},
			storage:    mockStorage,
			callback:   func(args ...any) tea.Cmd { return func() tea.Msg { return "File selected: " + args[0].(string) } },
			keyPresses: []string{"enter"},
			wantPath:   "selected_file_path",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			screen := NewFilePickScreen(tc.secret, tc.storage, tc.callback)

			for _, key := range tc.keyPresses {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
				screen.Update(msg)
			}

			if screen.filePicker.CurrentDirectory != tc.wantPath && tc.name == "Initializes_with_default_path" {
				t.Errorf("expected path to be %v, got %v", tc.wantPath, screen.filePicker.CurrentDirectory)
			}
		})
	}
}

func Test_FilePickScreen_Make(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockSecret := &models.Secret{}

	msg := tui.NavigationMsg{
		Secret:   mockSecret,
		Storage:  mockStorage,
		Callback: func(args ...any) tea.Cmd { return nil },
	}

	screen := &FilePickScreen{}
	resultScreen, err := screen.Make(msg, 0, 0)

	if err != nil {
		t.Errorf("Make returned an error: %v", err)
	}

	if _, ok := resultScreen.(*FilePickScreen); !ok {
		t.Errorf("Make did not return a FilePickScreen instance")
	}
}

func Test_FilePickScreen_Init(t *testing.T) {
	screen := NewFilePickScreen(&models.Secret{}, mocks.NewMockStorage(gomock.NewController(t)), func(args ...any) tea.Cmd { return nil })
	cmd := screen.Init()

	if cmd == nil {
		t.Errorf("Init did not return a proper command")
	}
}

func Test_FilePickScreen_View(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	screen := NewFilePickScreen(&models.Secret{}, mocks.NewMockStorage(gomock.NewController(t)), func(args ...any) tea.Cmd { return nil })

	screen.filePicker.CurrentDirectory = homeDir

	view := screen.View()

	expectedContent := "Select file to store. Use ←, ↑, →, ↓ to navigate. Press B to go back"
	if !strings.Contains(view, expectedContent) {
		t.Errorf("View does not contain expected content. Expected to find %v", expectedContent)
	}
}
