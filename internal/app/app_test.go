package app_test

import (
	"context"
	"testing"

	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/types"
)

// MockQuickNoteUI implements ui.QuickNoteUI for testing
type MockQuickNoteUI struct{}

func (m *MockQuickNoteUI) Show(ctx context.Context) error { return nil }
func (m *MockQuickNoteUI) GetInput() <-chan string        { return make(chan string) }

func TestAppInitialization(t *testing.T) {
	// Initialize logger for tests
	if err := logger.InitializeWithConfig(types.LogConfig{
		Level:  "debug",
		Output: []string{"stdout"},
	}); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Create test configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Path: ":memory:", // Use in-memory SQLite for tests
		},
		Logging: types.LogConfig{
			Level:  "debug",
			Output: []string{"stdout"},
		},
	}

	// Initialize app using wire-generated constructor
	testApp, err := app.InitializeAppWithConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize app: %v", err)
	}

	// Verify app components
	if testApp.GetTodoService() == nil {
		t.Error("TodoService not initialized")
	}
}
