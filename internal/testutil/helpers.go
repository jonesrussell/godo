// Package testutil provides testing utilities and mock implementations
package testutil

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
)

// WindowState tracks window visibility state for testing
type WindowState struct {
	mu      sync.RWMutex
	visible bool
}

func NewWindowState() *WindowState {
	return &WindowState{visible: false}
}

func (w *WindowState) Show() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.visible = true
}

func (w *WindowState) Hide() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.visible = false
}

func (w *WindowState) IsVisible() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.visible
}

// AssertWindowShown asserts that a window is shown
func AssertWindowShown(t *testing.T, state *WindowState) {
	assert.True(t, state.IsVisible())
}

// AssertWindowHidden asserts that a window is hidden
func AssertWindowHidden(t *testing.T, state *WindowState) {
	assert.False(t, state.IsVisible())
}

// WaitForWindowShown waits for a window to be shown with timeout
func WaitForWindowShown(t *testing.T, state *WindowState, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if state.IsVisible() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Error("Window was not shown within timeout")
}

// WaitForWindowHidden waits for a window to be hidden with timeout
func WaitForWindowHidden(t *testing.T, state *WindowState, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !state.IsVisible() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Error("Window was not hidden within timeout")
}

// AssertTaskExists asserts that a task exists in the store
func AssertTaskExists(t *testing.T, store *storage.MockStore, id string) {
	task, err := store.GetByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, id, task.ID)
}

// AssertTaskNotExists asserts that a task does not exist in the store
func AssertTaskNotExists(t *testing.T, store *storage.MockStore, id string) {
	_, err := store.GetByID(context.Background(), id)
	assert.Error(t, err)
}

// AssertLogContains asserts that a log message was recorded
func AssertLogContains(t *testing.T, logger logger.Logger, level string, message string) {
	mockLogger, ok := logger.(*logger.MockTestLogger)
	if !ok {
		t.Errorf("Expected logger to be *logger.MockTestLogger, got %T", logger)
		return
	}
	mockLogger.AssertCalled(t, level, message)
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
			Level: "debug",
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
				Width:       800,
				Height:      600,
				StartHidden: false,
			},
			QuickNote: config.WindowConfig{
				Width:       400,
				Height:      300,
				StartHidden: true,
			},
		},
	}
}

// WithTestContext creates a test context with timeout
func WithTestContext(t *testing.T) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
