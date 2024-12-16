package repository

import (
	"context"
	"errors"

	"github.com/jonesrussell/godo/internal/model"
)

var ErrNotFound = errors.New("todo not found")

// TodoRepository defines the interface for todo storage operations
type TodoRepository interface {
	Create(ctx context.Context, todo *model.Todo) error
	GetByID(ctx context.Context, id int64) (*model.Todo, error)
	List(ctx context.Context) ([]model.Todo, error)
	Update(ctx context.Context, todo *model.Todo) error
	Delete(ctx context.Context, id int64) error
}

// DB interface for database operations
type DB interface {
	Create(*model.Todo) error
	GetByID(int64) (*model.Todo, error)
	List() ([]model.Todo, error)
	Update(*model.Todo) error
	Delete(int64) error
}

// Repository implements TodoRepository
type Repository struct {
	db DB
}

func NewTodoRepository(db DB) TodoRepository {
	return &Repository{db: db}
}

// Implementation of TodoRepository interface
func (r *Repository) Create(ctx context.Context, todo *model.Todo) error {
	return r.db.Create(todo)
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*model.Todo, error) {
	return r.db.GetByID(id)
}

func (r *Repository) List(ctx context.Context) ([]model.Todo, error) {
	return r.db.List()
}

func (r *Repository) Update(ctx context.Context, todo *model.Todo) error {
	return r.db.Update(todo)
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	return r.db.Delete(id)
}
