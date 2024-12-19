package storage

import "github.com/jonesrussell/godo/internal/model"

// Store defines the interface for todo storage
type Store interface {
	Add(todo *model.Todo) error
	Get(id string) (*model.Todo, error)
	List() []*model.Todo
	Update(todo *model.Todo) error
	Delete(id string) error
}
