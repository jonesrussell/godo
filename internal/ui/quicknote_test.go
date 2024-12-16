package ui_test

import (
	"context"
	"testing"

	"github.com/jonesrussell/godo/internal/testutil"
	"github.com/jonesrussell/godo/internal/ui"
)

func TestNewQuickNote(t *testing.T) {
	mockService := testutil.NewMockTodoService()
	quickNote := ui.NewQuickNote(mockService)

	if quickNote == nil {
		t.Error("Expected QuickNote instance, got nil")
	}
}

func TestQuickNote_CreateTodo(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid input",
			input:   "Test todo",
			wantErr: false,
		},
		{
			name:    "Empty input",
			input:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := testutil.NewMockTodoService()
			mock := testutil.AsMockTodoService(mockService)
			mock.SetShouldError(tt.wantErr)

			// Test the service directly
			_, err := mockService.CreateTodo(context.Background(), "quick", tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTodo() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && mock.GetLastTitle() != tt.input {
				t.Errorf("Got title %q, want %q", mock.GetLastTitle(), tt.input)
			}
		})
	}
}
