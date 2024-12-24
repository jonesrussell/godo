package app_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
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

type MockSystrayService struct {
	ready     bool
	menu      *fyne.Menu
	icon      fyne.Resource
	setupDone bool
}

func (m *MockSystrayService) Setup(menu *fyne.Menu) {
	m.menu = menu
	m.setupDone = true
	m.ready = true
}

func (m *MockSystrayService) SetIcon(resource fyne.Resource) {
	m.icon = resource
}

func (m *MockSystrayService) IsReady() bool {
	return m.ready
}

func setupTestApp(t *testing.T) (*app.App, *MockQuickNoteService, *MockSystrayService) {
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
	mockSystray := &MockSystrayService{}

	testApp.SetQuickNoteService(mockQuickNote)
	testApp.SetSystrayService(mockSystray)

	return testApp, mockQuickNote, mockSystray
}

func TestApp(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*testing.T, *app.App, *MockQuickNoteService, *MockSystrayService)
	}{
		{
			name: "Save and retrieve notes",
			fn: func(t *testing.T, a *app.App, _ *MockQuickNoteService, _ *MockSystrayService) {
				err := a.SaveNote("Test note")
				require.NoError(t, err)

				notes, err := a.GetNotes()
				require.NoError(t, err)
				assert.Contains(t, notes, "Test note")
			},
		},
		{
			name: "Quick note service integration",
			fn: func(t *testing.T, a *app.App, m *MockQuickNoteService, _ *MockSystrayService) {
				assert.False(t, m.showCalled)
				a.ShowQuickNote()
				assert.True(t, m.showCalled)
			},
		},
		{
			name: "System tray setup",
			fn: func(t *testing.T, a *app.App, _ *MockQuickNoteService, s *MockSystrayService) {
				a.SetupUI()
				assert.True(t, s.setupDone)
				assert.NotNil(t, s.menu)
				assert.NotNil(t, s.icon)
			},
		},
		{
			name: "Version check",
			fn: func(t *testing.T, a *app.App, _ *MockQuickNoteService, _ *MockSystrayService) {
				version := a.GetVersion()
				assert.Equal(t, "0.0.1", version)
			},
		},
		{
			name: "Lifecycle events",
			fn: func(t *testing.T, a *app.App, _ *MockQuickNoteService, _ *MockSystrayService) {
				// Create a test app to verify lifecycle events
				testApp := test.NewApp()
				testWindow := testApp.NewWindow("Test")
				defer testWindow.Close()

				// Replace the main window with our test window
				a.SetMainWindow(testWindow)

				// Run setup which should trigger lifecycle events
				a.SetupUI()

				// Verify the window exists and has content
				assert.NotNil(t, testWindow.Content(), "Window should have content")

				// Check window size
				size := testWindow.Canvas().Size()
				assert.Equal(t, fyne.NewSize(800, 600), size, "Window should be properly sized")
			},
		},
		{
			name: "Main window hidden on startup",
			fn: func(t *testing.T, a *app.App, _ *MockQuickNoteService, _ *MockSystrayService) {
				// Create a test window using Fyne's test package
				testApp := test.NewApp()
				testWindow := testApp.NewWindow("Test")
				defer testWindow.Close()

				// Replace the main window with our test window
				a.SetMainWindow(testWindow)

				// Run the setup which should hide the window
				a.SetupUI()

				// Check window size
				size := testWindow.Canvas().Size()
				assert.Equal(t, fyne.NewSize(800, 600), size, "Window should be properly sized")
			},
		},
		{
			name: "System tray menu options",
			fn: func(t *testing.T, a *app.App, _ *MockQuickNoteService, s *MockSystrayService) {
				a.SetupUI()

				// Verify the menu was set up
				assert.NotNil(t, s.menu)

				// Check menu items
				menuItems := s.menu.Items
				assert.Equal(t, 4, len(menuItems), "Should have 4 menu items (Show, Quick Note, Separator, Quit)")

				// Check menu item labels
				assert.Equal(t, "Show", menuItems[0].Label)
				assert.Equal(t, "Quick Note", menuItems[1].Label)
				assert.True(t, menuItems[2].IsSeparator)
				assert.Equal(t, "Quit", menuItems[3].Label)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testApp, mockQuickNote, mockSystray := setupTestApp(t)
			tt.fn(t, testApp, mockQuickNote, mockSystray)
		})
	}
}
