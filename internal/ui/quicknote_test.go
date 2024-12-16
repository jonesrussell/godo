package ui_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/godo/internal/testutil"
	"github.com/jonesrussell/godo/internal/ui"
)

func TestQuickNoteUI(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTitle string
		wantErr   bool
	}{
		{
			name:      "Valid input",
			input:     "Test todo",
			wantTitle: "Test todo",
			wantErr:   false,
		},
		{
			name:      "Empty input",
			input:     "",
			wantTitle: "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new mock service for each test case
			mockService := testutil.NewMockTodoService()
			quickNote := ui.NewQuickNote(mockService)
			mock := testutil.AsMockTodoService(mockService)
			mock.SetShouldError(tt.wantErr)

			// Simulate typing
			for _, r := range tt.input {
				model, _ := quickNote.Update(tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{r},
				})
				quickNote = model.(*ui.QuickNoteUI)
			}

			// Simulate Enter
			_, _ = quickNote.Update(tea.KeyMsg{Type: tea.KeyEnter})

			if got := mock.GetLastTitle(); got != tt.wantTitle {
				t.Errorf("Got title %q, want %q", got, tt.wantTitle)
			}
		})
	}
}

func TestQuickNoteUI_Escape(t *testing.T) {
	mockService := testutil.NewMockTodoService()
	quickNote := ui.NewQuickNote(mockService)

	// Simulate Escape key
	model, cmd := quickNote.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if model == nil {
		t.Error("Expected model, got nil")
	}
	if cmd == nil {
		t.Error("Expected quit command, got nil")
	}
}
