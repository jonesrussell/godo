package sqlite

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunMigrations(t *testing.T) {
	// Create a temporary database file
	dbFile := "test.db"
	defer os.Remove(dbFile)

	// Open database connection
	db, err := sql.Open("sqlite", dbFile)
	require.NoError(t, err)
	defer db.Close()

	// Run migrations
	err = RunMigrations(db)
	require.NoError(t, err)

	// Verify table exists
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tasks'").Scan(&tableName)
	require.NoError(t, err)
	assert.Equal(t, "tasks", tableName)

	// Verify index exists
	var indexName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name='idx_tasks_created_at'").Scan(&indexName)
	require.NoError(t, err)
	assert.Equal(t, "idx_tasks_created_at", indexName)

	// Verify table schema
	rows, err := db.Query("PRAGMA table_info(tasks)")
	require.NoError(t, err)
	defer rows.Close()

	expectedColumns := map[string]string{
		"id":         "TEXT",
		"content":    "TEXT",
		"done":       "BOOLEAN",
		"created_at": "DATETIME",
		"updated_at": "DATETIME",
	}

	for rows.Next() {
		var cid int
		var name, typ string
		var notnull, pk int
		var defaultValue interface{}
		err := rows.Scan(&cid, &name, &typ, &notnull, &defaultValue, &pk)
		require.NoError(t, err)

		expectedType, ok := expectedColumns[name]
		assert.True(t, ok, "Unexpected column: %s", name)
		assert.Equal(t, expectedType, typ)
	}
}
