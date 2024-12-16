package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/service"
)

// MockTodoRepository implements repository.TodoRepository for testing
type MockTodoRepository struct {
	Todos     map[int64]*model.Todo
	NextID    int64
	LastError error
}

func NewMockTodoRepository() *MockTodoRepository {
	return &MockTodoRepository{
		Todos:  make(map[int64]*model.Todo),
		NextID: 1,
	}
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *model.Todo) error {
	if m.LastError != nil {
		return m.LastError
	}
	todo.ID = m.NextID
	m.Todos[todo.ID] = todo
	m.NextID++
	return nil
}

func (m *MockTodoRepository) GetByID(ctx context.Context, id int64) (*model.Todo, error) {
	if m.LastError != nil {
		return nil, m.LastError
	}
	if todo, exists := m.Todos[id]; exists {
		return todo, nil
	}
	return nil, nil
}

func (m *MockTodoRepository) List(ctx context.Context) ([]model.Todo, error) {
	if m.LastError != nil {
		return nil, m.LastError
	}
	todos := make([]model.Todo, 0, len(m.Todos))
	for _, todo := range m.Todos {
		todos = append(todos, *todo)
	}
	return todos, nil
}

func (m *MockTodoRepository) Update(ctx context.Context, todo *model.Todo) error {
	if m.LastError != nil {
		return m.LastError
	}
	if _, exists := m.Todos[todo.ID]; !exists {
		return service.ErrNotFound
	}
	m.Todos[todo.ID] = todo
	return nil
}

func (m *MockTodoRepository) Delete(ctx context.Context, id int64) error {
	if m.LastError != nil {
		return m.LastError
	}
	if _, exists := m.Todos[id]; !exists {
		return service.ErrNotFound
	}
	delete(m.Todos, id)
	return nil
}

// MockTodoService implements TodoServicer for testing
type MockTodoService struct {
	todos       map[int64]*model.Todo
	nextID      int64
	createCalls int
	lastTitle   string
	shouldErr   bool
}

func NewMockTodoService() *MockTodoService {
	return &MockTodoService{
		todos:  make(map[int64]*model.Todo),
		nextID: 1,
	}
}

func (m *MockTodoService) CreateTodo(ctx context.Context, title, description string) (*model.Todo, error) {
	m.createCalls++
	m.lastTitle = title

	if m.shouldErr {
		return nil, service.ErrEmptyTitle
	}

	todo := &model.Todo{
		ID:          m.nextID,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.todos[todo.ID] = todo
	m.nextID++
	return todo, nil
}

// ... implement other interface methods ...

// Helper methods for testing
func (m *MockTodoService) SetError(shouldErr bool) {
	m.shouldErr = shouldErr
}

func (m *MockTodoService) GetCreateCalls() int {
	return m.createCalls
}

func (m *MockTodoService) GetLastTitle() string {
	return m.lastTitle
}

func TestCreateTodo(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	svc := service.NewTodoService(mockRepo)
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
			todo, err := svc.CreateTodo(ctx, tt.title, tt.description)
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
	mockRepo := NewMockTodoRepository()
	svc := service.NewTodoService(mockRepo)
	ctx := context.Background()

	// Create a test todo
	testTodo := &model.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := mockRepo.Create(ctx, testTodo); err != nil {
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
			todo, err := svc.GetTodo(ctx, tt.id)
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
