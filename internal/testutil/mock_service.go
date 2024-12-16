package testutil

import (
	"context"

	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/service"
)

// MockTodoService implements service.TodoServicer for testing
type MockTodoService struct {
	todos       map[int64]*model.Todo
	nextID      int64
	lastTitle   string
	shouldError bool
}

func NewMockTodoService() service.TodoServicer {
	return &MockTodoService{
		todos:  make(map[int64]*model.Todo),
		nextID: 1,
	}
}

func (m *MockTodoService) CreateTodo(ctx context.Context, title, description string) (*model.Todo, error) {
	m.lastTitle = title
	if m.shouldError {
		return nil, service.ErrEmptyTitle
	}
	todo := &model.Todo{
		ID:          m.nextID,
		Title:       title,
		Description: description,
	}
	m.todos[todo.ID] = todo
	m.nextID++
	return todo, nil
}

func (m *MockTodoService) GetTodo(ctx context.Context, id int64) (*model.Todo, error) {
	if todo, exists := m.todos[id]; exists {
		return todo, nil
	}
	return nil, service.ErrNotFound
}

func (m *MockTodoService) ListTodos(ctx context.Context) ([]model.Todo, error) {
	todos := make([]model.Todo, 0, len(m.todos))
	for _, todo := range m.todos {
		todos = append(todos, *todo)
	}
	return todos, nil
}

func (m *MockTodoService) UpdateTodo(ctx context.Context, todo *model.Todo) error {
	if _, exists := m.todos[todo.ID]; !exists {
		return service.ErrNotFound
	}
	m.todos[todo.ID] = todo
	return nil
}

func (m *MockTodoService) ToggleTodoStatus(ctx context.Context, id int64) error {
	todo, exists := m.todos[id]
	if !exists {
		return service.ErrNotFound
	}
	todo.Completed = !todo.Completed
	return nil
}

func (m *MockTodoService) DeleteTodo(ctx context.Context, id int64) error {
	if _, exists := m.todos[id]; !exists {
		return service.ErrNotFound
	}
	delete(m.todos, id)
	return nil
}

// Test helper methods
func (m *MockTodoService) GetLastTitle() string {
	return m.lastTitle
}

func (m *MockTodoService) SetShouldError(should bool) {
	m.shouldError = should
}

// Helper method to access the mock service for assertions
func AsMockTodoService(s service.TodoServicer) *MockTodoService {
	return s.(*MockTodoService)
}
