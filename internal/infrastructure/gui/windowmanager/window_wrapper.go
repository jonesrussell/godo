// Package windowmanager provides centralized window management for the application
package windowmanager

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"

	"github.com/jonesrussell/godo/internal/infrastructure/logger"
)

// WindowWrapper wraps a Fyne window with enhanced functionality
type WindowWrapper struct {
	windowID   string
	windowType WindowType
	fyneWindow fyne.Window
	log        logger.Logger

	// State management
	state   WindowState
	stateMu sync.RWMutex

	// Lifecycle handlers
	lifecycleHandlers []WindowLifecycle
	lifecycleMu       sync.RWMutex

	// Focus management
	focusManager *FocusManager
}

// NewWindowWrapper creates a new window wrapper
func NewWindowWrapper(
	windowID string,
	windowType WindowType,
	fyneWindow fyne.Window,
	log logger.Logger,
) *WindowWrapper {
	ww := &WindowWrapper{
		windowID:     windowID,
		windowType:   windowType,
		fyneWindow:   fyneWindow,
		log:          log,
		focusManager: NewFocusManager(fyneWindow, log),
	}

	// Set up window event handlers
	ww.setupEventHandlers()

	return ww
}

// Show displays the window and notifies lifecycle handlers
func (ww *WindowWrapper) Show() {
	fyne.Do(func() {
		ww.fyneWindow.Show()
		ww.fyneWindow.CenterOnScreen()

		// Update state
		ww.stateMu.Lock()
		ww.state.Visible = true
		ww.state.LastPosition = time.Now()
		ww.stateMu.Unlock()

		// Notify lifecycle handlers
		ww.notifyLifecycleHandlers(func(handler WindowLifecycle) {
			handler.OnShow()
		})

		ww.log.Debug("Window shown", "window_id", ww.windowID)
	})
}

// Hide hides the window and notifies lifecycle handlers
func (ww *WindowWrapper) Hide() {
	fyne.Do(func() {
		ww.fyneWindow.Hide()

		// Update state
		ww.stateMu.Lock()
		ww.state.Visible = false
		ww.stateMu.Unlock()

		// Notify lifecycle handlers
		ww.notifyLifecycleHandlers(func(handler WindowLifecycle) {
			handler.OnHide()
		})

		ww.log.Debug("Window hidden", "window_id", ww.windowID)
	})
}

// Close closes the window and notifies lifecycle handlers
func (ww *WindowWrapper) Close() {
	fyne.Do(func() {
		// Notify lifecycle handlers before closing
		ww.notifyLifecycleHandlers(func(handler WindowLifecycle) {
			handler.OnClose()
		})

		ww.fyneWindow.Close()
		ww.log.Debug("Window closed", "window_id", ww.windowID)
	})
}

// GetWindow returns the underlying Fyne window
func (ww *WindowWrapper) GetWindow() fyne.Window {
	return ww.fyneWindow
}

// GetState returns the current window state
func (ww *WindowWrapper) GetState() WindowState {
	ww.stateMu.RLock()
	defer ww.stateMu.RUnlock()
	// Only store width/height, not position
	if ww.fyneWindow != nil {
		size := ww.fyneWindow.Canvas().Size()
		state := ww.state
		state.Width = int(size.Width)
		state.Height = int(size.Height)
		return state
	}
	return ww.state
}

// SetState sets the window state
func (ww *WindowWrapper) SetState(state WindowState) {
	ww.stateMu.Lock()
	ww.state = state
	ww.stateMu.Unlock()
	fyne.Do(func() {
		if ww.fyneWindow != nil {
			// Only set size
			ww.fyneWindow.Resize(fyne.NewSize(float32(state.Width), float32(state.Height)))
			if state.Visible {
				ww.fyneWindow.Show()
			} else {
				ww.fyneWindow.Hide()
			}
		}
	})
	ww.log.Debug("Window state set", "window_id", ww.windowID, "state", state)
}

// BringToFront brings the window to the front (no-op, not supported by Fyne)
func (ww *WindowWrapper) BringToFront() {
	// Not supported
}

// RequestFocus requests focus for the window
func (ww *WindowWrapper) RequestFocus() {
	ww.focusManager.RequestFocus()
}

// IsValid checks if the window is still valid
func (ww *WindowWrapper) IsValid() bool {
	return ww.fyneWindow != nil
}

// GetType returns the window type
func (ww *WindowWrapper) GetType() WindowType {
	return ww.windowType
}

// AddLifecycleHandler adds a lifecycle handler
func (ww *WindowWrapper) AddLifecycleHandler(handler WindowLifecycle) {
	ww.lifecycleMu.Lock()
	defer ww.lifecycleMu.Unlock()

	ww.lifecycleHandlers = append(ww.lifecycleHandlers, handler)
}

// RemoveLifecycleHandler removes a lifecycle handler
func (ww *WindowWrapper) RemoveLifecycleHandler(handler WindowLifecycle) {
	ww.lifecycleMu.Lock()
	defer ww.lifecycleMu.Unlock()

	for i, h := range ww.lifecycleHandlers {
		if h == handler {
			ww.lifecycleHandlers = append(ww.lifecycleHandlers[:i], ww.lifecycleHandlers[i+1:]...)
			break
		}
	}
}

// setupEventHandlers sets up window event handlers
func (ww *WindowWrapper) setupEventHandlers() {
	if ww.fyneWindow == nil {
		return
	}

	// Set up close intercept to use Hide instead of Close
	ww.fyneWindow.SetCloseIntercept(func() {
		ww.Hide()
	})

	// Set up focus events
	ww.fyneWindow.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		// Handle focus events
		if ke.Name == fyne.KeyTab {
			ww.focusManager.HandleTabKey(ke)
		}
	})
}

// notifyLifecycleHandlers notifies all lifecycle handlers
func (ww *WindowWrapper) notifyLifecycleHandlers(action func(WindowLifecycle)) {
	ww.lifecycleMu.RLock()
	handlers := ww.lifecycleHandlers
	ww.lifecycleMu.RUnlock()

	for _, handler := range handlers {
		action(handler)
	}
}
