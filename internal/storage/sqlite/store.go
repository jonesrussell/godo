package sqlite

import (
	"database/sql"

	"github.com/jonesrussell/godo/internal/storage"
	"go.uber.org/zap"
	_ "modernc.org/sqlite" // SQLite driver
)

// Store implements storage.Store using SQLite
type Store struct {
	db     *sql.DB
	logger *zap.Logger
}

// New creates a new SQLite store
func New(path string, logger *zap.Logger) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	store := &Store{
		db:     db,
		logger: logger,
	}

	if err := RunMigrations(db); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Add(task storage.Task) error {
	_, err := s.db.Exec(
		"INSERT INTO tasks (id, title, completed) VALUES (?, ?, ?)",
		task.ID, task.Title, task.Completed,
	)
	return err
}

func (s *Store) List() ([]storage.Task, error) {
	rows, err := s.db.Query("SELECT id, title, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []storage.Task
	for rows.Next() {
		var task storage.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (s *Store) Update(task storage.Task) error {
	result, err := s.db.Exec(
		"UPDATE tasks SET title = ?, completed = ? WHERE id = ?",
		task.Title, task.Completed, task.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return storage.ErrTaskNotFound
	}
	return nil
}

func (s *Store) Delete(id string) error {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return storage.ErrTaskNotFound
	}
	return nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}
