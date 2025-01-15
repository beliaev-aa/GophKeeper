package auth

import (
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/internal/client/tui/components"
	"beliaev-aa/GophKeeper/tests/mocks"
	"context"
	"errors"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCase struct {
	name      string
	setupMock func(client *mocks.MockClientGRPCInterface)
	mode      Mode
	login     string
	password  string
	expectErr string
}

func TestAuthenticateScreen_Submit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockClientGRPCInterface(ctrl)

	tests := []testCase{
		{
			name: "Submit_Login_Success",
			setupMock: func(client *mocks.MockClientGRPCInterface) {
				client.EXPECT().Login(context.Background(), "test", "password").Return("test-token", nil).Times(1)
				client.EXPECT().SetToken("test-token").Times(1)
				client.EXPECT().SetPassword("password").Times(1)
				client.EXPECT().GetPassword().Return("password").AnyTimes()
			},
			mode:      modeLogin,
			login:     "test",
			password:  "password",
			expectErr: "",
		},
		{
			name: "Submit_Login_Error",
			setupMock: func(client *mocks.MockClientGRPCInterface) {
				client.EXPECT().Login(context.Background(), "test", "password").Return("", errors.New("login error")).Times(1)
			},
			mode:      modeLogin,
			login:     "test",
			password:  "password",
			expectErr: "login error",
		},
		{
			name: "Submit_Register_Success",
			setupMock: func(client *mocks.MockClientGRPCInterface) {
				client.EXPECT().Register(context.Background(), "test", "password").Return("test-token", nil).Times(1)
				client.EXPECT().SetToken("test-token").Times(1)
				client.EXPECT().SetPassword("password").Times(1)
			},
			mode:      modeRegister,
			login:     "test",
			password:  "password",
			expectErr: "",
		},
		{
			name: "Submit_Register_Error",
			setupMock: func(client *mocks.MockClientGRPCInterface) {
				client.EXPECT().Register(context.Background(), "test", "password").Return("", errors.New("registration error")).Times(1)
			},
			mode:      modeRegister,
			login:     "test",
			password:  "password",
			expectErr: "registration error",
		},
		{
			name:      "Submit_Empty_Login",
			setupMock: func(client *mocks.MockClientGRPCInterface) {},
			mode:      modeLogin,
			login:     "",
			password:  "password",
			expectErr: "please enter login",
		},
		{
			name:      "Submit_Empty_Password",
			setupMock: func(client *mocks.MockClientGRPCInterface) {},
			mode:      modeLogin,
			login:     "test",
			password:  "",
			expectErr: "please enter password",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			screen := NewLoginScreen(client)
			screen.inputGroup.Inputs[0].SetValue(tc.login)
			screen.inputGroup.Inputs[1].SetValue(tc.password)

			if tc.setupMock != nil {
				tc.setupMock(client)
			}

			cmd := screen.Submit(tc.mode)
			if cmd != nil {
				msg := cmd()
				switch msg := msg.(type) {
				case error:
					if tc.expectErr != "" {
						assert.Contains(t, msg.Error(), tc.expectErr)
					} else {
						assert.NoError(t, msg)
					}
				default:
					assert.Empty(t, tc.expectErr)
				}
			} else {
				assert.Empty(t, tc.expectErr)
			}
		})
	}
}

func TestAuthenticateScreen_View(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *AuthenticateScreen
		expectStr string
	}{
		{
			name: "view_renders_correctly",
			setup: func() *AuthenticateScreen {
				inputs := []textinput.Model{
					newInput(inputOpts{placeholder: "Login"}),
					newInput(inputOpts{placeholder: "Password"}),
				}

				buttons := []components.Button{
					{Title: "[ Login ]", Cmd: nil},
					{Title: "[ Register ]", Cmd: nil},
				}

				inputGroup := components.NewInputGroup(inputs, buttons)
				return &AuthenticateScreen{
					inputGroup: inputGroup,
				}
			},
			expectStr: "Fill in credentials:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			screen := tc.setup()
			view := screen.View()

			assert.Contains(t, view, tc.expectStr)
			assert.Contains(t, view, "Login")
			assert.Contains(t, view, "Password")
			assert.Contains(t, view, "[ Login ]")
			assert.Contains(t, view, "[ Register ]")
		})
	}
}

func TestAuthenticateScreen_Init(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *AuthenticateScreen
	}{
		{
			name: "init_success",
			setup: func() *AuthenticateScreen {
				inputs := []textinput.Model{
					newInput(inputOpts{placeholder: "Login"}),
					newInput(inputOpts{placeholder: "Password"}),
				}

				buttons := []components.Button{
					{Title: "[ Login ]", Cmd: nil},
					{Title: "[ Register ]", Cmd: nil},
				}

				inputGroup := components.NewInputGroup(inputs, buttons)
				return &AuthenticateScreen{
					inputGroup: inputGroup,
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			screen := tc.setup()
			cmd := screen.Init()
			assert.NotNil(t, cmd)
		})
	}
}

func TestAuthenticateScreen_Update(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *AuthenticateScreen
		inputMsg  string
		focusIdx  int
		expectIdx int
	}{
		{
			name: "update_focus_changes_on_down_key",
			setup: func() *AuthenticateScreen {
				inputs := []textinput.Model{
					newInput(inputOpts{placeholder: "Login"}),
					newInput(inputOpts{placeholder: "Password"}),
				}

				buttons := []components.Button{
					{Title: "[ Login ]", Cmd: nil},
					{Title: "[ Register ]", Cmd: nil},
				}

				inputGroup := components.NewInputGroup(inputs, buttons)
				return &AuthenticateScreen{
					inputGroup: inputGroup,
				}
			},
			inputMsg:  "down",
			focusIdx:  0,
			expectIdx: 1,
		},
		{
			name: "update_focus_wraps_to_top",
			setup: func() *AuthenticateScreen {
				inputs := []textinput.Model{
					newInput(inputOpts{placeholder: "Login"}),
					newInput(inputOpts{placeholder: "Password"}),
				}

				buttons := []components.Button{
					{Title: "[ Login ]", Cmd: nil},
					{Title: "[ Register ]", Cmd: nil},
				}

				inputGroup := components.NewInputGroup(inputs, buttons)
				return &AuthenticateScreen{
					inputGroup: inputGroup,
				}
			},
			inputMsg:  "down",
			focusIdx:  3,
			expectIdx: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			screen := tc.setup()
			screen.inputGroup.FocusIndex = tc.focusIdx

			cmd := screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			assert.NotNil(t, cmd)
			assert.Equal(t, tc.expectIdx, screen.inputGroup.FocusIndex)
		})
	}
}

func TestAuthenticateScreen_Make(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)

	navigationMsg := tui.NavigationMsg{
		Client: mockClient,
	}

	screen := &AuthenticateScreen{}
	result, err := screen.Make(navigationMsg, 0, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	_, ok := result.(*AuthenticateScreen)
	assert.True(t, ok, "Expected result to be of type *AuthenticateScreen")
}
