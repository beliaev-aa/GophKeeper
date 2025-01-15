package secrets

import (
	"bytes"
	"github.com/charmbracelet/bubbles/list"
	"strings"
	"testing"
)

func Test_secretItemDelegate_Render(t *testing.T) {
	delegate := secretItemDelegate{}
	items := []list.Item{
		secretItem{id: 1, name: "First Secret"},
		secretItem{id: 2, name: "Second Secret"},
	}

	listModel := list.New(items, delegate, 20, 5)
	listModel.SetShowStatusBar(false)

	tests := []struct {
		name           string
		setupModel     func(m list.Model) list.Model
		index          int
		expectedOutput string
	}{
		{
			name: "Render_non_selected_item",
			setupModel: func(m list.Model) list.Model {
				m.Select(1)
				return m
			},
			index:          0,
			expectedOutput: "    First Secret",
		},
		{
			name: "Render_selected_item",
			setupModel: func(m list.Model) list.Model {
				m.Select(0)
				return m
			},
			index:          0,
			expectedOutput: "> First Secret",
		},
		{
			name: "Render_non_selected_second_item",
			setupModel: func(m list.Model) list.Model {
				m.Select(1)
				return m
			},
			index:          0,
			expectedOutput: "    First Secret",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			modifiedModel := tc.setupModel(listModel)
			var buf bytes.Buffer
			delegate.Render(&buf, modifiedModel, tc.index, items[tc.index])

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func Test_secretItem_FilterValue(t *testing.T) {
	item := secretItem{id: 1, name: "Sample Secret"}
	if value := item.FilterValue(); value != "" {
		t.Errorf("FilterValue expected to return empty string, got %q", value)
	}
}

func Test_secretItemDelegate_Update(t *testing.T) {
	delegate := secretItemDelegate{}
	model := list.New(nil, delegate, 0, 0)

	if cmd := delegate.Update(nil, &model); cmd != nil {
		t.Errorf("Update expected to return nil")
	}
}
