package remotes

import (
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/tests/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_RemoteOpenScreen_Make(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)
	maker := RemoteOpenScreenMaker{Client: mockClient}

	screen, err := maker.Make(tui.NavigationMsg{}, 0, 0)
	if err != nil {
		t.Errorf("Make failed with error: %v", err)
	}
	if _, ok := screen.(*RemoteOpenScreen); !ok {
		t.Errorf("Expected *RemoteOpenScreen, got %T", screen)
	}
}

func Test_RemoteOpenScreen_InstanceMake(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)

	msg := tui.NavigationMsg{
		Client: mockClient,
	}

	screenInstance := &RemoteOpenScreen{}

	resultScreen, err := screenInstance.Make(msg, 0, 0)

	if err != nil {
		t.Errorf("Make returned an unexpected error: %v", err)
	}

	if _, ok := resultScreen.(*RemoteOpenScreen); !ok {
		t.Errorf("Expected result to be *RemoteOpenScreen, got %T", resultScreen)
	}

	resultScreenInstance, _ := resultScreen.(*RemoteOpenScreen)
	if resultScreenInstance.client != mockClient {
		t.Errorf("The client in the result should be the mock client used in the test")
	}
}

func Test_RemoteOpenScreen_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)
	mockClient.EXPECT().GetPassword().Return("valid_password").AnyTimes()

	tests := []struct {
		name        string
		token       string
		expectedCmd string
		expectError bool
	}{
		{
			name:        "With_valid_token",
			token:       "valid_token",
			expectedCmd: "StorageBrowseScreen",
			expectError: false,
		},
		{
			name:        "With_invalid_token",
			token:       "",
			expectedCmd: "LoginScreen",
			expectError: false,
		},
		{
			name:        "Storage_initialization_error",
			token:       "valid_token",
			expectedCmd: "Error",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient.EXPECT().GetToken().Return(tc.token).Times(1)
			screen := NewRemoteOpenScreen(mockClient)
			_ = screen.Init()
		})
	}
}

func Test_RemoteOpenScreen_Update(t *testing.T) {
	mockClient := mocks.NewMockClientGRPCInterface(gomock.NewController(t))
	screen := NewRemoteOpenScreen(mockClient)
	cmd := screen.Update(nil)
	if cmd != nil {
		t.Errorf("Expected cmd to be nil")
	}
}

func Test_RemoteOpenScreen_View(t *testing.T) {
	mockClient := mocks.NewMockClientGRPCInterface(gomock.NewController(t))
	screen := NewRemoteOpenScreen(mockClient)
	output := screen.View()
	if output != "" {
		t.Errorf("Expected empty string from View, got '%v'", output)
	}
}
