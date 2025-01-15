package blobs

import (
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/pkg/models"
	"beliaev-aa/GophKeeper/tests/mocks"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testCase struct {
	name      string
	setupMock func(mockStorage *mocks.MockStorage)
	inputs    func(screen *BlobEditScreen)
	filePath  string
	expectErr string
}

func TestBlobEditScreen_Submit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)

	tests := []testCase{
		{
			name: "Submit_Success_Create",
			setupMock: func(mockStorage *mocks.MockStorage) {
				mockStorage.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			inputs: func(screen *BlobEditScreen) {
				screen.inputGroup.Inputs[blobTitle].SetValue("Test Title")
				screen.inputGroup.Inputs[blobMetadata].SetValue("Test Metadata")
			},
			filePath:  "test_file.txt",
			expectErr: "",
		},
		{
			name: "Submit_Success_Update",
			setupMock: func(mockStorage *mocks.MockStorage) {
				mockStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			inputs: func(screen *BlobEditScreen) {
				screen.inputGroup.Inputs[blobTitle].SetValue("Updated Title")
				screen.inputGroup.Inputs[blobMetadata].SetValue("Updated Metadata")
				screen.secret.ID = 1
			},
			filePath:  "test_file.txt",
			expectErr: "",
		},
		{
			name:      "Submit_Empty_Title",
			setupMock: func(mockStorage *mocks.MockStorage) {},
			inputs: func(screen *BlobEditScreen) {
				screen.inputGroup.Inputs[blobTitle].SetValue("")
				screen.inputGroup.Inputs[blobMetadata].SetValue("Test Metadata")
			},
			filePath:  "",
			expectErr: "please enter title",
		},
		{
			name:      "Submit_Empty_Metadata",
			setupMock: func(mockStorage *mocks.MockStorage) {},
			inputs: func(screen *BlobEditScreen) {
				screen.inputGroup.Inputs[blobTitle].SetValue("Test Title")
				screen.inputGroup.Inputs[blobMetadata].SetValue("")
			},
			filePath:  "",
			expectErr: "please enter metadata",
		},
		{
			name:      "Submit_ReadFile_Error",
			setupMock: func(mockStorage *mocks.MockStorage) {},
			inputs: func(screen *BlobEditScreen) {
				screen.inputGroup.Inputs[blobTitle].SetValue("Test Title")
				screen.inputGroup.Inputs[blobMetadata].SetValue("Test Metadata")
			},
			filePath:  "nonexistent_file.txt",
			expectErr: "open nonexistent_file.txt: no such file or directory",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем временный файл
			content := []byte("test content")
			err := os.WriteFile("test_file.txt", content, 0644)
			assert.NoError(t, err)
			defer os.Remove("test_file.txt")

			secret := &models.Secret{}
			screen := NewBlobEditScreen(secret, mockStorage)

			if tc.inputs != nil {
				tc.inputs(screen)
			}
			if tc.setupMock != nil {
				tc.setupMock(mockStorage)
			}

			err = screen.Submit(tc.filePath)
			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBlobEditScreen_ValidateInputs(t *testing.T) {
	tests := []testCase{
		{
			name: "ValidateInputs_Success",
			inputs: func(screen *BlobEditScreen) {
				screen.inputGroup.Inputs[blobTitle].SetValue("Test Title")
				screen.inputGroup.Inputs[blobMetadata].SetValue("Test Metadata")
			},
			expectErr: "",
		},
		{
			name: "ValidateInputs_Empty_Title",
			inputs: func(screen *BlobEditScreen) {
				screen.inputGroup.Inputs[blobTitle].SetValue("")
				screen.inputGroup.Inputs[blobMetadata].SetValue("Test Metadata")
			},
			expectErr: "please enter title",
		},
		{
			name: "ValidateInputs_Empty_Metadata",
			inputs: func(screen *BlobEditScreen) {
				screen.inputGroup.Inputs[blobTitle].SetValue("Test Title")
				screen.inputGroup.Inputs[blobMetadata].SetValue("")
			},
			expectErr: "please enter metadata",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			screen := NewBlobEditScreen(&models.Secret{}, nil)

			if tc.inputs != nil {
				tc.inputs(screen)
			}

			err := screen.validateInputs()
			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBlobEditScreen_View(t *testing.T) {
	screen := NewBlobEditScreen(&models.Secret{}, nil)
	view := screen.View()
	assert.Contains(t, view, "Fill in file details:")
}

func TestReadFileToBytes(t *testing.T) {
	t.Run("ReadFile_Success", func(t *testing.T) {
		content := []byte("test content")
		err := os.WriteFile("test_file.txt", content, 0644)
		assert.NoError(t, err)
		defer os.Remove("test_file.txt")

		data, err := readFileToBytes("test_file.txt")
		assert.NoError(t, err)
		assert.Equal(t, content, data)
	})

	t.Run("ReadFile_Error", func(t *testing.T) {
		_, err := readFileToBytes("nonexistent_file.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})
}

func TestBlobEditScreen_Init(t *testing.T) {
	screen := NewBlobEditScreen(&models.Secret{}, nil)
	cmd := screen.Init()

	assert.NotNil(t, cmd)
}

func TestBlobEditScreen_Update(t *testing.T) {
	tests := []struct {
		name      string
		msg       tea.Msg
		setup     func(screen *BlobEditScreen)
		assertion func(t *testing.T, screen *BlobEditScreen, cmd tea.Cmd)
	}{
		{
			name:  "Update_FocusNext",
			msg:   tea.KeyMsg{Type: tea.KeyDown},
			setup: func(screen *BlobEditScreen) {},
			assertion: func(t *testing.T, screen *BlobEditScreen, cmd tea.Cmd) {
				assert.NotNil(t, cmd)
				assert.Equal(t, 1, screen.inputGroup.FocusIndex)
			},
		},
		{
			name: "Update_Submit",
			msg:  tea.KeyMsg{Type: tea.KeyEnter},
			setup: func(screen *BlobEditScreen) {
				screen.inputGroup.FocusIndex = 2 // Фокус на кнопке
			},
			assertion: func(t *testing.T, screen *BlobEditScreen, cmd tea.Cmd) {
				assert.NotNil(t, cmd)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			screen := NewBlobEditScreen(&models.Secret{}, nil)
			if tc.setup != nil {
				tc.setup(screen)
			}

			cmd := screen.Update(tc.msg)
			if tc.assertion != nil {
				tc.assertion(t, screen, cmd)
			}
		})
	}
}

func TestBlobEditScreen_Make(t *testing.T) {
	msg := tui.NavigationMsg{
		Secret:  &models.Secret{},
		Storage: nil,
	}

	screen := &BlobEditScreen{}
	result, err := screen.Make(msg, 0, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.IsType(t, &BlobEditScreen{}, result)
}
