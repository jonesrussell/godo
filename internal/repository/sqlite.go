package repository

import (
	"database/sql"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteTodoRepository(db *sql.DB) TodoRepository {
	return NewTodoRepository(&SQLiteDB{db: db})
}

func (r *SQLiteDB) Create(todo *model.Todo) error {
	logger.Debug("Creating todo: %+v", todo)
	query := `
        INSERT INTO todos (title, description, completed, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `
	now := time.Now()
	result, err := r.db.Exec(query,
		todo.Title,
		todo.Description,
		todo.Completed,
		now,
		now,
	)
	if err != nil {
		logger.Error("Failed to create todo: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Failed to get last insert ID: %v", err)
		return err
	}

	todo.ID = id
	todo.CreatedAt = now
	todo.UpdatedAt = now
	logger.Debug("Successfully created todo with ID: %d", id)
	return nil
}

func (r *SQLiteDB) GetByID(id int64) (*model.Todo, error) {
	query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM todos
        WHERE id = ?
    `
	todo := &model.Todo{}
	err := r.db.QueryRow(query, id).Scan(
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

func (r *SQLiteDB) List() ([]model.Todo, error) {
	logger.Debug("Listing todos")
	query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM todos
        ORDER BY created_at DESC
    `
	rows, err := r.db.Query(query)
	if err != nil {
		logger.Error("Failed to query todos: %v", err)
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
			logger.Error("Failed to scan todo: %v", err)
			return nil, err
		}
		todos = append(todos, todo)
	}
	logger.Debug("Found %d todos", len(todos))
	return todos, rows.Err()
}

func (r *SQLiteDB) Update(todo *model.Todo) error {
	query := `
        UPDATE todos
        SET title = ?, description = ?, completed = ?, updated_at = ?
        WHERE id = ?
    `
	now := time.Now()
	result, err := r.db.Exec(query,
		todo.Title,
		todo.Description,
		todo.Completed,
		now,
		todo.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	todo.UpdatedAt = now
	return nil
}

func (r *SQLiteDB) Delete(id int64) error {
	query := `DELETE FROM todos WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
