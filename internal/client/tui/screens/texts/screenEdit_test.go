package texts

import (
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/tests/mocks"
	"github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"strings"
	"testing"
)

func Test_TextEditScreen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)

	type testCase struct {
		name           string
		secret         *models.Secret
		mockSetup      func()
		action         func(screen *TextEditScreen) tea.Cmd
		expectedErrMsg string
		expectedCmd    tui.Screen
	}

	var noScreen tui.Screen

	testCases := []testCase{
		{
			name: "Initialize_with_existing_secret",
			secret: &models.Secret{
				ID:       1,
				Title:    "Existing Title",
				Metadata: "Existing Metadata",
				Text:     &models.Text{Content: "Existing Content"},
			},
			mockSetup: func() {},
			action: func(screen *TextEditScreen) tea.Cmd {
				return screen.Init()
			},
			expectedCmd: noScreen,
		},
		{
			name: "Submit_valid_secret",
			secret: &models.Secret{
				ID:       0,
				Title:    "",
				Metadata: "",
				Text:     &models.Text{},
			},
			mockSetup: func() {
				mockStorage.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			action: func(screen *TextEditScreen) tea.Cmd {
				screen.inputGroup.Inputs[textTitle].SetValue("New Title")
				screen.inputGroup.Inputs[textMetadata].SetValue("New Metadata")
				screen.inputGroup.Inputs[textContent].SetValue("New Content")
				return screen.inputGroup.Buttons[0].Cmd()
			},
			expectedCmd: tui.StorageBrowseScreen,
		},
		{
			name: "Submit_invalid_title",
			secret: &models.Secret{
				ID:       0,
				Title:    "",
				Metadata: "",
				Text:     &models.Text{},
			},
			mockSetup: func() {},
			action: func(screen *TextEditScreen) tea.Cmd {
				screen.inputGroup.Inputs[textMetadata].SetValue("New Metadata")
				screen.inputGroup.Inputs[textContent].SetValue("New Content")
				return screen.inputGroup.Buttons[0].Cmd()
			},
			expectedErrMsg: "please enter title",
		},
		{
			name: "Submit_invalid_metadata",
			secret: &models.Secret{
				ID:       0,
				Title:    "",
				Metadata: "",
				Text:     &models.Text{},
			},
			mockSetup: func() {},
			action: func(screen *TextEditScreen) tea.Cmd {
				screen.inputGroup.Inputs[textTitle].SetValue("New Title")
				screen.inputGroup.Inputs[textContent].SetValue("New Content")
				return screen.inputGroup.Buttons[0].Cmd()
			},
			expectedErrMsg: "please enter metadata",
		},
		{
			name: "Submit_invalid_content",
			secret: &models.Secret{
				ID:       0,
				Title:    "",
				Metadata: "",
				Text:     &models.Text{},
			},
			mockSetup: func() {},
			action: func(screen *TextEditScreen) tea.Cmd {
				screen.inputGroup.Inputs[textTitle].SetValue("New Title")
				screen.inputGroup.Inputs[textMetadata].SetValue("New Metadata")
				return screen.inputGroup.Buttons[0].Cmd()
			},
			expectedErrMsg: "please enter content",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			screen := NewTextEditScreen(tc.secret, mockStorage)
			cmd := tc.action(screen)

			if cmd != nil {
				if tc.expectedErrMsg != "" {
					err, ok := cmd().(error)
					if !ok {
						t.Fatalf("Expected error, got %T", cmd())
					}
					if err.Error() != tc.expectedErrMsg {
						t.Errorf("Expected error message '%s', got '%s'", tc.expectedErrMsg, err.Error())
					}
				} else if tc.expectedCmd != noScreen {
					switch v := cmd().(type) {
					case tui.NavigationMsg:
						if v.Screen != tc.expectedCmd {
							t.Errorf("Expected screen %v, got %v", tc.expectedCmd, v.Screen)
						}
					default:
						t.Fatalf("Unexpected command type: %T", v)
					}
				}
			}
		})
	}
}

func Test_TextEditScreen_Make(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	secret := &models.Secret{
		ID:       1,
		Title:    "Test Title",
		Metadata: "Test Metadata",
		Text:     &models.Text{Content: "Test Content"},
	}

	msg := tui.NavigationMsg{Secret: secret, Storage: mockStorage}
	screen := TextEditScreen{}

	resultScreen, err := screen.Make(msg, 0, 0)
	if err != nil {
		t.Fatalf("Make returned an error: %v", err)
	}

	if _, ok := resultScreen.(*TextEditScreen); !ok {
		t.Errorf("Make did not return a *TextEditScreen, got %T", resultScreen)
	}
}

func Test_TextEditScreen_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	secret := &models.Secret{
		ID:       1,
		Title:    "Test Title",
		Metadata: "Test Metadata",
		Text:     &models.Text{Content: "Test Content"},
	}

	screen := NewTextEditScreen(secret, mockStorage)

	type testCase struct {
		name          string
		message       tea.Msg
		expectedState string
	}

	testCases := []testCase{
		{
			name: "Update_TextInput",
			message: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune("A"),
			},
			expectedState: "Test TitleA",
		},
		{
			name: "Handle_Unknown_Message",
			message: struct {
				msg string
			}{msg: "unknown"},
			expectedState: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			screen.Update(tc.message)

			if tc.expectedState != "" {
				state := screen.inputGroup.Inputs[textTitle].Value()
				if state != tc.expectedState {
					t.Errorf("Expected state '%s', got '%s'", tc.expectedState, state)
				}
			}
		})
	}
}

func Test_TextEditScreen_View(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	secret := &models.Secret{
		ID:       1,
		Title:    "Test Title",
		Metadata: "Test Metadata",
		Text:     &models.Text{Content: "Test Content"},
	}

	screen := NewTextEditScreen(secret, mockStorage)

	view := screen.View()
	expectedContent := "Fill in text details:"
	if !contains(view, expectedContent) {
		t.Errorf("View did not contain expected content '%s'", expectedContent)
	}
}

func contains(output, content string) bool {
	return strings.Contains(output, content)
}
