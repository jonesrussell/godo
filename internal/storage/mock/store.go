package mock

import "github.com/jonesrussell/godo/internal/model"

// Store is a mock implementation of storage.Store for testing
type Store struct {
	notes []string
	todos []*model.Todo
}

// NewStore creates a new mock store
func NewStore() *Store {
	return &Store{
		notes: make([]string, 0),
		todos: make([]*model.Todo, 0),
	}
}

// Add adds a todo to the store
func (s *Store) Add(todo *model.Todo) error {
	s.todos = append(s.todos, todo)
	return nil
}

// Update updates a todo in the store
func (s *Store) Update(todo *model.Todo) error {
	for i, t := range s.todos {
		if t.ID == todo.ID {
			s.todos[i] = todo
			return nil
		}
	}
	return nil
}

// Delete removes a todo from the store
func (s *Store) Delete(id string) error {
	for i, todo := range s.todos {
		if todo.ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return nil
		}
	}
	return nil
}

// Get retrieves a todo by ID
func (s *Store) Get(id string) (*model.Todo, error) {
	for _, todo := range s.todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return nil, nil
}

// List returns all todos
func (s *Store) List() []*model.Todo {
	return s.todos
}

// SaveNote adds a note to the store
func (s *Store) SaveNote(content string) error {
	s.notes = append(s.notes, content)
	return nil
}

// GetNotes returns all notes in the store
func (s *Store) GetNotes() ([]string, error) {
	return s.notes, nil
}

// Close is a no-op for the mock store
func (s *Store) Close() error {
	return nil
}
