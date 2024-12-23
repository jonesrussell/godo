package app_test

import (
	"testing"

	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockQuickNoteService struct {
	showCalled bool
	hideCalled bool
}

func (m *MockQuickNoteService) Show() {
	m.showCalled = true
}

func (m *MockQuickNoteService) Hide() {
	m.hideCalled = true
}

func setupTestApp(t *testing.T) (*app.App, *MockQuickNoteService) {
	t.Helper()

	// Setup logger
	log, err := logger.New(&common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	})
	require.NoError(t, err)

	// Setup in-memory store
	store := memory.New()

	// Setup config with in-memory database
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Test App",
			Version: "0.0.1",
		},
		Database: config.DatabaseConfig{
			Path: ":memory:",
		},
	}

	// Create app with store
	testApp := app.NewApp(cfg, store, log)
	mockQuickNote := &MockQuickNoteService{}
	testApp.SetQuickNoteService(mockQuickNote)

	return testApp, mockQuickNote
}

func TestApp(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*testing.T, *app.App, *MockQuickNoteService)
	}{
		{
			name: "Save and retrieve notes",
			fn: func(t *testing.T, a *app.App, m *MockQuickNoteService) {
				err := a.SaveNote("Test note")
				require.NoError(t, err)

				notes, err := a.GetNotes()
				require.NoError(t, err)
				assert.Contains(t, notes, "Test note")
			},
		},
		{
			name: "Quick note service integration",
			fn: func(t *testing.T, a *app.App, m *MockQuickNoteService) {
				assert.False(t, m.showCalled)
				a.ShowQuickNote()
				assert.True(t, m.showCalled)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testApp, mockQuickNote := setupTestApp(t)

			// Add assertions based on observable behavior
			tt.fn(t, testApp, mockQuickNote)
		})
	}
}
