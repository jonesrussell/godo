package service

import (
	"context"

	"github.com/jonesrussell/godo/internal/model"
)

// TodoServicer defines the interface for todo operations
type TodoServicer interface {
	CreateTodo(ctx context.Context, title, description string) (*model.Todo, error)
	GetTodo(ctx context.Context, id int64) (*model.Todo, error)
	ListTodos(ctx context.Context) ([]model.Todo, error)
	UpdateTodo(ctx context.Context, todo *model.Todo) error
	ToggleTodoStatus(ctx context.Context, id int64) error
	DeleteTodo(ctx context.Context, id int64) error
}
