// Package windowmanager provides centralized window management for the application
package windowmanager

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
)

// WindowType represents different types of windows in the application
type WindowType string

const (
	WindowTypeMain      WindowType = "main"
	WindowTypeQuickNote WindowType = "quicknote"
	WindowTypeDialog    WindowType = "dialog"
)

// WindowState represents the state of a window
type WindowState struct {
	X, Y, Width, Height int
	Maximized           bool
	Visible             bool
	LastPosition        time.Time
}

// Window represents a managed window with lifecycle events
type Window interface {
	Show()
	Hide()
	Close()
	GetWindow() fyne.Window
	GetState() WindowState
	SetState(state WindowState)
	BringToFront()
	RequestFocus()
	IsValid() bool
	GetType() WindowType
}

// WindowLifecycle defines lifecycle events for windows
type WindowLifecycle interface {
	OnShow()
	OnHide()
	OnClose()
	OnFocus()
	OnBlur()
}

// WindowManager manages all windows in the application
type WindowManager struct {
	app    fyne.App
	store  storage.TaskStore
	log    logger.Logger
	config *config.Config

	// Window registry
	windows   map[string]Window
	windowsMu sync.RWMutex

	// Active window tracking
	activeWindow Window
	activeMu     sync.RWMutex

	// State persistence
	stateStore WindowStateStore

	// Lifecycle handlers
	lifecycleHandlers map[WindowType][]WindowLifecycle
	lifecycleMu       sync.RWMutex
}

// WindowStateStore handles persistence of window states
type WindowStateStore interface {
	SaveState(windowID string, state WindowState) error
	LoadState(windowID string) (WindowState, error)
	DeleteState(windowID string) error
}

// NewWindowManager creates a new window manager
func NewWindowManager(
	app fyne.App,
	store storage.TaskStore,
	log logger.Logger,
	config *config.Config,
	stateStore WindowStateStore,
) *WindowManager {
	return &WindowManager{
		app:               app,
		store:             store,
		log:               log,
		config:            config,
		windows:           make(map[string]Window),
		stateStore:        stateStore,
		lifecycleHandlers: make(map[WindowType][]WindowLifecycle),
	}
}

// CreateWindow creates a new window of the specified type
func (wm *WindowManager) CreateWindow(windowType WindowType, config WindowConfig) (Window, error) {
	wm.windowsMu.Lock()
	defer wm.windowsMu.Unlock()

	windowID := fmt.Sprintf("%s_%d", windowType, time.Now().UnixNano())

	var window Window
	var err error

	switch windowType {
	case WindowTypeMain:
		window, err = wm.createMainWindow(windowID, config)
	case WindowTypeQuickNote:
		window, err = wm.createQuickNoteWindow(windowID, config)
	case WindowTypeDialog:
		window, err = wm.createDialogWindow(windowID, config)
	default:
		return nil, fmt.Errorf("unsupported window type: %s", windowType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create window: %w", err)
	}

	wm.windows[windowID] = window

	// Restore window state if available
	if wm.stateStore != nil {
		if state, err := wm.stateStore.LoadState(windowID); err == nil {
			window.SetState(state)
		}
	}

	wm.log.Debug("Window created", "type", windowType, "id", windowID)
	return window, nil
}

// DestroyWindow properly destroys a window and cleans up resources
func (wm *WindowManager) DestroyWindow(window Window) error {
	wm.windowsMu.Lock()
	defer wm.windowsMu.Unlock()

	// Find window ID
	var windowID string
	for id, w := range wm.windows {
		if w == window {
			windowID = id
			break
		}
	}

	if windowID == "" {
		return fmt.Errorf("window not found in registry")
	}

	// Save state before destroying
	if wm.stateStore != nil && window.IsValid() {
		state := window.GetState()
		if err := wm.stateStore.SaveState(windowID, state); err != nil {
			wm.log.Warn("Failed to save window state", "error", err)
		}
	}

	// Notify lifecycle handlers
	wm.notifyLifecycleHandlers(window.GetType(), func(handler WindowLifecycle) {
		handler.OnClose()
	})

	// Close window
	window.Close()

	// Remove from registry
	delete(wm.windows, windowID)

	// Clear active window if it's the one being destroyed
	wm.activeMu.Lock()
	if wm.activeWindow == window {
		wm.activeWindow = nil
	}
	wm.activeMu.Unlock()

	wm.log.Debug("Window destroyed", "id", windowID)
	return nil
}

// GetActiveWindow returns the currently active window
func (wm *WindowManager) GetActiveWindow() Window {
	wm.activeMu.RLock()
	defer wm.activeMu.RUnlock()
	return wm.activeWindow
}

// SetActiveWindow sets the active window and notifies lifecycle handlers
func (wm *WindowManager) SetActiveWindow(window Window) error {
	wm.activeMu.Lock()
	defer wm.activeMu.Unlock()

	// Notify previous active window
	if wm.activeWindow != nil && wm.activeWindow != window {
		wm.notifyLifecycleHandlers(wm.activeWindow.GetType(), func(handler WindowLifecycle) {
			handler.OnBlur()
		})
	}

	wm.activeWindow = window

	// Notify new active window
	if window != nil {
		wm.notifyLifecycleHandlers(window.GetType(), func(handler WindowLifecycle) {
			handler.OnFocus()
		})
	}

	return nil
}

// GetWindowByType returns all windows of a specific type
func (wm *WindowManager) GetWindowByType(windowType WindowType) []Window {
	wm.windowsMu.RLock()
	defer wm.windowsMu.RUnlock()

	var windows []Window
	for _, window := range wm.windows {
		if window.GetType() == windowType {
			windows = append(windows, window)
		}
	}
	return windows
}

// CloseAllWindows closes all managed windows
func (wm *WindowManager) CloseAllWindows() error {
	wm.windowsMu.Lock()
	defer wm.windowsMu.Unlock()

	var errors []error
	for windowID, window := range wm.windows {
		if err := wm.DestroyWindow(window); err != nil {
			errors = append(errors, fmt.Errorf("failed to destroy window %s: %w", windowID, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing windows: %v", errors)
	}

	return nil
}

// AddLifecycleHandler adds a lifecycle handler for a window type
func (wm *WindowManager) AddLifecycleHandler(windowType WindowType, handler WindowLifecycle) {
	wm.lifecycleMu.Lock()
	defer wm.lifecycleMu.Unlock()

	wm.lifecycleHandlers[windowType] = append(wm.lifecycleHandlers[windowType], handler)
}

// RemoveLifecycleHandler removes a lifecycle handler
func (wm *WindowManager) RemoveLifecycleHandler(windowType WindowType, handler WindowLifecycle) {
	wm.lifecycleMu.Lock()
	defer wm.lifecycleMu.Unlock()

	handlers := wm.lifecycleHandlers[windowType]
	for i, h := range handlers {
		if h == handler {
			wm.lifecycleHandlers[windowType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// notifyLifecycleHandlers notifies all lifecycle handlers for a window type
func (wm *WindowManager) notifyLifecycleHandlers(windowType WindowType, action func(WindowLifecycle)) {
	wm.lifecycleMu.RLock()
	handlers := wm.lifecycleHandlers[windowType]
	wm.lifecycleMu.RUnlock()

	for _, handler := range handlers {
		action(handler)
	}
}

// createMainWindow creates a main window
func (wm *WindowManager) createMainWindow(windowID string, config WindowConfig) (Window, error) {
	fyneWindow := wm.app.NewWindow(config.Title)
	if fyneWindow == nil {
		return nil, fmt.Errorf("failed to create Fyne window")
	}
	if config.Width > 0 && config.Height > 0 {
		fyneWindow.Resize(fyne.NewSize(float32(config.Width), float32(config.Height)))
	}
	window := NewWindowWrapper(windowID, WindowTypeMain, fyneWindow, wm.log)
	initialState := WindowState{
		Width:   config.Width,
		Height:  config.Height,
		Visible: !config.StartHidden,
	}
	window.SetState(initialState)
	return window, nil
}

func (wm *WindowManager) createQuickNoteWindow(windowID string, config WindowConfig) (Window, error) {
	fyneWindow := wm.app.NewWindow(config.Title)
	if fyneWindow == nil {
		return nil, fmt.Errorf("failed to create Fyne window")
	}
	if config.Width > 0 && config.Height > 0 {
		fyneWindow.Resize(fyne.NewSize(float32(config.Width), float32(config.Height)))
	}
	window := NewWindowWrapper(windowID, WindowTypeQuickNote, fyneWindow, wm.log)
	initialState := WindowState{
		Width:   config.Width,
		Height:  config.Height,
		Visible: !config.StartHidden,
	}
	window.SetState(initialState)
	return window, nil
}

func (wm *WindowManager) createDialogWindow(windowID string, config WindowConfig) (Window, error) {
	fyneWindow := wm.app.NewWindow(config.Title)
	if fyneWindow == nil {
		return nil, fmt.Errorf("failed to create Fyne window")
	}
	if config.Width > 0 && config.Height > 0 {
		fyneWindow.Resize(fyne.NewSize(float32(config.Width), float32(config.Height)))
	}
	window := NewWindowWrapper(windowID, WindowTypeDialog, fyneWindow, wm.log)
	initialState := WindowState{
		Width:   config.Width,
		Height:  config.Height,
		Visible: !config.StartHidden,
	}
	window.SetState(initialState)
	return window, nil
}

// WindowConfig represents configuration for window creation
type WindowConfig struct {
	Title       string
	Width       int
	Height      int
	StartHidden bool
	Modal       bool
	Resizable   bool
	MinSize     fyne.Size
	MaxSize     fyne.Size
}
