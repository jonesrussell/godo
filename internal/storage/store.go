package storage

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
)

// Task represents a todo item
type Task struct {
	ID        string
	Title     string
	Completed bool
}

// Store defines the interface for task storage
type Store interface {
	Add(task Task) error
	List() ([]Task, error)
	Update(task Task) error
	Delete(id string) error
}
