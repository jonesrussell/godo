package ui_test

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/godo/internal/testutil"
	"github.com/jonesrussell/godo/internal/ui"
)

func TestTodoUI(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(_ *testing.T, _ *ui.TodoUI, _ *testutil.MockTodoService)
		msg         tea.Msg
		wantQuit    bool
		wantCommand bool
	}{
		{
			name:        "Quit on q",
			setup:       func(_ *testing.T, _ *ui.TodoUI, _ *testutil.MockTodoService) {},
			msg:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantQuit:    true,
			wantCommand: true,
		},
		{
			name:        "Add mode on a",
			setup:       func(_ *testing.T, _ *ui.TodoUI, _ *testutil.MockTodoService) {},
			msg:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			wantQuit:    false,
			wantCommand: false,
		},
		{
			name: "Toggle todo on space",
			setup: func(t *testing.T, ui *ui.TodoUI, mock *testutil.MockTodoService) {
				// Create a test todo
				_, err := mock.CreateTodo(context.Background(), "Test Todo", "")
				if err != nil {
					t.Fatalf("Failed to create test todo: %v", err)
				}

				// Execute the load command
				initCmd := ui.Init()
				if initCmd == nil {
					t.Fatal("Init() returned nil command")
				}

				msg := initCmd()
				model, _ := ui.Update(msg)
				if model == nil {
					t.Fatal("Update returned nil model")
				}
			},
			msg:         tea.KeyMsg{Type: tea.KeySpace},
			wantQuit:    false,
			wantCommand: true,
		},
		{
			name: "Delete todo on d",
			setup: func(t *testing.T, ui *ui.TodoUI, mock *testutil.MockTodoService) {
				// Create a test todo
				_, err := mock.CreateTodo(context.Background(), "Test Todo", "")
				if err != nil {
					t.Fatalf("Failed to create test todo: %v", err)
				}

				// Execute the load command to update UI state
				initCmd := ui.Init()
				if initCmd == nil {
					t.Fatal("Init() returned nil command")
				}

				msg := initCmd()
				model, _ := ui.Update(msg)
				if model == nil {
					t.Fatal("Update returned nil model")
				}
			},
			msg:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			wantQuit:    false,
			wantCommand: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh instances for each test
			mockService := testutil.NewMockTodoService()
			todoUI := ui.New(mockService)
			mock := testutil.AsMockTodoService(mockService)

			// Initialize UI
			initCmd := todoUI.Init()
			if initCmd == nil {
				t.Fatal("Init() returned nil command")
			}

			loadMsg := initCmd()
			model, loadCmd := todoUI.Update(loadMsg)
			if model == nil {
				t.Fatal("Update after Init returned nil model")
			}
			todoUI = model.(*ui.TodoUI)

			if loadCmd != nil {
				msg := loadCmd()
				model, _ = todoUI.Update(msg)
				if model != nil {
					todoUI = model.(*ui.TodoUI)
				}
			}

			// Run test-specific setup
			tt.setup(t, todoUI, mock)

			// Run the test
			model, cmd := todoUI.Update(tt.msg)

			if cmd == nil && tt.wantCommand {
				t.Errorf("Expected a command but got nil")
			}

			if cmd != nil {
				if tt.wantQuit {
					msg := cmd()
					if _, ok := msg.(tea.QuitMsg); !ok {
						t.Error("Expected quit command")
					}
				}
				if !tt.wantCommand {
					t.Error("Got unexpected command")
				}
			}

			if model == nil {
				t.Error("Update() returned nil model")
			}
		})
	}
}

func TestTodoUI_Navigation(t *testing.T) {
	mockService := testutil.NewMockTodoService()
	todoUI := ui.New(mockService)

	// Create some test todos
	mock := testutil.AsMockTodoService(mockService)
	mock.CreateTodo(context.Background(), "Test 1", "")
	mock.CreateTodo(context.Background(), "Test 2", "")

	// Test navigation keys
	navTests := []struct {
		name     string
		keyType  tea.KeyType
		keyRunes []rune
	}{
		{
			name:     "Move down with down arrow",
			keyType:  tea.KeyDown,
			keyRunes: nil,
		},
		{
			name:     "Move up with up arrow",
			keyType:  tea.KeyUp,
			keyRunes: nil,
		},
		{
			name:     "Move down with j",
			keyType:  tea.KeyRunes,
			keyRunes: []rune{'j'},
		},
		{
			name:     "Move up with k",
			keyType:  tea.KeyRunes,
			keyRunes: []rune{'k'},
		},
	}

	for _, tt := range navTests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{
				Type:  tt.keyType,
				Runes: tt.keyRunes,
			}

			model, _ := todoUI.Update(msg)
			if model == nil {
				t.Error("Update() returned nil model")
			}
		})
	}
}
