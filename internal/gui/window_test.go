package gui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestWindowImplementation(t *testing.T) {
	// Create test window
	window := test.NewWindow(nil)
	defer window.Close()

	// Create window implementation
	impl := NewWindow(window)
	visible := false

	// Test Show
	t.Run("Show", func(t *testing.T) {
		impl.Show()
		visible = true
		assert.True(t, visible)
	})

	// Test Hide
	t.Run("Hide", func(t *testing.T) {
		impl.Hide()
		visible = false
		assert.False(t, visible)
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
	visible := false

	// Test Show
	t.Run("Show", func(t *testing.T) {
		manager.Show()
		visible = true
		assert.True(t, visible)
	})

	// Test Hide
	t.Run("Hide", func(t *testing.T) {
		manager.Hide()
		visible = false
		assert.False(t, visible)
	})

	// Test Close
	t.Run("Close", func(t *testing.T) {
		manager.Close()
		visible = false
		assert.False(t, visible)
	})
}
