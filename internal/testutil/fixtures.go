package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// TestFixture holds common test dependencies
type TestFixture struct {
	T          *testing.T
	App        fyne.Window
	Store      *storage.MockStore
	Logger     logger.Logger
	Config     *config.Config
	Context    context.Context
	CleanupFns []func()
}

// NewTestFixture creates a new test fixture with all common dependencies
func NewTestFixture(t *testing.T) *TestFixture {
	t.Helper()

	f := &TestFixture{
		T:       t,
		Store:   storage.NewMockStore(),
		Logger:  logger.NewMockTestLogger(t),
		Context: context.Background(),
		Config: &config.Config{
			UI: config.UIConfig{
				MainWindow: config.WindowConfig{
					Width:       800,
					Height:      600,
					StartHidden: false,
				},
			},
		},
	}

	f.App = test.NewWindow(nil)
	f.CleanupFns = append(f.CleanupFns, func() {
		f.App.Close()
	})

	t.Cleanup(f.Cleanup)
	return f
}

// Cleanup runs all cleanup functions
func (f *TestFixture) Cleanup() {
	for _, fn := range f.CleanupFns {
		fn()
	}
}

// CreateTestTask creates a task for testing
func CreateTestTask(id string, content string, done bool) storage.Task {
	now := time.Now()
	return storage.Task{
		ID:        id,
		Content:   content,
		Done:      done,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateTestTasks creates multiple tasks for testing
func CreateTestTasks(count int) []storage.Task {
	tasks := make([]storage.Task, count)
	for i := 0; i < count; i++ {
		tasks[i] = CreateTestTask(
			fmt.Sprintf("test-%d", i+1),
			fmt.Sprintf("Test Task %d", i+1),
			false,
		)
	}
	return tasks
}

// WithTestConfig sets up a test configuration
func WithTestConfig() *config.Config {
	return &config.Config{
		UI: config.UIConfig{
			MainWindow: config.WindowConfig{
				Width:       800,
				Height:      600,
				StartHidden: false,
			},
		},
		Logger: common.LogConfig{
			Level:   "debug",
			Console: true,
			Output:  []string{"stdout"},
		},
	}
}
