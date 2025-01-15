package top

import (
	"beliaev-aa/GophKeeper/internal/client/config"
	"beliaev-aa/GophKeeper/internal/client/tui"
	"beliaev-aa/GophKeeper/internal/client/tui/styles"
	"beliaev-aa/GophKeeper/tests/mocks"
	"errors"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/golang/mock/gomock"
	"reflect"
	"strings"
	"testing"
)

func Test_keyMapToSlice(t *testing.T) {
	type testCase struct {
		name         string
		input        any
		expectedKeys []string
		expectNil    bool
		expectPanic  bool
	}

	testCases := []testCase{
		{
			name: "Valid_Struct_With_Bindings",
			input: struct {
				Binding1 key.Binding
				Binding2 key.Binding
			}{
				Binding1: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action A")),
				Binding2: key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "action B")),
			},
			expectedKeys: []string{"a", "b"},
			expectNil:    false,
		},
		{
			name: "Non_Struct_Input",
			input: []key.Binding{
				key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "action C")),
			},
			expectNil:   true,
			expectPanic: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.expectPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tc.expectPanic {
					t.Error("Expected panic, but got none")
				}
			}()

			result := keyMapToSlice(tc.input)
			if tc.expectNil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
				return
			}

			var keys []string
			for _, b := range result {
				keys = append(keys, b.Keys()...)
			}

			if !reflect.DeepEqual(keys, tc.expectedKeys) {
				t.Errorf("Expected keys %v, got %v", tc.expectedKeys, keys)
			}
		})
	}
}

func Test_removeDuplicateBindings(t *testing.T) {
	type testCase struct {
		name           string
		input          []key.Binding
		expectedOutput []key.Binding
	}

	testCases := []testCase{
		{
			name: "No_Duplicates",
			input: []key.Binding{
				key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action A")),
				key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "action B")),
			},
			expectedOutput: []key.Binding{
				key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action A")),
				key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "action B")),
			},
		},
		{
			name: "With_Duplicates",
			input: []key.Binding{
				key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action A")),
				key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action A (duplicate)")),
				key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "action B")),
				key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action A (another duplicate)")),
			},
			expectedOutput: []key.Binding{
				key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action A")),
				key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "action B")),
			},
		},
		{
			name:           "Empty_Input",
			input:          []key.Binding{},
			expectedOutput: []key.Binding{},
		},
		{
			name: "All_Duplicates",
			input: []key.Binding{
				key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action A")),
				key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action A")),
				key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action A")),
			},
			expectedOutput: []key.Binding{
				key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "action A")),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := removeDuplicateBindings(tc.input)

			if len(result) != len(tc.expectedOutput) {
				t.Fatalf("Expected %d bindings, got %d", len(tc.expectedOutput), len(result))
			}

			for i, binding := range result {
				if !strings.EqualFold(strings.Join(binding.Keys(), " "), strings.Join(tc.expectedOutput[i].Keys(), " ")) {
					t.Errorf("Binding keys mismatch at index %d: expected %v, got %v", i, tc.expectedOutput[i].Keys(), binding.Keys())
				}
				if binding.Help().Desc != tc.expectedOutput[i].Help().Desc {
					t.Errorf("Binding description mismatch at index %d: expected %q, got %q", i, tc.expectedOutput[i].Help().Desc, binding.Help().Desc)
				}
			}
		})
	}
}

func Test_Model_AvailableFooterMsgWidth(t *testing.T) {
	type testCase struct {
		name          string
		width         int
		helpWidget    string
		versionWidget string
		expectedWidth int
	}

	testCases := []testCase{
		{
			name:          "Sufficient_Width",
			width:         100,
			helpWidget:    lipgloss.NewStyle().Width(20).Render("Help"),
			versionWidget: lipgloss.NewStyle().Width(30).Render("Version"),
			expectedWidth: 50,
		},
		{
			name:          "Minimal_Width",
			width:         50,
			helpWidget:    lipgloss.NewStyle().Width(20).Render("Help"),
			versionWidget: lipgloss.NewStyle().Width(20).Render("Version"),
			expectedWidth: 10,
		},
		{
			name:          "Insufficient_Width",
			width:         30,
			helpWidget:    lipgloss.NewStyle().Width(20).Render("Help"),
			versionWidget: lipgloss.NewStyle().Width(20).Render("Version"),
			expectedWidth: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := &Model{
				width:         tc.width,
				helpWidget:    tc.helpWidget,
				versionWidget: tc.versionWidget,
			}

			result := model.availableFooterMsgWidth()
			if result != tc.expectedWidth {
				t.Errorf("Expected width %d, got %d", tc.expectedWidth, result)
			}
		})
	}
}

func Test_Model_ViewHeight(t *testing.T) {
	type testCase struct {
		name           string
		height         int
		mode           mode
		showHelp       bool
		expectedHeight int
	}

	testCases := []testCase{
		{
			name:           "Normal_Mode_Without_Help",
			height:         30,
			mode:           normalMode,
			showHelp:       false,
			expectedHeight: 29,
		},
		{
			name:           "Prompt_Mode_Without_Help",
			height:         30,
			mode:           promptMode,
			showHelp:       false,
			expectedHeight: 26,
		},
		{
			name:           "Normal_Mode_With_Help",
			height:         30,
			mode:           normalMode,
			showHelp:       true,
			expectedHeight: 17,
		},
		{
			name:           "Prompt_Mode_With_Help",
			height:         30,
			mode:           promptMode,
			showHelp:       true,
			expectedHeight: 14,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := &Model{
				height:   tc.height,
				mode:     tc.mode,
				showHelp: tc.showHelp,
			}

			result := model.viewHeight()
			if result != tc.expectedHeight {
				t.Errorf("Expected height %d, got %d", tc.expectedHeight, result)
			}
		})
	}
}

func Test_Model_ViewWidth(t *testing.T) {
	type testCase struct {
		name          string
		width         int
		expectedWidth int
	}

	testCases := []testCase{
		{
			name:          "Width_Above_Min",
			width:         100,
			expectedWidth: 100,
		},
		{
			name:          "Width_Equal_To_Min",
			width:         MinContentWidth,
			expectedWidth: MinContentWidth,
		},
		{
			name:          "Width_Below_Min",
			width:         50,
			expectedWidth: MinContentWidth,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := &Model{
				width: tc.width,
			}

			result := model.viewWidth()
			if result != tc.expectedWidth {
				t.Errorf("Expected width %d, got %d", tc.expectedWidth, result)
			}
		})
	}
}

func Test_Model_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientGRPCInterface(ctrl)
	mockConfig := &config.Config{
		BuildVersion: "1.0.0",
		BuildDate:    "2025-01-15",
		BuildCommit:  "abc123",
	}

	model, err := NewModel(mockConfig, mockClient)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	cmd := model.Init()
	if cmd == nil {
		t.Fatalf("Expected Init to return a command, but got nil")
	}
}

func Test_Model_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGRPCClient := mocks.NewMockClientGRPCInterface(ctrl)

	cfg := &config.Config{
		BuildVersion: "1.0.0",
		BuildDate:    "2025-01-15",
		BuildCommit:  "abc123",
	}

	type testCase struct {
		name         string
		msg          tea.Msg
		initialMode  mode
		expectedMode mode
		expectedHelp bool
		expectedInfo string
		expectedErr  error
	}

	testCases := []testCase{
		{
			name:         "PromptMsg_Changes_Mode_To_PromptMode",
			msg:          tui.PromptMsg{Prompt: "Prompt Content"},
			initialMode:  normalMode,
			expectedMode: promptMode,
		},
		{
			name:         "KeyMsg_Help_Toggles_Help",
			msg:          tea.KeyMsg{Type: tea.KeyCtrlH},
			initialMode:  normalMode,
			expectedHelp: true,
		},
		{
			name:         "KeyMsg_Quit_Triggers_YesNoPrompt",
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")},
			initialMode:  normalMode,
			expectedMode: normalMode,
		},
		{
			name:        "ErrorMsg_Sets_Error",
			msg:         tui.ErrorMsg(errors.New("test error")),
			initialMode: normalMode,
			expectedErr: errors.New("test error"),
		},
		{
			name:         "InfoMsg_Sets_Info",
			msg:          tui.InfoMsg("info message"),
			initialMode:  normalMode,
			expectedInfo: "info message",
		},
		{
			name:        "WindowSizeMsg_Updates_Dimensions",
			msg:         tea.WindowSizeMsg{Width: 100, Height: 50},
			initialMode: normalMode,
		},
		{
			name:         "UnknownMsg_No_State_Change",
			msg:          tea.Msg("unknown"),
			initialMode:  normalMode,
			expectedMode: normalMode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model, _ := NewModel(cfg, mockGRPCClient)
			model.mode = tc.initialMode
			model.showHelp = false
			model.info = ""
			model.err = nil

			updatedModel, _ := model.Update(tc.msg)

			if updatedModel.(*Model).mode != tc.expectedMode {
				t.Errorf("Expected mode %v, got %v", tc.expectedMode, updatedModel.(*Model).mode)
			}

			if updatedModel.(*Model).showHelp != tc.expectedHelp {
				t.Errorf("Expected help state %v, got %v", tc.expectedHelp, updatedModel.(*Model).showHelp)
			}

			if updatedModel.(*Model).info != tc.expectedInfo {
				t.Errorf("Expected info message '%s', got '%s'", tc.expectedInfo, updatedModel.(*Model).info)
			}

			if tc.expectedErr != nil {
				if updatedModel.(*Model).err == nil || updatedModel.(*Model).err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error '%v', got '%v'", tc.expectedErr, updatedModel.(*Model).err)
				}
			} else if updatedModel.(*Model).err != nil {
				t.Errorf("Expected no error, got '%v'", updatedModel.(*Model).err)
			}
		})
	}
}

func Test_Model_View(t *testing.T) {
	cfg := &config.Config{
		BuildVersion: "1.0.0",
		BuildDate:    "2025-01-15",
		BuildCommit:  "abc123",
	}

	mockGRPCClient := mocks.NewMockClientGRPCInterface(nil)

	type testCase struct {
		name           string
		mode           mode
		showHelp       bool
		info           string
		err            error
		expectedFooter string
		initPrompt     bool
	}

	testCases := []testCase{
		{
			name:           "Normal_Mode_No_Help_No_Errors",
			mode:           normalMode,
			showHelp:       false,
			info:           "",
			err:            nil,
			expectedFooter: styles.Padded.Foreground(styles.Black).Background(styles.EvenLighterGrey).Render(""),
		},
		{
			name:           "Prompt_Mode_No_Help_No_Errors",
			mode:           promptMode,
			showHelp:       false,
			info:           "",
			err:            nil,
			expectedFooter: styles.Padded.Foreground(styles.Black).Background(styles.EvenLighterGrey).Render(""),
			initPrompt:     true,
		},
		{
			name:           "Normal_Mode_Help_Displayed",
			mode:           normalMode,
			showHelp:       true,
			info:           "",
			err:            nil,
			expectedFooter: styles.Padded.Foreground(styles.Black).Background(styles.EvenLighterGrey).Render(""),
		},
		{
			name:           "Normal_Mode_Info_Displayed",
			mode:           normalMode,
			showHelp:       false,
			info:           "Operation successful",
			err:            nil,
			expectedFooter: styles.Padded.Foreground(styles.Black).Background(styles.LightGreen).Render("Operation successful"),
		},
		{
			name:           "Normal_Mode_Error_Displayed",
			mode:           normalMode,
			showHelp:       false,
			info:           "",
			err:            errors.New("something went wrong"),
			expectedFooter: styles.Regular.Padding(0, 1).Background(styles.Red).Foreground(styles.White).Render("something went wrong"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model, _ := NewModel(cfg, mockGRPCClient)

			model.mode = tc.mode
			model.showHelp = tc.showHelp
			model.info = tc.info
			model.err = tc.err

			if tc.initPrompt {
				model.prompt, _ = tui.NewPrompt(tui.PromptMsg{Prompt: "Test prompt"})
			}

			view := model.View()

			if !strings.Contains(view, tc.expectedFooter) {
				t.Errorf("Expected footer to contain '%s', got '%s'", tc.expectedFooter, view)
			}
		})
	}
}
