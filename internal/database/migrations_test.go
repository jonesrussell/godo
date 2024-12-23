package database

import (
	"database/sql"
	"testing"

	"github.com/jonesrussell/godo/internal/logger"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestRunMigrations(t *testing.T) {
	// Create an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	log := logger.NewTestLogger(t)

	// Run migrations
	err = RunMigrations(db, log)
	assert.NoError(t, err)

	// Verify the todos table was created with correct schema
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='todos'").Scan(&tableName)
	assert.NoError(t, err)
	assert.Equal(t, "todos", tableName)

	// Verify table schema
	rows, err := db.Query("PRAGMA table_info(todos)")
	assert.NoError(t, err)
	defer rows.Close()

	expectedColumns := map[string]struct{}{
		"id":          {},
		"title":       {},
		"description": {},
		"completed":   {},
		"created_at":  {},
		"updated_at":  {},
	}

	for rows.Next() {
		var (
			cid     int
			name    string
			typ     string
			notnull int
			dfltVal sql.NullString
			pk      int
		)
		err := rows.Scan(&cid, &name, &typ, &notnull, &dfltVal, &pk)
		assert.NoError(t, err)
		_, exists := expectedColumns[name]
		assert.True(t, exists, "Unexpected column: %s", name)
		delete(expectedColumns, name)
	}
	assert.NoError(t, rows.Err())
	assert.Empty(t, expectedColumns, "Missing columns: %v", expectedColumns)
}

func TestRunMigrations_Error(t *testing.T) {
	// Create an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	log := logger.NewTestLogger(t)

	// Close the database to force an error
	db.Close()

	// Run migrations should fail
	err = RunMigrations(db, log)
	assert.Error(t, err)
}
