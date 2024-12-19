package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSQLiteDB(t *testing.T) {
	// Initialize logger with test config
	logConfig := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}

	if _, err := logger.Initialize(logConfig); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	tests := []struct {
		name    string
		dbPath  string
		wantErr bool
	}{
		{
			name:    "creates new database",
			dbPath:  dbPath,
			wantErr: false,
		},
		{
			name:    "fails with invalid path",
			dbPath:  filepath.Join(os.TempDir(), "nonexistent", "path", "test.db"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			db, err := NewSQLiteDB(tt.dbPath)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, db)

				// Verify we can create and query todos
				_, err = db.Exec(`INSERT INTO todos (title, description, completed) 
					VALUES (?, ?, ?)`, "test todo", "test description", false)
				assert.NoError(t, err)

				var count int
				err = db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&count)
				assert.NoError(t, err)
				assert.Equal(t, 1, count)

				// Cleanup
				db.Close()
			}
		})
	}
}

func TestRunMigrations(t *testing.T) {
	// Initialize logger with test config
	logConfig := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}

	if _, err := logger.Initialize(logConfig); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "migrations_test.db")
	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	// Test
	err = RunMigrations(db)
	assert.NoError(t, err)

	// Verify migrations by checking table structure
	var tableName string
	err = db.QueryRow(`SELECT name FROM sqlite_master 
		WHERE type='table' AND name='todos'`).Scan(&tableName)
	assert.NoError(t, err)
	assert.Equal(t, "todos", tableName)

	// Verify column structure
	columns, err := db.Query(`PRAGMA table_info(todos)`)
	require.NoError(t, err)
	defer columns.Close()

	expectedColumns := map[string]string{
		"id":          "INTEGER",
		"title":       "TEXT",
		"description": "TEXT",
		"completed":   "BOOLEAN",
		"created_at":  "DATETIME",
		"updated_at":  "DATETIME",
	}

	for columns.Next() {
		var (
			cid      int
			name     string
			dataType string
			notNull  bool
			dfltVal  sql.NullString
			pk       bool
		)
		err := columns.Scan(&cid, &name, &dataType, &notNull, &dfltVal, &pk)
		require.NoError(t, err)

		expectedType, exists := expectedColumns[name]
		assert.True(t, exists, "Unexpected column: %s", name)
		assert.Equal(t, expectedType, dataType, "Wrong type for column %s", name)
	}
}

func TestEnsureDataDir(t *testing.T) {
	tmpDir := t.TempDir()
	tests := []struct {
		name    string
		dbPath  string
		wantErr bool
	}{
		{
			name:    "creates directory if not exists",
			dbPath:  filepath.Join(tmpDir, "subdir", "test.db"),
			wantErr: false,
		},
		{
			name:    "succeeds with existing directory",
			dbPath:  filepath.Join(tmpDir, "test.db"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ensureDataDir(tt.dbPath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.DirExists(t, filepath.Dir(tt.dbPath))
			}
		})
	}
}
