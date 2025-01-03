package gui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestWindowImplementation(t *testing.T) {
	// Create test window
	window := test.NewWindow(nil)
	defer window.Close()

	// Create window implementation
	impl := NewWindow(window)

	// Test Show
	t.Run("Show", func(t *testing.T) {
		impl.Show()
		assert.True(t, window.Visible())
	})

	// Test Hide
	t.Run("Hide", func(t *testing.T) {
		impl.Hide()
		assert.False(t, window.Visible())
	})

	// Test GetWindow
	t.Run("GetWindow", func(t *testing.T) {
		assert.Equal(t, window, impl.GetWindow())
	})

	// Test SetOnClosed
	t.Run("SetOnClosed", func(t *testing.T) {
		called := false
		impl.SetOnClosed(func() {
			called = true
		})

		window.Close()
		assert.True(t, called)
	})
}

func TestWindowManager(t *testing.T) {
	// Create test window
	window := test.NewWindow(nil)
	defer window.Close()

	// Create window manager
	manager := NewWindow(window)

	// Test Show
	t.Run("Show", func(t *testing.T) {
		manager.Show()
		assert.True(t, window.Visible())
	})

	// Test Hide
	t.Run("Hide", func(t *testing.T) {
		manager.Hide()
		assert.False(t, window.Visible())
	})

	// Test Close
	t.Run("Close", func(t *testing.T) {
		manager.Close()
		assert.False(t, window.Visible())
	})
}
