package app

import (
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	// Create test dependencies
	cfg := &config.Config{
		Logger: common.LogConfig{
			Level:   "debug",
			Console: true,
		},
		App: config.AppConfig{
			Name:    "Test App",
			Version: "0.1.0",
		},
	}
	log := logger.NewTestLogger(t)
	store := memory.New()

	// Create app
	app := NewApp(cfg, store, log)
	require.NotNil(t, app)

	// Test note operations
	t.Run("note operations", func(t *testing.T) {
		// Add note
		err := app.SaveNote("Test note")
		require.NoError(t, err)

		// Get notes
		notes, err := app.GetNotes()
		require.NoError(t, err)
		assert.Contains(t, notes, "Test note")

		// Verify version
		assert.Equal(t, "0.1.0", app.GetVersion())
	})

	// Clean up
	app.Cleanup()
}
