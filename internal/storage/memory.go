package storage

import (
	"sync"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
)

// MemoryStore implements in-memory todo storage
type MemoryStore struct {
	todos map[string]*model.Todo
	mu    sync.RWMutex
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		todos: make(map[string]*model.Todo),
	}
}

// Add adds a new todo to storage
func (s *MemoryStore) Add(todo *model.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.todos[todo.ID] = todo
	logger.Debug("Added todo", "id", todo.ID, "content", todo.Content)
	return nil
}

// Get retrieves a todo by ID
func (s *MemoryStore) Get(id string) (*model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, exists := s.todos[id]
	if !exists {
		logger.Debug("Todo not found", "id", id)
		return nil, ErrTodoNotFound
	}

	return todo, nil
}

// List returns all todos
func (s *MemoryStore) List() []*model.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]*model.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}

	logger.Debug("Listed todos", "count", len(todos))
	return todos
}

// Update updates an existing todo
func (s *MemoryStore) Update(todo *model.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.todos[todo.ID]; !exists {
		logger.Debug("Todo not found for update", "id", todo.ID)
		return ErrTodoNotFound
	}

	s.todos[todo.ID] = todo
	logger.Debug("Updated todo", "id", todo.ID)
	return nil
}

// Delete removes a todo by ID
func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.todos[id]; !exists {
		logger.Debug("Todo not found for deletion", "id", id)
		return ErrTodoNotFound
	}

	delete(s.todos, id)
	logger.Debug("Deleted todo", "id", id)
	return nil
}
