package memory

import (
	"sync"

	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
)

// Store implements storage.Store interface with in-memory storage
type Store struct {
	mu    sync.RWMutex
	todos map[string]*model.Todo
	notes []string
}

// New creates a new memory store
func New() *Store {
	return &Store{
		todos: make(map[string]*model.Todo),
		notes: make([]string, 0),
	}
}

// Add note-related methods to implement the full Store interface
func (s *Store) SaveNote(content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes = append(s.notes, content)
	return nil
}

func (s *Store) GetNotes() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	notes := make([]string, len(s.notes))
	copy(notes, s.notes)
	return notes, nil
}

// Add adds a new todo
func (s *Store) Add(todo *model.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.todos[todo.ID] = todo
	return nil
}

// Get retrieves a todo by ID
func (s *Store) Get(id string) (*model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if todo, exists := s.todos[id]; exists {
		return todo, nil
	}
	return nil, storage.ErrTodoNotFound
}

// List returns all todos
func (s *Store) List() []*model.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	todos := make([]*model.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos
}

// Update updates an existing todo
func (s *Store) Update(todo *model.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.todos[todo.ID]; !exists {
		return storage.ErrTodoNotFound
	}
	s.todos[todo.ID] = todo
	return nil
}

// Delete removes a todo
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.todos[id]; !exists {
		return storage.ErrTodoNotFound
	}
	delete(s.todos, id)
	return nil
}
