package sqlite

import (
	"database/sql"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

// Store implements todo storage using SQLite
type Store struct {
	db *sql.DB
}

// New creates a new SQLite store
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &Store{db: db}
	if err := store.migrate(); err != nil {
		return nil, err
	}

	return store, nil
}

// migrate creates the necessary database tables
func (s *Store) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		content TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		logger.Error("Failed to create todos table", "error", err)
		return err
	}

	logger.Info("Database migration completed")
	return nil
}

// Add adds a new todo to storage
func (s *Store) Add(todo *model.Todo) error {
	query := `
	INSERT INTO todos (id, content, done, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, todo.ID, todo.Content, todo.Done, todo.CreatedAt, todo.UpdatedAt)
	if err != nil {
		logger.Error("Failed to add todo", "error", err)
		return err
	}

	return nil
}

// Get retrieves a todo by ID
func (s *Store) Get(id string) (*model.Todo, error) {
	query := `SELECT id, content, done, created_at, updated_at FROM todos WHERE id = ?`

	todo := &model.Todo{}
	err := s.db.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Content,
		&todo.Done,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Debug("Todo not found", "id", id)
		return nil, storage.ErrTodoNotFound
	}
	if err != nil {
		logger.Error("Failed to get todo", "error", err)
		return nil, err
	}

	return todo, nil
}

// List returns all todos
func (s *Store) List() []*model.Todo {
	query := `SELECT id, content, done, created_at, updated_at FROM todos ORDER BY created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		logger.Error("Failed to list todos", "error", err)
		return nil
	}
	defer rows.Close()

	var todos []*model.Todo
	for rows.Next() {
		todo := &model.Todo{}
		err := rows.Scan(
			&todo.ID,
			&todo.Content,
			&todo.Done,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			logger.Error("Failed to scan todo row", "error", err)
			continue
		}
		todos = append(todos, todo)
	}

	return todos
}

// Delete removes a todo by ID
func (s *Store) Delete(id string) error {
	query := `DELETE FROM todos WHERE id = ?`

	result, err := s.db.Exec(query, id)
	if err != nil {
		logger.Error("Failed to delete todo", "error", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("Failed to get rows affected", "error", err)
		return err
	}

	if rows == 0 {
		logger.Debug("Todo not found for deletion", "id", id)
		return storage.ErrTodoNotFound
	}

	logger.Debug("Deleted todo", "id", id)
	return nil
}

// Close closes the database connection
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Update updates an existing todo
func (s *Store) Update(todo *model.Todo) error {
	query := `
	UPDATE todos 
	SET content = ?, done = ?, updated_at = datetime('now')
	WHERE id = ?`

	result, err := s.db.Exec(query, todo.Content, todo.Done, todo.ID)
	if err != nil {
		logger.Error("Failed to update todo", "error", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("Failed to get rows affected", "error", err)
		return err
	}

	if rows == 0 {
		logger.Debug("Todo not found for update", "id", todo.ID)
		return storage.ErrTodoNotFound
	}

	logger.Debug("Updated todo", "id", todo.ID)
	return nil
}
