package systray

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct {
	logger.Logger
	debugCalls []string
	warnCalls  []string
}

func (m *mockLogger) Debug(msg string, args ...interface{}) {
	m.debugCalls = append(m.debugCalls, msg)
}

func (m *mockLogger) Warn(msg string, args ...interface{}) {
	m.warnCalls = append(m.warnCalls, msg)
}

func TestNew(t *testing.T) {
	t.Run("creates service with app", func(t *testing.T) {
		mockApp := test.NewApp()
		mockLog := &mockLogger{}

		service := New(mockApp, mockLog)

		assert.NotNil(t, service)
		assert.Equal(t, mockApp, service.app)
		assert.Equal(t, mockLog, service.log)
		assert.False(t, service.ready)
		assert.Nil(t, service.menu)
		assert.Nil(t, service.icon)
	})
}

func TestSetup(t *testing.T) {
	t.Run("sets menu and ready state", func(t *testing.T) {
		mockApp := test.NewApp()
		mockLog := &mockLogger{}
		service := New(mockApp, mockLog)

		menu := fyne.NewMenu("Test Menu")
		service.Setup(menu)

		assert.True(t, service.ready)
		assert.Equal(t, menu, service.menu)
	})
}

func TestSetIcon(t *testing.T) {
	t.Run("sets icon", func(t *testing.T) {
		mockApp := test.NewApp()
		mockLog := &mockLogger{}
		service := New(mockApp, mockLog)

		icon := fyne.NewStaticResource("test.png", []byte{})
		service.SetIcon(icon)

		assert.Equal(t, icon, service.icon)
	})
}

func TestIsReady(t *testing.T) {
	t.Run("returns ready state", func(t *testing.T) {
		mockApp := test.NewApp()
		mockLog := &mockLogger{}
		service := New(mockApp, mockLog)

		assert.False(t, service.IsReady())

		service.ready = true
		assert.True(t, service.IsReady())
	})
}

func TestIntegration(t *testing.T) {
	t.Run("full systray lifecycle", func(t *testing.T) {
		mockApp := test.NewApp()
		mockLog := &mockLogger{}
		service := New(mockApp, mockLog)

		// Initial state
		assert.False(t, service.IsReady())
		assert.Nil(t, service.menu)
		assert.Nil(t, service.icon)

		// Setup menu
		menu := fyne.NewMenu("Test Menu",
			fyne.NewMenuItem("Item 1", nil),
			fyne.NewMenuItem("Item 2", nil),
		)
		service.Setup(menu)

		// Verify menu setup
		assert.True(t, service.IsReady())
		assert.Equal(t, menu, service.menu)
		assert.Len(t, menu.Items, 2)

		// Set icon
		icon := fyne.NewStaticResource("test.png", []byte{})

		service.SetIcon(icon)

		// Verify icon setup
		assert.Equal(t, icon, service.icon)
	})
}
