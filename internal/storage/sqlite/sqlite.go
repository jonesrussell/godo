package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

const (
	SchemaVersion = 1
	Schema        = `
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		content TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
)

type Store struct {
	db     *sql.DB
	logger logger.Logger
}

// validatePath checks if the path contains invalid characters
func validatePath(dbPath string) error {
	// Allow in-memory database paths
	if strings.HasPrefix(dbPath, "file::memory:") {
		return nil
	}

	// Invalid characters in Windows paths, excluding the drive letter colon
	invalidChars := []string{"<", ">", "\"", "|", "?", "*"}

	// For Windows paths, allow colon only after drive letter
	if strings.Count(dbPath, ":") > 1 || (strings.Contains(dbPath, ":") && len(dbPath) < 2) {
		return fmt.Errorf("invalid path: multiple colons or invalid drive format")
	}

	for _, char := range invalidChars {
		if strings.Contains(dbPath, char) {
			return fmt.Errorf("path contains invalid character: %s", char)
		}
	}
	return nil
}

func New(dbPath string, log logger.Logger) (*Store, error) {
	if log == nil {
		return nil, fmt.Errorf("logger is required")
	}

	log.Info("Opening database", "path", dbPath)

	// Validate the path first
	if err := validatePath(dbPath); err != nil {
		log.Error("Invalid database path", "error", err)
		return nil, err
	}

	// Only create directories for non-memory databases
	if !strings.HasPrefix(dbPath, "file::memory:") {
		// Try to create the database directory
		if err := ensureDataDir(dbPath); err != nil {
			log.Error("Failed to create database directory", "error", err)
			return nil, err
		}

		// Verify we can write to the database path
		if err := verifyDatabaseAccess(dbPath); err != nil {
			log.Error("Failed to verify database access", "error", err)
			return nil, err
		}
	}

	// Open the database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Error("Failed to open database", "error", err)
		return nil, err
	}

	store := &Store{
		db:     db,
		logger: log,
	}

	// Try to initialize the database
	if err := store.initialize(); err != nil {
		// Close the database connection before returning error
		if closeErr := db.Close(); closeErr != nil {
			log.Error("Failed to close database after initialization error", "error", closeErr)
		}
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

func (s *Store) initSchema() error {
	// Create schema_version table
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		s.logger.Error("Failed to create schema_version table", "error", err)
		return err
	}

	// Check current version
	var version int
	err = s.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
	if err != nil {
		s.logger.Error("Failed to get schema version", "error", err)
		return err
	}

	s.logger.Debug("Schema versions", "current", version, "target", SchemaVersion)

	if version < SchemaVersion {
		tx, err := s.db.Begin()
		if err != nil {
			s.logger.Error("Failed to begin transaction", "error", err)
			return err
		}

		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				s.logger.Error("Failed to rollback transaction", "error", err)
			}
		}()

		// Create or update tables
		if _, err := tx.Exec(Schema); err != nil {
			s.logger.Error("Failed to create todos table", "error", err)
			return err
		}

		// Update schema version
		if _, err := tx.Exec("INSERT INTO schema_version (version) VALUES (?)", SchemaVersion); err != nil {
			s.logger.Error("Failed to update schema version", "error", err)
			return err
		}

		if err := tx.Commit(); err != nil {
			s.logger.Error("Failed to commit schema changes", "error", err)
			return err
		}

		s.logger.Info("Schema updated", "version", SchemaVersion)
	}

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
	// Clean and normalize the path
	dbPath = filepath.Clean(dbPath)
	dir := filepath.Dir(dbPath)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0o755); err != nil {
		if !os.IsPermission(err) && !os.IsExist(err) {
			return err
		}
	}
	return nil
}

// verifyDatabaseAccess checks if we can write to the database file
func verifyDatabaseAccess(dbPath string) error {
	// Try to create/open the file
	f, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	return f.Close()
}

func (s *Store) SaveNote(content string) error {
	query := `INSERT INTO notes (content) VALUES (?)`
	_, err := s.db.Exec(query, content)
	if err != nil {
		s.logger.Error("Failed to save note", "error", err)
		return err
	}
	s.logger.Info("Note saved successfully")
	return nil
}

func (s *Store) GetNotes() ([]string, error) {
	query := `SELECT content FROM notes ORDER BY created_at DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Error("Failed to get notes", "error", err)
		return nil, err
	}
	defer rows.Close()

	var notes []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			s.logger.Error("Failed to scan note", "error", err)
			continue
		}
		notes = append(notes, content)
	}
	return notes, nil
}
