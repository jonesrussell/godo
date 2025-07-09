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
	// Register registers the hotkey with the system
	Register() error
	// Unregister removes the hotkey registration from the system
	Unregister() error
	// Start begins listening for hotkey events
	Start() error
	// Stop ends the hotkey listening and unregisters the hotkey
	Stop() error
	// SetQuickNote configures the quick note service and hotkey binding
	SetQuickNote(quickNote QuickNoteService, binding *config.HotkeyBinding)
	// SetQuickNoteFactory configures a factory function to create the quick note service on demand
	SetQuickNoteFactory(factory func() QuickNoteService, binding *config.HotkeyBinding)
}

// QuickNoteService defines quick note operations that can be triggered by hotkeys
type QuickNoteService interface {
	Show()
	Hide()
}

// HotkeyManager manages hotkeys using the simplified library approach
type HotkeyManager struct {
	log           logger.Logger
	quickNote     QuickNoteService
	quickNoteFunc func() QuickNoteService // Factory function for backward compatibility
	binding       *config.HotkeyBinding
	hk            *hotkey.Hotkey
	stopChan      chan struct{}
	running       bool
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
	m.quickNote = quickNote
	m.binding = binding
}

// SetQuickNoteFactory configures a factory function to create the quick note service on demand
func (m *HotkeyManager) SetQuickNoteFactory(factory func() QuickNoteService, binding *config.HotkeyBinding) {
	m.log.Debug("Setting quick note factory and binding",
		"binding", fmt.Sprintf("%+v", binding),
		"factory_nil", factory == nil)
	m.quickNoteFunc = factory
	m.binding = binding
}

// Register registers the configured hotkey with the system
func (m *HotkeyManager) Register() error {
	if m.binding == nil {
		return fmt.Errorf("hotkey binding not set")
	}

	if m.hk != nil {
		return fmt.Errorf("hotkey already registered")
	}

	// Platform-specific checks
	if err := m.checkPlatformRequirements(); err != nil {
		return err
	}

	m.log.Info("Registering hotkey",
		"modifiers", strings.Join(m.binding.Modifiers, "+"),
		"key", m.binding.Key,
		"os", runtime.GOOS)

	mods, err := m.convertModifiers()
	if err != nil {
		return fmt.Errorf("failed to convert modifiers: %w", err)
	}

	key, err := m.convertKey()
	if err != nil {
		return fmt.Errorf("failed to convert key: %w", err)
	}

	// Create and register hotkey using library's simple API
	m.hk = hotkey.New(mods, key)
	if registerErr := m.hk.Register(); registerErr != nil {
		return fmt.Errorf("failed to register hotkey: %w", registerErr)
	}

	m.log.Info("Successfully registered hotkey",
		"modifiers", strings.Join(m.binding.Modifiers, "+"),
		"key", m.binding.Key)
	return nil
}

// Unregister removes the hotkey registration
func (m *HotkeyManager) Unregister() error {
	if m.hk == nil {
		return nil
	}

	if err := m.hk.Unregister(); err != nil {
		return fmt.Errorf("failed to unregister hotkey: %w", err)
	}

	m.hk = nil
	m.log.Info("Hotkey unregistered")
	return nil
}

// Start begins listening for hotkey events using the library's channel-based approach
func (m *HotkeyManager) Start() error {
	if m.hk == nil {
		return fmt.Errorf("hotkey not registered")
	}

	if m.quickNote == nil && m.quickNoteFunc == nil {
		return fmt.Errorf("quick note service not set")
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

		for {
			select {
			case <-m.stopChan:
				return
			case <-m.hk.Keydown():
				m.log.Info("Hotkey triggered")
				// Get quick note service - either from existing instance or factory
				var quickNoteService QuickNoteService
				if m.quickNote != nil {
					quickNoteService = m.quickNote
				} else if m.quickNoteFunc != nil {
					quickNoteService = m.quickNoteFunc()
				}

				if quickNoteService != nil {
					quickNoteService.Show()
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
func (m *HotkeyManager) convertModifiers() ([]hotkey.Modifier, error) {
	modifierMap := map[string]hotkey.Modifier{
		"Ctrl":  hotkey.ModCtrl,
		"Shift": hotkey.ModShift,
		"Alt":   m.getAltModifier(),
	}

	var mods []hotkey.Modifier
	for _, modStr := range m.binding.Modifiers {
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
func (m *HotkeyManager) convertKey() (hotkey.Key, error) {
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

	if key, exists := keyMap[m.binding.Key]; exists {
		return key, nil
	}

	return 0, fmt.Errorf("unsupported key: %s", m.binding.Key)
}
