package app_test

import (
	"testing"

	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/config"
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
	application, err := app.InitializeAppWithConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize app: %v", err)
	}

	// Verify app components
	if application.GetTodoService() == nil {
		t.Error("TodoService not initialized")
	}

	if application.GetHotkeyManager() == nil {
		t.Error("HotkeyManager not initialized")
	}
}
