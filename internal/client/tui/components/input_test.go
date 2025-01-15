package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCase struct {
	name      string
	setup     func() InputGroup
	expectErr string
	testFunc  func(t *testing.T, inputGroup InputGroup)
}

func TestInputGroup_Init(t *testing.T) {
	tests := []testCase{
		{
			name: "Init_Success",
			setup: func() InputGroup {
				inputs := []textinput.Model{
					textinput.New(),
					textinput.New(),
				}
				return NewInputGroup(inputs, nil)
			},
			testFunc: func(t *testing.T, inputGroup InputGroup) {
				cmd := inputGroup.Init()
				assert.NotNil(t, cmd)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inputGroup := tc.setup()
			tc.testFunc(t, inputGroup)
		})
	}
}

func TestInputGroup_Update(t *testing.T) {
	tests := []testCase{
		{
			name: "Update_FocusChange",
			setup: func() InputGroup {
				input1 := textinput.New()
				input1.Placeholder = "Input 1"
				input2 := textinput.New()
				input2.Placeholder = "Input 2"

				inputs := []textinput.Model{input1, input2}
				return NewInputGroup(inputs, nil)
			},
			testFunc: func(t *testing.T, inputGroup InputGroup) {
				msg := tea.KeyMsg{Type: tea.KeyDown}
				newModel, _ := inputGroup.Update(msg)
				updatedGroup := newModel.(InputGroup)

				assert.Equal(t, 1, updatedGroup.FocusIndex)
			},
		},
		{
			name: "Update_ButtonCommand",
			setup: func() InputGroup {
				input1 := textinput.New()
				input1.Placeholder = "Input 1"

				button := Button{
					Title: "Submit",
					Cmd:   func() tea.Cmd { return func() tea.Msg { return "Command Executed" } },
				}

				inputs := []textinput.Model{input1}
				buttons := []Button{button}
				return NewInputGroup(inputs, buttons)
			},
			testFunc: func(t *testing.T, inputGroup InputGroup) {
				inputGroup.FocusIndex = 1 // Фокус на кнопке
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				newModel, cmd := inputGroup.Update(msg)
				updatedGroup := newModel.(InputGroup)

				assert.NotNil(t, cmd)
				assert.Equal(t, 1, updatedGroup.FocusIndex)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inputGroup := tc.setup()
			tc.testFunc(t, inputGroup)
		})
	}
}

func TestInputGroup_View(t *testing.T) {
	tests := []testCase{
		{
			name: "View_Render",
			setup: func() InputGroup {
				input1 := textinput.New()
				input1.Placeholder = "Input 1"
				input2 := textinput.New()
				input2.Placeholder = "Input 2"

				inputs := []textinput.Model{input1, input2}
				buttons := []Button{{Title: "Submit"}}
				return NewInputGroup(inputs, buttons)
			},
			testFunc: func(t *testing.T, inputGroup InputGroup) {
				view := inputGroup.View()
				assert.Contains(t, view, "Input 1")
				assert.Contains(t, view, "Input 2")
				assert.Contains(t, view, "Submit")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inputGroup := tc.setup()
			tc.testFunc(t, inputGroup)
		})
	}
}
