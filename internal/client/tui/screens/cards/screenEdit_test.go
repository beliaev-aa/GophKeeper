package cards

import (
	"beliaev-aa/GophKeeper/internal/client/storage"
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/tests/mocks"
	"github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"strings"
	"testing"
)

func Test_CardEditScreen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockSecret := &models.Secret{
		ID: 1,
		Card: &models.Card{
			Number:   "1234567890123456",
			ExpYear:  2023,
			ExpMonth: 12,
			CVV:      123,
		},
		Title:    "Test Card",
		Metadata: "Test Metadata",
	}

	tests := []struct {
		name          string
		secret        *models.Secret
		storage       storage.Storage
		inputs        []tea.Msg
		expectedError string
	}{
		{
			name:    "Initialize_with_secret",
			secret:  mockSecret,
			storage: mockStorage,
		},
		{
			name:    "Submit_empty_title",
			secret:  &models.Secret{},
			storage: mockStorage,
			inputs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyEnter},
			},
			expectedError: "please enter metadata",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			screen := NewCardEditScreen(tc.secret, tc.storage)
			for _, input := range tc.inputs {
				screen.Update(input)
			}
			err := screen.Submit()
			if (err != nil && err.Error() != tc.expectedError) || (err == nil && tc.expectedError != "") {
				t.Errorf("Test %v failed: expected error '%v', got '%v'", tc.name, tc.expectedError, err)
			}
		})
	}
}

func Test_CardEditScreen_View(t *testing.T) {
	mockSecret := &models.Secret{
		ID: 1,
		Card: &models.Card{
			Number:   "1234567890123456",
			ExpYear:  2023,
			ExpMonth: 12,
			CVV:      123,
		},
		Title:    "Test Card",
		Metadata: "Test Metadata",
	}

	mockStorage := mocks.NewMockStorage(gomock.NewController(t))
	screen := NewCardEditScreen(mockSecret, mockStorage)

	viewOutput := screen.View()
	expectedContent := "Fill in card details:"
	if !strings.Contains(viewOutput, expectedContent) {
		t.Errorf("View test failed: Expected to find '%v'", expectedContent)
	}
}

func Test_CardEditScreen_Make(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	testSecret := &models.Secret{
		Title: "Test Card",
	}

	msg := tui.NavigationMsg{
		Secret:  testSecret,
		Storage: mockStorage,
	}

	screen := &CardEditScreen{}
	result, err := screen.Make(msg, 0, 0)
	if err != nil {
		t.Errorf("Make method returned an error: %v", err)
	}

	if _, ok := result.(*CardEditScreen); !ok {
		t.Errorf("Expected result to be *CardEditScreen, got %T", result)
	}
}

func Test_CardEditScreen_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	testSecret := &models.Secret{
		Title: "Test Card",
	}

	screen := NewCardEditScreen(testSecret, mockStorage)
	cmd := screen.Init()

	if cmd == nil {
		t.Error("Init should return a non-nil command")
	}
}
