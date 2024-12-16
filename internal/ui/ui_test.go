package ui

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/service"
)

// MockTodoRepository implements repository.TodoRepository for testing
type MockTodoRepository struct {
	todos     []model.Todo
	lastError error
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *model.Todo) error {
	if m.lastError != nil {
		return m.lastError
	}
	todo.ID = int64(len(m.todos) + 1)
	m.todos = append(m.todos, *todo)
	return nil
}

func (m *MockTodoRepository) GetByID(ctx context.Context, id int64) (*model.Todo, error) {
	if m.lastError != nil {
		return nil, m.lastError
	}
	for _, todo := range m.todos {
		if todo.ID == id {
			return &todo, nil
		}
	}
	return nil, nil
}

func (m *MockTodoRepository) List(ctx context.Context) ([]model.Todo, error) {
	if m.lastError != nil {
		return nil, m.lastError
	}
	return m.todos, nil
}

func (m *MockTodoRepository) Delete(ctx context.Context, id int64) error {
	if m.lastError != nil {
		return m.lastError
	}
	for i, todo := range m.todos {
		if todo.ID == id {
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			return nil
		}
	}
	return service.ErrNotFound
}

func (m *MockTodoRepository) Update(ctx context.Context, todo *model.Todo) error {
	if m.lastError != nil {
		return m.lastError
	}
	for i, t := range m.todos {
		if t.ID == todo.ID {
			m.todos[i] = *todo
			return nil
		}
	}
	return service.ErrNotFound
}

func NewMockTodoService() *service.TodoService {
	mockRepo := &MockTodoRepository{
		todos: make([]model.Todo, 0),
	}
	return service.NewTodoService(mockRepo)
}

func TestTodoUI_Init(t *testing.T) {
	mockService := NewMockTodoService()
	ui := New(mockService)
	cmd := ui.Init()
	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

func TestTodoUI_Update(t *testing.T) {
	tests := []struct {
		name     string
		msg      tea.Msg
		wantCmd  bool
		wantQuit bool
	}{
		{
			name:     "Quit message",
			msg:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantCmd:  true,
			wantQuit: true,
		},
		{
			name:    "Add message",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			wantCmd: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockTodoService()
			ui := New(mockService)
			model, cmd := ui.Update(tt.msg)

			if (cmd != nil) != tt.wantCmd {
				t.Errorf("Update() cmd = %v, want %v", cmd != nil, tt.wantCmd)
			}

			if tt.wantQuit {
				if cmd == nil {
					t.Error("Update() should return a command when quitting")
				}
			}

			if model == nil {
				t.Error("Update() should return a model")
			}
		})
	}
}

func TestTodoUI_View(t *testing.T) {
	mockService := NewMockTodoService()
	ui := New(mockService)
	ui.todos = []model.Todo{
		{ID: 1, Title: "Test Todo", Completed: false},
		{ID: 2, Title: "Completed Todo", Completed: true},
	}

	view := ui.View()
	if view == "" {
		t.Error("View() should return non-empty string")
	}

	// Test adding mode
	ui.adding = true
	addView := ui.View()
	if addView == "" {
		t.Error("View() in adding mode should return non-empty string")
	}
}

func TestTodoUI_renderTodoItem(t *testing.T) {
	ui := &TodoUI{cursor: 0}
	todo := model.Todo{
		ID:        1,
		Title:     "Test Todo",
		Completed: false,
	}

	// Test uncompleted todo
	result := ui.renderTodoItem(0, todo)
	expected := "  > ☐ Test Todo\n"
	if result != expected {
		t.Errorf("renderTodoItem() = %q, want %q", result, expected)
	}

	// Test completed todo
	todo.Completed = true
	result = ui.renderTodoItem(0, todo)
	expected = "  > ☑ Test Todo\n"
	if result != expected {
		t.Errorf("renderTodoItem() = %q, want %q", result, expected)
	}
}

func TestTodoUI_getSelectedTodoID(t *testing.T) {
	tests := []struct {
		name    string
		todos   []model.Todo
		cursor  int
		wantID  int64
		wantErr bool
	}{
		{
			name:    "Empty todos",
			todos:   []model.Todo{},
			cursor:  0,
			wantID:  0,
			wantErr: true,
		},
		{
			name: "Valid selection",
			todos: []model.Todo{
				{ID: 1, Title: "Test Todo"},
			},
			cursor:  0,
			wantID:  1,
			wantErr: false,
		},
		{
			name: "Cursor out of bounds",
			todos: []model.Todo{
				{ID: 1, Title: "Test Todo"},
			},
			cursor:  1,
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &TodoUI{
				todos:  tt.todos,
				cursor: tt.cursor,
			}

			gotID, err := ui.getSelectedTodoID()
			if (err != nil) != tt.wantErr {
				t.Errorf("getSelectedTodoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("getSelectedTodoID() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}
