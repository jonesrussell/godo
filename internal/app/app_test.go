package di_test

import (
	"testing"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/di"
)

func TestAppInitialization(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Path: ":memory:", // Use in-memory SQLite for tests
		},
		Logging: config.LoggingConfig{
			Level:  "debug",
			Output: []string{"stdout"},
		},
	}

	// Initialize app with config
	app, err := di.InitializeAppWithConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize app: %v", err)
	}

	// Verify app components
	if app.GetTodoService() == nil {
		t.Error("TodoService not initialized")
	}

	if app.GetHotkeyManager() == nil {
		t.Error("HotkeyManager not initialized")
	}
}
