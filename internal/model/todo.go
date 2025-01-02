// Package model defines the core domain models for the application
package model

import (
	"time"

	"github.com/google/uuid"
)

// Todo represents a todo item
type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// NewTodo creates a new todo item
func NewTodo(title string) *Todo {
	now := time.Now().Unix()
	return &Todo{
		ID:        uuid.New().String(),
		Title:     title,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ToggleCompleted toggles the completed status of the todo
func (t *Todo) ToggleCompleted() {
	t.Completed = !t.Completed
	t.UpdatedAt = time.Now().Unix()
}

// UpdateTitle updates the title of the todo
func (t *Todo) UpdateTitle(title string) {
	t.Title = title
	t.UpdatedAt = time.Now().Unix()
}
