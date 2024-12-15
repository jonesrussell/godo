package repository

import (
    "context"
    "database/sql"
    "time"

    "github.com/jonesrussell/godo/internal/model"
)

// TodoRepository defines the interface for todo storage operations
type TodoRepository interface {
    Create(ctx context.Context, todo *model.Todo) error
    GetByID(ctx context.Context, id int64) (*model.Todo, error)
    List(ctx context.Context) ([]model.Todo, error)
    Update(ctx context.Context, todo *model.Todo) error
    Delete(ctx context.Context, id int64) error
}

// SQLiteTodoRepository implements TodoRepository using SQLite
type SQLiteTodoRepository struct {
    db *sql.DB
}

// NewSQLiteTodoRepository creates a new SQLite repository instance
func NewSQLiteTodoRepository(db *sql.DB) *SQLiteTodoRepository {
    return &SQLiteTodoRepository{db: db}
}

// Create inserts a new todo item
func (r *SQLiteTodoRepository) Create(ctx context.Context, todo *model.Todo) error {
    query := `
        INSERT INTO todos (title, description, completed, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `
    now := time.Now()
    todo.CreatedAt = now
    todo.UpdatedAt = now

    result, err := r.db.ExecContext(ctx, query,
        todo.Title,
        todo.Description,
        todo.Completed,
        todo.CreatedAt,
        todo.UpdatedAt,
    )
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }

    todo.ID = id
    return nil
}

// GetByID retrieves a todo by its ID
func (r *SQLiteTodoRepository) GetByID(ctx context.Context, id int64) (*model.Todo, error) {
    query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM todos WHERE id = ?
    `
    
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
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    return todo, nil
}

// List returns all todos
func (r *SQLiteTodoRepository) List(ctx context.Context) ([]model.Todo, error) {
    query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM todos
        ORDER BY created_at DESC
    `
    
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

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return todos, nil
}

// Update modifies an existing todo
func (r *SQLiteTodoRepository) Update(ctx context.Context, todo *model.Todo) error {
    query := `
        UPDATE todos 
        SET title = ?, description = ?, completed = ?, updated_at = ?
        WHERE id = ?
    `
    
    todo.UpdatedAt = time.Now()
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
        return sql.ErrNoRows
    }

    return nil
}

// Delete removes a todo by its ID
func (r *SQLiteTodoRepository) Delete(ctx context.Context, id int64) error {
    query := `DELETE FROM todos WHERE id = ?`
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rows == 0 {
        return sql.ErrNoRows
    }

    return nil
}
