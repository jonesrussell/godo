package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*MemoryStore)
		validate func(*testing.T, *MemoryStore)
	}{
		{
			name:  "new store is empty",
			setup: func(s *MemoryStore) {},
			validate: func(t *testing.T, s *MemoryStore) {
				tasks, err := s.List()
				assert.NoError(t, err)
				assert.Empty(t, tasks)
			},
		},
		{
			name: "add and retrieve task",
			setup: func(s *MemoryStore) {
				task := Task{
					ID:        "test-1",
					Content:   "Test Task",
					Done:      false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				err := s.Add(task)
				require.NoError(t, err)
			},
			validate: func(t *testing.T, s *MemoryStore) {
				tasks, err := s.List()
				assert.NoError(t, err)
				assert.Len(t, tasks, 1)
				assert.Equal(t, "test-1", tasks[0].ID)
				assert.Equal(t, "Test Task", tasks[0].Content)
				assert.False(t, tasks[0].Done)
			},
		},
		{
			name: "update existing task",
			setup: func(s *MemoryStore) {
				task := Task{
					ID:        "test-1",
					Content:   "Original Content",
					Done:      false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				err := s.Add(task)
				require.NoError(t, err)

				task.Content = "Updated Content"
				task.Done = true
				err = s.Update(task)
				require.NoError(t, err)
			},
			validate: func(t *testing.T, s *MemoryStore) {
				task, err := s.GetByID("test-1")
				assert.NoError(t, err)
				assert.Equal(t, "Updated Content", task.Content)
				assert.True(t, task.Done)
			},
		},
		{
			name: "delete task",
			setup: func(s *MemoryStore) {
				task := Task{
					ID:        "test-1",
					Content:   "Test Task",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				err := s.Add(task)
				require.NoError(t, err)
				err = s.Delete("test-1")
				require.NoError(t, err)
			},
			validate: func(t *testing.T, s *MemoryStore) {
				tasks, err := s.List()
				assert.NoError(t, err)
				assert.Empty(t, tasks)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryStore()
			tt.setup(store)
			tt.validate(t, store)
		})
	}
}

func TestMemoryStoreEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		op      func(*MemoryStore) error
		wantErr error
	}{
		{
			name: "update non-existent task",
			op: func(s *MemoryStore) error {
				return s.Update(Task{ID: "nonexistent"})
			},
			wantErr: ErrTaskNotFound,
		},
		{
			name: "delete non-existent task",
			op: func(s *MemoryStore) error {
				return s.Delete("nonexistent")
			},
			wantErr: ErrTaskNotFound,
		},
		{
			name: "get non-existent task",
			op: func(s *MemoryStore) error {
				_, err := s.GetByID("nonexistent")
				return err
			},
			wantErr: ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryStore()
			err := tt.op(store)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestMemoryStoreConcurrent(t *testing.T) {
	store := NewMemoryStore()
	const numTasks = 100

	// Test concurrent reads and writes
	t.Run("concurrent operations", func(t *testing.T) {
		done := make(chan bool)

		// Writer goroutine
		go func() {
			for i := 0; i < numTasks; i++ {
				task := Task{
					ID:        fmt.Sprintf("task-%d", i),
					Content:   fmt.Sprintf("Task %d", i),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				err := store.Add(task)
				assert.NoError(t, err)
			}
			done <- true
		}()

		// Reader goroutine
		go func() {
			for i := 0; i < numTasks; i++ {
				_, _ = store.List() // Ignore errors, just testing for race conditions
			}
			done <- true
		}()

		// Updater goroutine
		go func() {
			for i := 0; i < numTasks; i++ {
				task := Task{
					ID:        fmt.Sprintf("task-%d", i),
					Content:   fmt.Sprintf("Updated Task %d", i),
					Done:      true,
					UpdatedAt: time.Now(),
				}
				_ = store.Update(task) // Ignore errors, just testing for race conditions
			}
			done <- true
		}()

		// Wait for all goroutines to complete
		for i := 0; i < 3; i++ {
			<-done
		}

		// Verify final state
		tasks, err := store.List()
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(tasks), numTasks)
	})
}

func TestMemoryStoreValidation(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()

	tests := []struct {
		name    string
		task    Task
		op      string
		wantErr error
	}{
		{
			name: "valid task",
			task: Task{
				ID:        "test-1",
				Content:   "Test Task",
				CreatedAt: now,
				UpdatedAt: now,
			},
			op: "add",
		},
		{
			name: "duplicate task",
			task: Task{
				ID:        "test-1",
				Content:   "Duplicate Task",
				CreatedAt: now,
				UpdatedAt: now,
			},
			op:      "add",
			wantErr: ErrDuplicateID,
		},
		{
			name: "empty content",
			task: Task{
				ID:        "test-2",
				Content:   "",
				CreatedAt: now,
				UpdatedAt: now,
			},
			op: "add",
		},
		{
			name: "zero timestamps",
			task: Task{
				ID:      "test-3",
				Content: "Test Task",
			},
			op: "add",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch tt.op {
			case "add":
				err = store.Add(tt.task)
			case "update":
				err = store.Update(tt.task)
			}

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				// Verify task was stored correctly
				stored, err := store.GetByID(tt.task.ID)
				assert.NoError(t, err)
				assert.Equal(t, tt.task.ID, stored.ID)
				assert.Equal(t, tt.task.Content, stored.Content)
			}
		})
	}
}

func TestMemoryStoreClose(t *testing.T) {
	store := NewMemoryStore()

	// Add some tasks
	task := Task{
		ID:        "test-1",
		Content:   "Test Task",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := store.Add(task)
	require.NoError(t, err)

	// Close should be a no-op but shouldn't fail
	err = store.Close()
	assert.NoError(t, err)

	// Store should still be usable after close
	tasks, err := store.List()
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
}
