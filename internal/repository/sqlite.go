package repository

import (
	"context"
	"database/sql"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
)

type SQLiteTodoRepository struct {
	db *sql.DB
}

func NewSQLiteTodoRepository(db *sql.DB) TodoRepository {
	return &SQLiteTodoRepository{db: db}
}

func (r *SQLiteTodoRepository) Create(ctx context.Context, todo *model.Todo) error {
	logger.Debug("Creating todo in repository",
		"title", todo.Title,
		"description", todo.Description)

	query := `INSERT INTO todos (title, description, completed, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.CreatedAt,
		todo.UpdatedAt,
	)
	if err != nil {
		logger.Error("Failed to create todo",
			"title", todo.Title,
			"error", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Failed to get last insert ID", "error", err)
		return err
	}

	todo.ID = id
	logger.Debug("Successfully created todo", "id", id)
	return nil
}

func (r *SQLiteTodoRepository) GetByID(ctx context.Context, id int64) (*model.Todo, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at
			  FROM todos WHERE id = ?`

	todo := &model.Todo{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *SQLiteTodoRepository) List(ctx context.Context) ([]model.Todo, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at
			  FROM todos ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, rows.Err()
}

func (r *SQLiteTodoRepository) Update(ctx context.Context, todo *model.Todo) error {
	query := `UPDATE todos 
			  SET title = ?, description = ?, completed = ?, updated_at = ?
			  WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.UpdatedAt,
		todo.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *SQLiteTodoRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
