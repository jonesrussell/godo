// internal/repository/todo_test.go
package repository_test

import (
	"context"
	"testing"

	"github.com/jonesrussell/godo/internal/database"
	"github.com/jonesrussell/godo/internal/repository"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*database.DB, func()) {
	db, err := database.NewTestDB()
	assert.NoError(t, err)

	return db, func() {
		db.Close()
	}
}

func TestTodoRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewTodoRepository(db)
	ctx := context.Background()

	// Test cases here
}
