// internal/repository/todo_test.go
package repository_test

import (
	"context"
	"testing"

	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/repository"
	"github.com/stretchr/testify/assert"
)

type testDB struct {
	todos map[int64]*model.Todo
}

func newTestDB() *testDB {
	return &testDB{
		todos: make(map[int64]*model.Todo),
	}
}

func TestTodoRepository_Create(t *testing.T) {
	db := newTestDB()
	repo := repository.NewTodoRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		todo    *model.Todo
		wantErr bool
	}{
		{
			name: "valid todo",
			todo: &model.Todo{
				Title:       "Test Todo",
				Description: "Test Description",
			},
			wantErr: false,
		},
		// Add more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.todo)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotZero(t, tt.todo.ID)
		})
	}
}

// Implement testDB methods to satisfy repository.DB interface
func (db *testDB) Create(todo *model.Todo) error {
	todo.ID = int64(len(db.todos) + 1)
	db.todos[todo.ID] = todo
	return nil
}

func (db *testDB) GetByID(id int64) (*model.Todo, error) {
	if todo, exists := db.todos[id]; exists {
		return todo, nil
	}
	return nil, repository.ErrNotFound
}

func (db *testDB) List() ([]model.Todo, error) {
	todos := make([]model.Todo, 0, len(db.todos))
	for _, todo := range db.todos {
		todos = append(todos, *todo)
	}
	return todos, nil
}

func (db *testDB) Update(todo *model.Todo) error {
	if _, exists := db.todos[todo.ID]; !exists {
		return repository.ErrNotFound
	}
	db.todos[todo.ID] = todo
	return nil
}

func (db *testDB) Delete(id int64) error {
	if _, exists := db.todos[id]; !exists {
		return repository.ErrNotFound
	}
	delete(db.todos, id)
	return nil
}
