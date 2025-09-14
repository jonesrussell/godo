//go:build linux || windows
// +build linux windows

package hotkey

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/csturiale/hotkey"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/platform"
)

// Manager defines the interface for hotkey management
type Manager interface {
	// Register registers all configured hotkeys with the system
	Register() error
	// Unregister removes all hotkey registrations from the system
	Unregister() error
	// Start begins listening for hotkey events
	Start() error
	// Stop ends the hotkey listening and unregisters all hotkeys
	Stop() error
	// SetQuickNote configures the quick note service and hotkey binding
	SetQuickNote(quickNote QuickNoteService, binding *config.HotkeyBinding)
	// SetQuickNoteFactory configures a factory function to create the quick note service on demand
	SetQuickNoteFactory(factory func() QuickNoteService, binding *config.HotkeyBinding)
	// SetMainWindow configures the main window service and hotkey binding
	SetMainWindow(mainWindow MainWindowService, binding *config.HotkeyBinding)
	// SetMainWindowFactory configures a factory function to create the main window service on demand
	SetMainWindowFactory(factory func() MainWindowService, binding *config.HotkeyBinding)
}

// QuickNoteService defines quick note operations that can be triggered by hotkeys
type QuickNoteService interface {
	Show()
	Hide()
}

// MainWindowService defines main window operations that can be triggered by hotkeys
type MainWindowService interface {
	Show()
	Hide()
}

// HotkeyEntry represents a single hotkey configuration
type HotkeyEntry struct {
	hotkey     *hotkey.Hotkey
	binding    *config.HotkeyBinding
	quickNote  QuickNoteService
	mainWindow MainWindowService
}

// HotkeyManager manages hotkeys using the simplified library approach
type HotkeyManager struct {
	log      logger.Logger
	hotkeys  []HotkeyEntry
	stopChan chan struct{}
	running  bool
}

// NewManager creates a new HotkeyManager instance
func NewManager(log logger.Logger, hotkeyConfig *config.HotkeyConfig) (Manager, error) {
	return &HotkeyManager{
		log:      log,
		stopChan: make(chan struct{}),
	}, nil
}

// SetQuickNote configures the quick note service and hotkey binding
func (m *HotkeyManager) SetQuickNote(quickNote QuickNoteService, binding *config.HotkeyBinding) {
	m.log.Debug("Setting quick note service and binding",
		"binding", fmt.Sprintf("%+v", binding),
		"quicknote_nil", quickNote == nil)

	// Create hotkey entry for quick note
	entry := HotkeyEntry{
		binding:   binding,
		quickNote: quickNote,
	}
	m.hotkeys = append(m.hotkeys, entry)
}

// SetQuickNoteFactory configures a factory function to create the quick note service on demand
func (m *HotkeyManager) SetQuickNoteFactory(factory func() QuickNoteService, binding *config.HotkeyBinding) {
	m.log.Debug("Setting quick note factory and binding",
		"binding", fmt.Sprintf("%+v", binding),
		"factory_nil", factory == nil)

	// Create hotkey entry for quick note with factory
	entry := HotkeyEntry{
		binding:   binding,
		quickNote: factory(), // Call factory immediately to get service
	}
	m.hotkeys = append(m.hotkeys, entry)
}

// SetMainWindow configures the main window service and hotkey binding
func (m *HotkeyManager) SetMainWindow(mainWindow MainWindowService, binding *config.HotkeyBinding) {
	m.log.Debug("Setting main window service and binding",
		"binding", fmt.Sprintf("%+v", binding),
		"mainwindow_nil", mainWindow == nil)

	// Create hotkey entry for main window
	entry := HotkeyEntry{
		binding:    binding,
		mainWindow: mainWindow,
	}
	m.hotkeys = append(m.hotkeys, entry)
}

// SetMainWindowFactory configures a factory function to create the main window service on demand
func (m *HotkeyManager) SetMainWindowFactory(factory func() MainWindowService, binding *config.HotkeyBinding) {
	m.log.Debug("Setting main window factory and binding",
		"binding", fmt.Sprintf("%+v", binding),
		"factory_nil", factory == nil)

	// Create hotkey entry for main window with factory
	entry := HotkeyEntry{
		binding:    binding,
		mainWindow: factory(), // Call factory immediately to get service
	}
	m.hotkeys = append(m.hotkeys, entry)
}

// Register registers all configured hotkeys with the system
func (m *HotkeyManager) Register() error {
	if len(m.hotkeys) == 0 {
		return fmt.Errorf("no hotkeys configured")
	}

	// Platform-specific checks
	if err := m.checkPlatformRequirements(); err != nil {
		return err
	}

	for i, entry := range m.hotkeys {
		if entry.binding == nil {
			return fmt.Errorf("hotkey binding not set for entry %d", i)
		}

		m.log.Info("Registering hotkey",
			"entry", i,
			"modifiers", strings.Join(entry.binding.Modifiers, "+"),
			"key", entry.binding.Key,
			"os", runtime.GOOS)

		mods, err := m.convertModifiers(entry.binding)
		if err != nil {
			return fmt.Errorf("failed to convert modifiers for entry %d: %w", i, err)
		}

		key, err := m.convertKey(entry.binding)
		if err != nil {
			return fmt.Errorf("failed to convert key for entry %d: %w", i, err)
		}

		// Create and register hotkey using library's simple API
		hk := hotkey.New(mods, key)
		if registerErr := hk.Register(); registerErr != nil {
			return fmt.Errorf("failed to register hotkey for entry %d: %w", i, registerErr)
		}

		// Store the hotkey in the entry
		m.hotkeys[i].hotkey = hk

		m.log.Info("Successfully registered hotkey",
			"entry", i,
			"modifiers", strings.Join(entry.binding.Modifiers, "+"),
			"key", entry.binding.Key)
	}

	return nil
}

// Unregister removes all hotkey registrations
func (m *HotkeyManager) Unregister() error {
	var lastErr error

	for i, entry := range m.hotkeys {
		if entry.hotkey != nil {
			if err := entry.hotkey.Unregister(); err != nil {
				m.log.Error("Failed to unregister hotkey", "entry", i, "error", err)
				lastErr = fmt.Errorf("failed to unregister hotkey %d: %w", i, err)
			} else {
				m.log.Info("Hotkey unregistered", "entry", i)
			}
			m.hotkeys[i].hotkey = nil
		}
	}

	return lastErr
}

// Start begins listening for hotkey events using the library's channel-based approach
func (m *HotkeyManager) Start() error {
	if len(m.hotkeys) == 0 {
		return fmt.Errorf("no hotkeys registered")
	}

	if m.running {
		return nil
	}

	m.running = true
	m.log.Info("Starting hotkey listener")

	// Use the library's recommended channel-based event handling
	go func() {
		defer func() {
			m.running = false
			m.log.Info("Hotkey listener stopped")
		}()

		// Create a slice of channels for all hotkeys
		channels := make([]<-chan hotkey.Event, len(m.hotkeys))
		for i, entry := range m.hotkeys {
			if entry.hotkey != nil {
				channels[i] = entry.hotkey.Keydown()
			}
		}

		for {
			select {
			case <-m.stopChan:
				return
			default:
				// Check each hotkey channel
				for i, ch := range channels {
					if ch != nil {
						select {
						case event := <-ch:
							m.log.Info("Hotkey triggered", "entry", i, "event", event)
							entry := m.hotkeys[i]

							// Trigger the appropriate service
							if entry.quickNote != nil {
								entry.quickNote.Show()
							} else if entry.mainWindow != nil {
								entry.mainWindow.Show()
							}
						default:
							// No event on this channel, continue
						}
					}
				}
			}
		}
	}()

	return nil
}

// Stop ends the hotkey listening
func (m *HotkeyManager) Stop() error {
	if !m.running {
		return nil
	}

	close(m.stopChan)
	m.running = false

	if err := m.Unregister(); err != nil {
		return fmt.Errorf("failed to unregister during stop: %w", err)
	}

	return nil
}

// checkPlatformRequirements performs basic platform checks
func (m *HotkeyManager) checkPlatformRequirements() error {
	if runtime.GOOS == "linux" {
		if platform.IsWSL2() {
			return ErrWSL2NotSupported
		}
		if platform.IsHeadless() {
			return ErrX11NotAvailable
		}
	}
	return nil
}

// convertModifiers converts string modifiers to hotkey.Modifier using a simple map
func (m *HotkeyManager) convertModifiers(binding *config.HotkeyBinding) ([]hotkey.Modifier, error) {
	modifierMap := map[string]hotkey.Modifier{
		"Ctrl":  hotkey.ModCtrl,
		"Shift": hotkey.ModShift,
		"Alt":   m.getAltModifier(),
	}

	var mods []hotkey.Modifier
	for _, modStr := range binding.Modifiers {
		if mod, exists := modifierMap[modStr]; exists {
			mods = append(mods, mod)
		} else {
			return nil, fmt.Errorf("unknown modifier: %s", modStr)
		}
	}

	return mods, nil
}

// getAltModifier returns the platform-specific Alt modifier
func (m *HotkeyManager) getAltModifier() hotkey.Modifier {
	switch runtime.GOOS {
	case "windows":
		return hotkey.Modifier(0x1)
	case "linux", "darwin":
		return hotkey.Modifier(8)
	default:
		// Fallback to a reasonable default
		return hotkey.Modifier(8)
	}
}

// convertKey converts string key to hotkey.Key using a simple map
func (m *HotkeyManager) convertKey(binding *config.HotkeyBinding) (hotkey.Key, error) {
	// Simple map for common keys - only include keys that actually exist in the library
	keyMap := map[string]hotkey.Key{
		// Letters
		"A": hotkey.KeyA, "B": hotkey.KeyB, "C": hotkey.KeyC, "D": hotkey.KeyD,
		"E": hotkey.KeyE, "F": hotkey.KeyF, "G": hotkey.KeyG, "H": hotkey.KeyH,
		"I": hotkey.KeyI, "J": hotkey.KeyJ, "K": hotkey.KeyK, "L": hotkey.KeyL,
		"M": hotkey.KeyM, "N": hotkey.KeyN, "O": hotkey.KeyO, "P": hotkey.KeyP,
		"Q": hotkey.KeyQ, "R": hotkey.KeyR, "S": hotkey.KeyS, "T": hotkey.KeyT,
		"U": hotkey.KeyU, "V": hotkey.KeyV, "W": hotkey.KeyW, "X": hotkey.KeyX,
		"Y": hotkey.KeyY, "Z": hotkey.KeyZ,
		// Numbers
		"0": hotkey.Key0, "1": hotkey.Key1, "2": hotkey.Key2, "3": hotkey.Key3,
		"4": hotkey.Key4, "5": hotkey.Key5, "6": hotkey.Key6, "7": hotkey.Key7,
		"8": hotkey.Key8, "9": hotkey.Key9,
		// Function keys
		"F1": hotkey.KeyF1, "F2": hotkey.KeyF2, "F3": hotkey.KeyF3, "F4": hotkey.KeyF4,
		"F5": hotkey.KeyF5, "F6": hotkey.KeyF6, "F7": hotkey.KeyF7, "F8": hotkey.KeyF8,
		"F9": hotkey.KeyF9, "F10": hotkey.KeyF10, "F11": hotkey.KeyF11, "F12": hotkey.KeyF12,
		// Common keys
		"Space": hotkey.KeySpace, "Enter": hotkey.KeyReturn, "Tab": hotkey.KeyTab,
		"Escape": hotkey.KeyEscape,
		"Up":     hotkey.KeyUp, "Down": hotkey.KeyDown, "Left": hotkey.KeyLeft, "Right": hotkey.KeyRight,
	}

	if key, exists := keyMap[binding.Key]; exists {
		return key, nil
	}

	return 0, fmt.Errorf("unsupported key: %s", binding.Key)
}
