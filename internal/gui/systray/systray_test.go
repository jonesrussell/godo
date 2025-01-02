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

type mockMainWindow struct {
	fyne.Window
	mock.Mock
}

func (m *mockMainWindow) Show() {
	m.Called()
}

func (m *mockMainWindow) Hide() {
	m.Called()
}

func (m *mockMainWindow) CenterOnScreen() {
	m.Called()
}

func (m *mockMainWindow) SetContent(content fyne.CanvasObject) {
	m.Called(content)
}

func (m *mockMainWindow) Resize(size fyne.Size) {
	m.Called(size)
}

func (m *mockMainWindow) GetWindow() fyne.Window {
	args := m.Called()
	return args.Get(0).(fyne.Window)
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

func (m *mockQuickNote) CenterOnScreen() {
	m.Called()
}

func TestSetupSystray_QuickNoteMenuItem(t *testing.T) {
	// Create mocks
	app := &mockDesktopApp{
		App: test.NewApp(),
	}
	mainWindow := &mockMainWindow{
		Window: test.NewWindow(nil),
	}
	quickNote := &mockQuickNote{}

	// Set expectations
	app.On("SetSystemTrayIcon", mock.Anything).Return()
	app.On("SetSystemTrayMenu", mock.Anything).Return()
	mainWindow.On("Show").Return()
	mainWindow.On("GetWindow").Return(test.NewWindow(nil))
	mainWindow.On("CenterOnScreen").Return()
	quickNote.On("Show").Return()
	quickNote.On("CenterOnScreen").Return()

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
	mainWindow := &mockMainWindow{
		Window: test.NewWindow(nil),
	}
	quickNote := &mockQuickNote{}

	// Set expectations
	mainWindow.On("GetWindow").Return(test.NewWindow(nil))

	// This should not panic
	SetupSystray(app, mainWindow, quickNote)

	// Quick Note should not be called
	quickNote.AssertNotCalled(t, "Show")
}
