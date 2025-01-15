package secrets

import (
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/tests/mocks"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"strings"
	"testing"
)

func Test_SecretTypeScreen_Make(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	msg := tui.NavigationMsg{Storage: mockStorage}

	screen := &SecretTypeScreen{}
	resultScreen, err := screen.Make(msg, 0, 0)
	if err != nil {
		t.Errorf("Make returned an error: %v", err)
	}

	if _, ok := resultScreen.(*SecretTypeScreen); !ok {
		t.Errorf("Expected result to be *SecretTypeScreen, got %T", resultScreen)
	}
}

func Test_SecretTypeScreen_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	screen := NewSecretTypeScreen(mockStorage)

	if screen.storage != mockStorage {
		t.Errorf("Storage not set correctly in NewSecretTypeScreen")
	}
}

func Test_SecretTypeScreen_PrepareSecretListModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	screen := NewSecretTypeScreen(mockStorage)
	screen.prepareSecretListModel()

	if len(screen.list.Items()) != 5 {
		t.Errorf("Expected 5 items in the list, got %d", len(screen.list.Items()))
	}
}

func Test_SecretTypeScreen_View(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	screen := NewSecretTypeScreen(mockStorage)

	view := screen.View()
	expectedContent := "Select type of secret:"
	if !strings.Contains(view, expectedContent) {
		t.Errorf("View did not contain expected content. Got: %v", view)
	}
}

func Test_SecretTypeScreen_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	screen := NewSecretTypeScreen(mockStorage)
	screen.prepareSecretListModel()

	tests := []struct {
		name         string
		selectIndex  int
		expectedFunc string
	}{
		{name: "Select_Back", selectIndex: selectBack, expectedFunc: "SetBodyPane to StorageBrowseScreen"},
		{name: "Select_Credentials", selectIndex: selectCredential, expectedFunc: "SetBodyPane"},
		{name: "Select_Text", selectIndex: selectText, expectedFunc: "SetBodyPane"},
		{name: "Select_Card", selectIndex: selectCard, expectedFunc: "SetBodyPane"},
		{name: "Select_Blob", selectIndex: selectBlob, expectedFunc: "SetBodyPane"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			screen.list.Select(tc.selectIndex)

			msg := tea.KeyMsg{Type: tea.KeyEnter}
			cmd := screen.Update(msg)

			if cmd == nil {
				t.Fatalf("%s: Expected a command, got nil", tc.name)
			}
		})
	}
}

func Test_SecretTypeScreen_Update_WindowSizeMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	screen := NewSecretTypeScreen(mockStorage)
	screen.prepareSecretListModel()

	testWidth := 100
	testHeight := 40

	windowSizeMsg := tea.WindowSizeMsg{
		Width:  testWidth,
		Height: testHeight,
	}

	screen.Update(windowSizeMsg)

	if screen.list.Width() != testWidth {
		t.Errorf("List width was not set correctly: got %d, want %d", screen.list.Width(), testWidth)
	}
	if screen.list.Height() != testHeight-4 {
		t.Errorf("List height was not set correctly: got %d, want %d", screen.list.Height(), testHeight-4)
	}
}

func Test_SecretTypeScreen_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	screen := NewSecretTypeScreen(mockStorage)

	cmd := screen.Init()

	if cmd == nil {
		t.Fatal("Init did not return a command")
	}
}
