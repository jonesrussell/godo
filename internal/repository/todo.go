package repository

import (
    "context"
    "database/sql"
    "time"

    "github.com/jonesrussell/godo/internal/logger"
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
    logger.Debug("Creating new SQLite todo repository")
    return &SQLiteTodoRepository{db: db}
}

// Create inserts a new todo item
func (r *SQLiteTodoRepository) Create(ctx context.Context, todo *model.Todo) error {
    logger.Debug("Creating new todo: %s", todo.Title)
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
        logger.Error("Failed to create todo: %v", err)
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        logger.Error("Failed to get last insert ID: %v", err)
        return err
    }

    todo.ID = id
    logger.Debug("Successfully created todo with ID: %d", id)
    return nil
}

// GetByID retrieves a todo by its ID
func (r *SQLiteTodoRepository) GetByID(ctx context.Context, id int64) (*model.Todo, error) {
    logger.Debug("Getting todo by ID: %d", id)
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
        logger.Debug("No todo found with ID: %d", id)
        return nil, nil
    }
    if err != nil {
        logger.Error("Error getting todo by ID: %v", err)
        return nil, err
    }
    
    logger.Debug("Successfully retrieved todo: %s", todo.Title)
    return todo, nil
}

// List returns all todos
func (r *SQLiteTodoRepository) List(ctx context.Context) ([]model.Todo, error) {
    logger.Debug("Listing all todos")
    query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM todos
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        logger.Error("Error listing todos: %v", err)
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
            logger.Error("Error scanning todo row: %v", err)
            return nil, err
        }
        todos = append(todos, todo)
    }

    if err = rows.Err(); err != nil {
        logger.Error("Error after scanning todos: %v", err)
        return nil, err
    }

    logger.Debug("Successfully retrieved %d todos", len(todos))
    return todos, nil
}

// Update modifies an existing todo
func (r *SQLiteTodoRepository) Update(ctx context.Context, todo *model.Todo) error {
    logger.Debug("Updating todo ID: %d", todo.ID)
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
        logger.Error("Error updating todo: %v", err)
        return err
    }

    rows, err := result.RowsAffected()
    if err != nil {
        logger.Error("Error getting rows affected: %v", err)
        return err
    }
    if rows == 0 {
        logger.Debug("No todo found to update with ID: %d", todo.ID)
        return sql.ErrNoRows
    }

    logger.Debug("Successfully updated todo ID: %d", todo.ID)
    return nil
}

// Delete removes a todo by its ID
func (r *SQLiteTodoRepository) Delete(ctx context.Context, id int64) error {
    logger.Debug("Deleting todo ID: %d", id)
    query := `DELETE FROM todos WHERE id = ?`
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        logger.Error("Error deleting todo: %v", err)
        return err
    }

    rows, err := result.RowsAffected()
    if err != nil {
        logger.Error("Error getting rows affected: %v", err)
        return err
    }
    if rows == 0 {
        logger.Debug("No todo found to delete with ID: %d", id)
        return sql.ErrNoRows
    }

    logger.Debug("Successfully deleted todo ID: %d", id)
    return nil
}
