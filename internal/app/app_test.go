package app_test

import (
	"testing"

	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
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
	log, err := logger.NewZapLogger(&logger.Config{
		Level:   "debug",
		Console: true,
	})
	require.NoError(t, err)

	// Setup config
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Test App",
			Version: "0.0.1",
		},
	}

	// Create app
	testApp := app.NewApp(cfg, nil, log)
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
			name: "Setup UI initializes components",
			fn: func(t *testing.T, a *app.App, m *MockQuickNoteService) {
				a.SetupUI()
				// Add assertions based on observable behavior
			},
		},
		{
			name: "Cleanup closes resources",
			fn: func(t *testing.T, a *app.App, m *MockQuickNoteService) {
				a.Cleanup()
				// Add assertions for cleanup
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testApp, mockQuickNote := setupTestApp(t)
			tt.fn(t, testApp, mockQuickNote)
		})
	}
}

func TestApp_QuickNoteIntegration(t *testing.T) {
	testApp, mockQuickNote := setupTestApp(t)

	// Test that quick note service is properly integrated
	assert.False(t, mockQuickNote.showCalled)
	testApp.ShowQuickNote() // Add this method to App
	assert.True(t, mockQuickNote.showCalled)
}
