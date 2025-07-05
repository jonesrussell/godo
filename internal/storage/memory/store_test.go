package memory_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/memory"
)

func TestMemoryStore(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*memory.Store)
		validate func(*testing.T, *memory.Store)
	}{
		{
			name:  "new store is empty",
			setup: func(_ *memory.Store) {},
			validate: func(t *testing.T, s *memory.Store) {
				tasks, err := s.List(context.Background())
				assert.NoError(t, err)
				assert.Empty(t, tasks)
			},
		},
		{
			name: "add and retrieve task",
			setup: func(s *memory.Store) {
				task := storage.Task{
					ID:        "test-1",
					Content:   "Test Task",
					Done:      false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				err := s.Add(context.Background(), task)
				require.NoError(t, err)
			},
			validate: func(t *testing.T, s *memory.Store) {
				tasks, err := s.List(context.Background())
				assert.NoError(t, err)
				assert.Len(t, tasks, 1)
				assert.Equal(t, "test-1", tasks[0].ID)
				assert.Equal(t, "Test Task", tasks[0].Content)
				assert.False(t, tasks[0].Done)
			},
		},
		{
			name: "update existing task",
			setup: func(s *memory.Store) {
				task := storage.Task{
					ID:        "test-1",
					Content:   "Original Content",
					Done:      false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				err := s.Add(context.Background(), task)
				require.NoError(t, err)

				task.Content = "Updated Content"
				task.Done = true
				err = s.Update(context.Background(), task)
				require.NoError(t, err)
			},
			validate: func(t *testing.T, s *memory.Store) {
				task, err := s.GetByID(context.Background(), "test-1")
				assert.NoError(t, err)
				assert.Equal(t, "Updated Content", task.Content)
				assert.True(t, task.Done)
			},
		},
		{
			name: "delete task",
			setup: func(s *memory.Store) {
				task := storage.Task{
					ID:        "test-1",
					Content:   "Test Task",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				err := s.Add(context.Background(), task)
				require.NoError(t, err)
				err = s.Delete(context.Background(), "test-1")
				require.NoError(t, err)
			},
			validate: func(t *testing.T, s *memory.Store) {
				tasks, err := s.List(context.Background())
				assert.NoError(t, err)
				assert.Empty(t, tasks)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := memory.New()
			tt.setup(store)
			tt.validate(t, store)
		})
	}
}

func TestMemoryStoreEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		op      func(*memory.Store) error
		wantErr error
	}{
		{
			name: "update non-existent task",
			op: func(s *memory.Store) error {
				return s.Update(context.Background(), storage.Task{ID: "nonexistent"})
			},
			wantErr: storage.ErrTaskNotFound,
		},
		{
			name: "delete non-existent task",
			op: func(s *memory.Store) error {
				return s.Delete(context.Background(), "nonexistent")
			},
			wantErr: storage.ErrTaskNotFound,
		},
		{
			name: "get non-existent task",
			op: func(s *memory.Store) error {
				_, err := s.GetByID(context.Background(), "nonexistent")
				return err
			},
			wantErr: storage.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := memory.New()
			err := tt.op(store)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestMemoryStoreConcurrent(t *testing.T) {
	store := memory.New()
	const numTasks = 100

	// Test concurrent reads and writes
	t.Run("concurrent operations", func(t *testing.T) {
		done := make(chan bool)
		ctx := context.Background()

		// Writer goroutine
		go func() {
			for i := 0; i < numTasks; i++ {
				task := storage.Task{
					ID:        fmt.Sprintf("task-%d", i),
					Content:   fmt.Sprintf("Task %d", i),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				err := store.Add(ctx, task)
				assert.NoError(t, err)
			}
			done <- true
		}()

		// Reader goroutine
		go func() {
			for i := 0; i < numTasks; i++ {
				_, _ = store.List(ctx) // Ignore errors, just testing for race conditions
			}
			done <- true
		}()

		// Updater goroutine
		go func() {
			for i := 0; i < numTasks; i++ {
				task := storage.Task{
					ID:        fmt.Sprintf("task-%d", i),
					Content:   fmt.Sprintf("Updated Task %d", i),
					Done:      true,
					UpdatedAt: time.Now(),
				}
				_ = store.Update(ctx, task) // Ignore errors, just testing for race conditions
			}
			done <- true
		}()

		// Wait for all goroutines to complete
		for i := 0; i < 3; i++ {
			<-done
		}

		// Verify final state
		tasks, err := store.List(ctx)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(tasks), numTasks)
	})
}

func TestMemoryStoreValidation(t *testing.T) {
	store := memory.New()
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name    string
		task    storage.Task
		op      string
		wantErr error
	}{
		{
			name: "valid task",
			task: storage.Task{
				ID:        "test-1",
				Content:   "Test Task",
				CreatedAt: now,
				UpdatedAt: now,
			},
			op: "add",
		},
		{
			name: "duplicate task",
			task: storage.Task{
				ID:        "test-1",
				Content:   "Duplicate Task",
				CreatedAt: now,
				UpdatedAt: now,
			},
			op:      "add",
			wantErr: storage.ErrDuplicateID,
		},
		{
			name: "empty content",
			task: storage.Task{
				ID:        "test-2",
				Content:   "",
				CreatedAt: now,
				UpdatedAt: now,
			},
			op: "add",
		},
		{
			name: "zero timestamps",
			task: storage.Task{
				ID:      "test-3",
				Content: "Test Task",
			},
			op: "add",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.op == "add" {
				err = store.Add(ctx, tt.task)
				if tt.wantErr != nil {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestMemoryStoreClose(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	// Add some tasks
	for i := 0; i < 5; i++ {
		task := storage.Task{
			ID:        fmt.Sprintf("task-%d", i),
			Content:   fmt.Sprintf("Task %d", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := store.Add(ctx, task)
		require.NoError(t, err)
	}

	// Close the store
	err := store.Close()
	assert.NoError(t, err)

	// Verify store is empty after close
	tasks, err := store.List(ctx)
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}
