package credentials

import (
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/tests/mocks"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"strings"
	"testing"
)

func Test_CredentialEditScreen_Make(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockSecret := &models.Secret{}

	msg := tui.NavigationMsg{
		Secret:  mockSecret,
		Storage: mockStorage,
	}

	screen := &CredentialEditScreen{}
	result, err := screen.Make(msg, 0, 0)
	if err != nil {
		t.Errorf("Make returned an error: %v", err)
	}

	if _, ok := result.(*CredentialEditScreen); !ok {
		t.Errorf("Expected result to be *CredentialEditScreen, got %T", result)
	}
}

func Test_CredentialEditScreen_Init(t *testing.T) {
	mockSecret := &models.Secret{}
	mockStorage := mocks.NewMockStorage(gomock.NewController(t))

	screen := NewCredentialEditScreen(mockSecret, mockStorage)
	cmd := screen.Init()
	if cmd == nil {
		t.Errorf("Init did not return a valid command")
	}
}

func Test_CredentialEditScreen_Update(t *testing.T) {
	mockSecret := &models.Secret{}
	mockStorage := mocks.NewMockStorage(gomock.NewController(t))

	screen := NewCredentialEditScreen(mockSecret, mockStorage)
	cmd := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	if cmd == nil {
		t.Errorf("Update did not handle input correctly")
	}
}

func Test_CredentialEditScreen_Submit(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStorage := mocks.NewMockStorage(ctrl)
	mockStorage.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockSecret := &models.Secret{}
	screen := NewCredentialEditScreen(mockSecret, mockStorage)
	screen.inputGroup.Inputs[credTitle].SetValue("Test")
	screen.inputGroup.Inputs[credMetadata].SetValue("Test metadata")
	screen.inputGroup.Inputs[credLogin].SetValue("user")
	screen.inputGroup.Inputs[credPassword].SetValue("pass")

	err := screen.Submit()
	if err != nil {
		t.Errorf("Submit failed with error: %v", err)
	}
}

func Test_CredentialEditScreen_View(t *testing.T) {
	mockSecret := &models.Secret{}
	mockStorage := mocks.NewMockStorage(gomock.NewController(t))

	screen := NewCredentialEditScreen(mockSecret, mockStorage)
	view := screen.View()
	expectedContent := "Fill in credential details:"
	if !contains(view, expectedContent) {
		t.Errorf("View output did not contain expected content. Got: %v", view)
	}
}

// Helper function to check if a string contains another string
func contains(input, match string) bool {
	return strings.Contains(input, match)
}
