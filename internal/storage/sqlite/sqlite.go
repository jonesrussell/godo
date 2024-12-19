package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db     *sql.DB
	logger logger.Logger
}

func New(dbPath string, log logger.Logger) (*Store, error) {
	log.Info("Opening database", "path", dbPath)

	if err := ensureDataDir(dbPath); err != nil {
		log.Error("Failed to create database directory", "error", err)
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Error("Failed to open database", "error", err)
		return nil, err
	}

	store := &Store{
		db:     db,
		logger: log,
	}

	if err := store.initialize(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) initialize() error {
	s.db.SetMaxOpenConns(1)
	s.db.SetMaxIdleConns(1)

	if _, err := s.db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		s.logger.Error("Failed to enable foreign keys", "error", err)
		return err
	}

	if err := s.db.Ping(); err != nil {
		s.logger.Error("Database ping failed", "error", err)
		return err
	}
	s.logger.Info("Database connection successful")

	s.logger.Info("Initializing database schema...")
	if err := s.initSchema(); err != nil {
		s.logger.Error("Schema initialization failed", "error", err)
		return err
	}
	s.logger.Info("Schema initialized successfully")

	return nil
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
		s.logger.Error("Failed to create todos table", "error", err)
		return err
	}

	s.logger.Info("Database migration completed")
	return nil
}

// Add adds a new todo to storage
func (s *Store) Add(todo *model.Todo) error {
	query := `
	INSERT INTO todos (id, content, done, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, todo.ID, todo.Content, todo.Done, todo.CreatedAt, todo.UpdatedAt)
	if err != nil {
		s.logger.Error("Failed to add todo", "error", err)
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
		s.logger.Debug("Todo not found", "id", id)
		return nil, storage.ErrTodoNotFound
	}
	if err != nil {
		s.logger.Error("Failed to get todo", "error", err)
		return nil, err
	}

	return todo, nil
}

// List returns all todos
func (s *Store) List() []*model.Todo {
	query := `SELECT id, content, done, created_at, updated_at FROM todos ORDER BY created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Error("Failed to list todos", "error", err)
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
			s.logger.Error("Failed to scan todo row", "error", err)
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
		s.logger.Error("Failed to delete todo", "error", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("Failed to get rows affected", "error", err)
		return err
	}

	if rows == 0 {
		s.logger.Debug("Todo not found for deletion", "id", id)
		return storage.ErrTodoNotFound
	}

	s.logger.Debug("Deleted todo", "id", id)
	return nil
}

// Close closes the database connection
func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		s.logger.Error("Error closing database", "error", err)
		return err
	}
	s.logger.Info("Database closed successfully")
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
		s.logger.Error("Failed to update todo", "error", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("Failed to get rows affected", "error", err)
		return err
	}

	if rows == 0 {
		s.logger.Debug("Todo not found for update", "id", todo.ID)
		return storage.ErrTodoNotFound
	}

	s.logger.Debug("Updated todo", "id", todo.ID)
	return nil
}

// ensureDataDir creates the database directory if it doesn't exist
func ensureDataDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	return os.MkdirAll(dir, 0755)
}

// Add this method to the Store struct
func (s *Store) initSchema() error {
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
		s.logger.Error("Failed to initialize schema", "error", err)
		return err
	}

	s.logger.Info("Schema initialized successfully")
	return nil
}
