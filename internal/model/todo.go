package model

import (
	"time"

	"github.com/google/uuid"
)

// Todo represents a single todo item
type Todo struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewTodo creates a new Todo item
func NewTodo(content string) *Todo {
	now := time.Now()
	return &Todo{
		ID:        uuid.New().String(),
		Content:   content,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ToggleDone toggles the done status of the todo
func (t *Todo) ToggleDone() {
	t.Done = !t.Done
	t.UpdatedAt = time.Now()
}

// UpdateContent updates the content of the todo
func (t *Todo) UpdateContent(content string) {
	t.Content = content
	t.UpdatedAt = time.Now()
}
