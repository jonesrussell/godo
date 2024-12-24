package systray

import (
	"net/url"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
)

type mockResource struct {
	name    string
	content []byte
}

func (m *mockResource) Name() string {
	return m.name
}

func (m *mockResource) Content() []byte {
	return m.content
}

type mockDesktopApp struct {
	app  fyne.App
	menu *fyne.Menu
	icon fyne.Resource
}

func newMockDesktopApp() *mockDesktopApp {
	return &mockDesktopApp{
		app: test.NewApp(),
	}
}

func (m *mockDesktopApp) NewWindow(title string) fyne.Window {
	return m.app.NewWindow(title)
}

func (m *mockDesktopApp) Run() {
	m.app.Run()
}

func (m *mockDesktopApp) Quit() {
	m.app.Quit()
}

func (m *mockDesktopApp) Driver() fyne.Driver {
	return m.app.Driver()
}

func (m *mockDesktopApp) UniqueID() string {
	return m.app.UniqueID()
}

func (m *mockDesktopApp) SendNotification(*fyne.Notification) {
	// no-op for tests
}

func (m *mockDesktopApp) Settings() fyne.Settings {
	return m.app.Settings()
}

func (m *mockDesktopApp) Storage() fyne.Storage {
	return m.app.Storage()
}

func (m *mockDesktopApp) Preferences() fyne.Preferences {
	return m.app.Preferences()
}

func (m *mockDesktopApp) CloudProvider() fyne.CloudProvider {
	return m.app.CloudProvider()
}

func (m *mockDesktopApp) Icon() fyne.Resource {
	return m.app.Icon()
}

func (m *mockDesktopApp) Lifecycle() fyne.Lifecycle {
	return m.app.Lifecycle()
}

func (m *mockDesktopApp) Metadata() fyne.AppMetadata {
	return fyne.AppMetadata{
		ID:      "test.app",
		Name:    "Test App",
		Version: "1.0.0",
	}
}

func (m *mockDesktopApp) OpenURL(_ *url.URL) error {
	return nil
}

func (m *mockDesktopApp) SetCloudProvider(_ fyne.CloudProvider) {
	// no-op for tests
}

func (m *mockDesktopApp) SetIcon(icon fyne.Resource) {
	m.app.SetIcon(icon)
}

func (m *mockDesktopApp) SetSystemTrayMenu(menu *fyne.Menu) {
	m.menu = menu
}

func (m *mockDesktopApp) SetSystemTrayIcon(icon fyne.Resource) {
	m.icon = icon
}

func TestInterface(t *testing.T) {
	t.Run("implementation satisfies interface", func(_ *testing.T) {
		var _ Interface = (*Systray)(nil)
	})
}

func TestSystray(t *testing.T) {
	mockApp := newMockDesktopApp()
	defer mockApp.Quit()

	log, err := logger.New(&common.LogConfig{
		Level:   "debug",
		Console: true,
	})
	assert.NoError(t, err)

	t.Run("Setup sets menu", func(t *testing.T) {
		svc := New(mockApp, log)
		menu := &fyne.Menu{
			Label: "Test Menu",
			Items: []*fyne.MenuItem{
				{Label: "Test Item"},
			},
		}

		svc.Setup(menu)
		assert.True(t, svc.IsReady())
		assert.Equal(t, menu, svc.menu)
		assert.Equal(t, menu, mockApp.menu)
	})

	t.Run("SetIcon sets icon", func(t *testing.T) {
		svc := New(mockApp, log)
		icon := &mockResource{
			name:    "test.png",
			content: []byte("test content"),
		}

		svc.SetIcon(icon)
		assert.Equal(t, icon, svc.icon)
		assert.Equal(t, icon, mockApp.icon)
	})

	t.Run("IsReady returns ready state", func(t *testing.T) {
		svc := New(mockApp, log)
		assert.False(t, svc.IsReady())

		svc.Setup(&fyne.Menu{})
		assert.True(t, svc.IsReady())
	})
}

func TestIntegration(t *testing.T) {
	t.Run("full lifecycle", func(t *testing.T) {
		mockApp := newMockDesktopApp()
		defer mockApp.Quit()

		log, err := logger.New(&common.LogConfig{
			Level:   "debug",
			Console: true,
		})
		assert.NoError(t, err)

		svc := New(mockApp, log)
		menu := &fyne.Menu{
			Label: "Test Menu",
			Items: []*fyne.MenuItem{
				{Label: "Test Item"},
			},
		}
		icon := &mockResource{
			name:    "test.png",
			content: []byte("test content"),
		}

		// Test initial state
		assert.False(t, svc.IsReady())
		assert.Nil(t, svc.menu)
		assert.Nil(t, svc.icon)

		// Test setup
		svc.Setup(menu)
		assert.True(t, svc.IsReady())
		assert.Equal(t, menu, svc.menu)
		assert.Equal(t, menu, mockApp.menu)

		// Test icon
		svc.SetIcon(icon)
		assert.Equal(t, icon, svc.icon)
		assert.Equal(t, icon, mockApp.icon)
		assert.True(t, svc.IsReady()) // Ready state should not be affected by icon
	})
}
