package systray

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockDesktopApp implements desktop.App
type mockDesktopApp struct {
	fyne.App
	mock.Mock
	menu *fyne.Menu
}

func (m *mockDesktopApp) SetSystemTrayMenu(menu *fyne.Menu) {
	m.menu = menu
	m.Called(menu)
}

func (m *mockDesktopApp) SetSystemTrayIcon(icon fyne.Resource) {
	m.Called(icon)
}

type mockQuickNote struct {
	mock.Mock
}

func (m *mockQuickNote) Show() {
	m.Called()
}

func (m *mockQuickNote) Hide() {
	m.Called()
}

func TestSetupSystray_QuickNoteMenuItem(t *testing.T) {
	// Create mocks
	app := &mockDesktopApp{
		App: test.NewApp(),
	}
	mainWindow := test.NewWindow(nil)
	quickNote := &mockQuickNote{}

	// Set expectations
	app.On("SetSystemTrayIcon", mock.Anything).Return()
	app.On("SetSystemTrayMenu", mock.Anything).Return()
	quickNote.On("Show").Return()

	// Setup systray
	SetupSystray(app, mainWindow, quickNote)

	// Verify menu was set
	assert.NotNil(t, app.menu, "Systray menu should be set")

	// Find Quick Note menu item
	var quickNoteItem *fyne.MenuItem
	for _, item := range app.menu.Items {
		if item.Label == "Quick Note" {
			quickNoteItem = item
			break
		}
	}

	// Verify Quick Note menu item exists
	assert.NotNil(t, quickNoteItem, "Quick Note menu item should exist")

	// Trigger Quick Note menu item
	if quickNoteItem != nil && quickNoteItem.Action != nil {
		quickNoteItem.Action()
	}

	// Verify Quick Note was shown
	quickNote.AssertCalled(t, "Show")
}

func TestSetupSystray_NonDesktopApp(t *testing.T) {
	// Test with non-desktop app
	app := test.NewApp()
	mainWindow := test.NewWindow(nil)
	quickNote := &mockQuickNote{}

	// This should not panic
	SetupSystray(app, mainWindow, quickNote)

	// Quick Note should not be called
	quickNote.AssertNotCalled(t, "Show")
}
