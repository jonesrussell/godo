package service

import (
	"context"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/repository"
)

// MockTodoRepository implements TodoRepository interface for testing
type MockTodoRepository struct {
	todos  map[int64]*model.Todo
	nextID int64
}

func NewMockTodoRepository() *MockTodoRepository {
	return &MockTodoRepository{
		todos:  make(map[int64]*model.Todo),
		nextID: 1,
	}
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *model.Todo) error {
	todo.ID = m.nextID
	m.todos[todo.ID] = todo
	m.nextID++
	return nil
}

func (m *MockTodoRepository) GetByID(ctx context.Context, id int64) (*model.Todo, error) {
	if todo, exists := m.todos[id]; exists {
		return todo, nil
	}
	return nil, nil
}

func (m *MockTodoRepository) List(ctx context.Context) ([]model.Todo, error) {
	todos := make([]model.Todo, 0, len(m.todos))
	for _, todo := range m.todos {
		todos = append(todos, *todo)
	}
	return todos, nil
}

func (m *MockTodoRepository) Update(ctx context.Context, todo *model.Todo) error {
	if _, exists := m.todos[todo.ID]; !exists {
		return ErrNotFound
	}
	m.todos[todo.ID] = todo
	return nil
}

func (m *MockTodoRepository) Delete(ctx context.Context, id int64) error {
	if _, exists := m.todos[id]; !exists {
		return ErrNotFound
	}
	delete(m.todos, id)
	return nil
}

func TestCreateTodo(t *testing.T) {
	var repo repository.TodoRepository = NewMockTodoRepository()
	service := NewTodoService(repo)
	ctx := context.Background()

	tests := []struct {
		name        string
		title       string
		description string
		wantErr     bool
	}{
		{
			name:        "Valid todo",
			title:       "Test Todo",
			description: "Test Description",
			wantErr:     false,
		},
		{
			name:        "Empty title",
			title:       "",
			description: "Test Description",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := service.CreateTodo(ctx, tt.title, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && todo == nil {
				t.Error("CreateTodo() returned nil todo when no error expected")
			}
			if !tt.wantErr {
				if todo.Title != tt.title {
					t.Errorf("CreateTodo() title = %v, want %v", todo.Title, tt.title)
				}
				if todo.Description != tt.description {
					t.Errorf("CreateTodo() description = %v, want %v", todo.Description, tt.description)
				}
			}
		})
	}
}

func TestGetTodo(t *testing.T) {
	var repo repository.TodoRepository = NewMockTodoRepository()
	service := NewTodoService(repo)
	ctx := context.Background()

	// Create a test todo
	testTodo := &model.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := repo.Create(ctx, testTodo); err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "Existing todo",
			id:      1,
			wantErr: false,
		},
		{
			name:    "Non-existing todo",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := service.GetTodo(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && todo == nil {
				t.Error("GetTodo() returned nil todo when no error expected")
			}
		})
	}
}
