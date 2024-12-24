package sqlite

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestRunMigrations(t *testing.T) {
	t.Run("executes migrations successfully", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		defer db.Close()

		err = RunMigrations(db)
		assert.NoError(t, err)

		var tableName string
		err = db.QueryRow(`
			SELECT name FROM sqlite_master 
			WHERE type='table' AND name='tasks'
		`).Scan(&tableName)

		assert.NoError(t, err)
		assert.Equal(t, "tasks", tableName)
	})

	t.Run("handles invalid SQL", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		require.NoError(t, err)
		defer db.Close()

		// Create a migration set with invalid SQL
		ms := &migrationSet{
			migrations: []string{"INVALID SQL"},
		}

		// Run migrations using the invalid set
		err = ms.RunMigrations(db)
		assert.Error(t, err)
	})
}
