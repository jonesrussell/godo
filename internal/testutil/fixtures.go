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

const (
	// DefaultMainWindowWidth is the default width for the main window in tests
	DefaultMainWindowWidth = 800
	// DefaultMainWindowHeight is the default height for the main window in tests
	DefaultMainWindowHeight = 600
	// DefaultQuickNoteWidth is the default width for the quick note window in tests
	DefaultQuickNoteWidth = 400
	// DefaultQuickNoteHeight is the default height for the quick note window in tests
	DefaultQuickNoteHeight = 300
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
					Width:       DefaultMainWindowWidth,
					Height:      DefaultMainWindowHeight,
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
func CreateTestTask(id, content string, done bool) storage.Task {
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

// CreateTestConfig creates a test configuration with default values
func CreateTestConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			Name:    "Godo Test",
			Version: "0.1.0",
			ID:      "io.github.jonesrussell.godo.test",
		},
		Logger: common.LogConfig{
			Level:   "debug",
			Console: true,
		},
		Hotkeys: config.HotkeyConfig{
			QuickNote: common.HotkeyBinding{
				Modifiers: []string{"Ctrl"},
				Key:       "Space",
			},
		},
		Database: config.DatabaseConfig{
			Path: ":memory:",
		},
		UI: config.UIConfig{
			MainWindow: config.WindowConfig{
				Width:       DefaultMainWindowWidth,
				Height:      DefaultMainWindowHeight,
				StartHidden: false,
			},
			QuickNote: config.WindowConfig{
				Width:       DefaultQuickNoteWidth,
				Height:      DefaultQuickNoteHeight,
				StartHidden: true,
			},
		},
	}
}
