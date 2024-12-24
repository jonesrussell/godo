//go:build !docker
// +build !docker

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

// mockMainWindow is a test implementation of MainWindow
type mockMainWindow struct {
	shown       bool
	setupCalled bool
}

func (m *mockMainWindow) Show() {
	m.shown = true
}

func (m *mockMainWindow) Setup() {
	m.setupCalled = true
}

// mockQuickNote is a test implementation of QuickNote
type mockQuickNote struct {
	shown bool
}

func (m *mockQuickNote) Show() {
	m.shown = true
}

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
	app := NewApp(cfg, log, store)
	require.NotNil(t, app)

	// Replace windows with mocks
	mockMain := &mockMainWindow{}
	app.mainWin = mockMain
	mockQN := &mockQuickNote{}
	app.quickNote = mockQN

	// Test UI setup
	t.Run("ui setup", func(t *testing.T) {
		app.SetupUI()
		assert.True(t, mockMain.setupCalled)
	})

	// Test version
	t.Run("version", func(t *testing.T) {
		assert.Equal(t, "0.1.0", app.GetVersion())
	})

	// Test run
	t.Run("run", func(t *testing.T) {
		err := app.Run()
		require.NoError(t, err)
		assert.True(t, mockMain.shown)
	})
}
