package testutil

import (
	"github.com/jonesrussell/godo/internal/model"
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

// ... rest of repository mock implementation ...

// Helper methods
func (m *MockTodoRepository) SetError(err error) {
	m.LastError = err
}

func (m *MockTodoRepository) Reset() {
	m.Todos = make(map[int64]*model.Todo)
	m.NextID = 1
	m.LastError = nil
}
