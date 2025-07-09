//go:build linux || windows
// +build linux windows

package hotkey

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/csturiale/hotkey"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/platform"
)

// UnifiedManager manages hotkeys for both Linux and Windows
type UnifiedManager struct {
	log       logger.Logger
	quickNote QuickNoteService
	binding   *config.HotkeyBinding
	hk        *hotkey.Hotkey
	stopChan  chan struct{}
	running   bool
	mu        sync.Mutex
	config    *config.HotkeyConfig
}

// NewUnifiedManager creates a new UnifiedManager instance
func NewUnifiedManager(log logger.Logger, hotkeyConfig *config.HotkeyConfig) (*UnifiedManager, error) {
	log.Debug("Creating unified hotkey manager",
		"os", runtime.GOOS,
		"arch", runtime.GOARCH,
		"pid", os.Getpid(),
		"display", os.Getenv("DISPLAY"))
	return &UnifiedManager{
		log:      log,
		stopChan: make(chan struct{}),
		config:   hotkeyConfig,
	}, nil
}

// newPlatformManager creates a platform-specific manager
func newPlatformManager(quickNote QuickNoteService, binding *config.HotkeyBinding) Manager {
	manager := &UnifiedManager{
		log:       logger.NewNoopLogger(),
		quickNote: quickNote,
		binding:   binding,
		stopChan:  make(chan struct{}),
	}
	return manager
}

// SetLogger sets the logger for this manager
func (m *UnifiedManager) SetLogger(log logger.Logger) {
	m.log = log
}

// SetQuickNote sets the quick note service and hotkey binding
func (m *UnifiedManager) SetQuickNote(quickNote QuickNoteService, binding *config.HotkeyBinding) {
	m.log.Debug("Setting quick note service and binding",
		"binding", fmt.Sprintf("%+v", binding),
		"quicknote_nil", quickNote == nil)
	m.quickNote = quickNote
	m.binding = binding
}

// SetHotkey sets the hotkey instance (used for testing)
func (m *UnifiedManager) SetHotkey(hk *hotkey.Hotkey) {
	m.hk = hk
}

// checkPlatformSpecificRequirements performs platform-specific checks
func (m *UnifiedManager) checkPlatformSpecificRequirements() error {
	if runtime.GOOS == "linux" {
		// Linux-specific checks
		if platform.IsWSL2() {
			return ErrWSL2NotSupported
		}
		if platform.IsHeadless() {
			return ErrX11NotAvailable
		}
	}
	return nil
}

// Register registers the configured hotkey with the system
func (m *UnifiedManager) Register() error {
	if m.binding == nil {
		m.log.Error("Hotkey binding not set")
		return fmt.Errorf("hotkey binding not set")
	}

	if m.hk != nil {
		m.log.Error("Hotkey already registered")
		return fmt.Errorf("hotkey already registered")
	}

	// Platform-specific checks
	if err := m.checkPlatformSpecificRequirements(); err != nil {
		return err
	}

	m.log.Debug("Starting hotkey registration",
		"modifiers", strings.Join(m.binding.Modifiers, "+"),
		"key", m.binding.Key,
		"os", runtime.GOOS,
		"arch", runtime.GOARCH,
		"pid", os.Getpid(),
		"display", os.Getenv("DISPLAY"))

	mods, modErr := m.convertModifiers()
	if modErr != nil {
		return modErr
	}

	key, keyErr := m.convertKey()
	if keyErr != nil {
		return keyErr
	}

	// Create hotkey directly
	m.hk = hotkey.New(mods, key)

	if regErr := m.registerWithRetries(); regErr != nil {
		return regErr
	}

	m.log.Debug("Successfully registered hotkey",
		"modifiers", strings.Join(m.binding.Modifiers, "+"),
		"key", m.binding.Key,
		"os", runtime.GOOS,
		"arch", runtime.GOARCH,
		"pid", os.Getpid(),
		"display", os.Getenv("DISPLAY"))

	return nil
}

func (m *UnifiedManager) convertModifiers() ([]hotkey.Modifier, error) {
	var mods []hotkey.Modifier
	m.log.Debug("Converting modifiers", "raw_modifiers", m.binding.Modifiers)
	for _, mod := range m.binding.Modifiers {
		switch mod {
		case "Ctrl":
			m.log.Debug("Adding Ctrl modifier")
			mods = append(mods, hotkey.ModCtrl)
		case "Shift":
			m.log.Debug("Adding Shift modifier")
			mods = append(mods, hotkey.ModShift)
		case "Alt":
			m.log.Debug("Adding Alt modifier")
			// Platform-specific Alt modifier
			switch runtime.GOOS {
			case "windows":
				mods = append(mods, hotkey.Modifier(0x1))
			case "linux", "darwin":
				mods = append(mods, hotkey.Modifier(8))
			default:
				m.log.Error("Unsupported platform for Alt modifier", "platform", runtime.GOOS)
				return nil, fmt.Errorf("unsupported platform for Alt modifier: %s", runtime.GOOS)
			}
		default:
			m.log.Error("Unknown modifier", "modifier", mod)
			return nil, fmt.Errorf("unknown modifier: %s", mod)
		}
	}
	m.log.Debug("Converted modifiers", "count", len(mods))
	return mods, nil
}

func (m *UnifiedManager) convertKey() (hotkey.Key, error) {
	m.log.Debug("Converting key", "raw_key", m.binding.Key)

	// Handle single character keys (A-Z, 0-9)
	if len(m.binding.Key) == 1 {
		char := m.binding.Key[0]
		if char >= 'A' && char <= 'Z' {
			switch char {
			case 'A':
				return hotkey.KeyA, nil
			case 'B':
				return hotkey.KeyB, nil
			case 'C':
				return hotkey.KeyC, nil
			case 'D':
				return hotkey.KeyD, nil
			case 'E':
				return hotkey.KeyE, nil
			case 'F':
				return hotkey.KeyF, nil
			case 'G':
				return hotkey.KeyG, nil
			case 'H':
				return hotkey.KeyH, nil
			case 'I':
				return hotkey.KeyI, nil
			case 'J':
				return hotkey.KeyJ, nil
			case 'K':
				return hotkey.KeyK, nil
			case 'L':
				return hotkey.KeyL, nil
			case 'M':
				return hotkey.KeyM, nil
			case 'N':
				return hotkey.KeyN, nil
			case 'O':
				return hotkey.KeyO, nil
			case 'P':
				return hotkey.KeyP, nil
			case 'Q':
				return hotkey.KeyQ, nil
			case 'R':
				return hotkey.KeyR, nil
			case 'S':
				return hotkey.KeyS, nil
			case 'T':
				return hotkey.KeyT, nil
			case 'U':
				return hotkey.KeyU, nil
			case 'V':
				return hotkey.KeyV, nil
			case 'W':
				return hotkey.KeyW, nil
			case 'X':
				return hotkey.KeyX, nil
			case 'Y':
				return hotkey.KeyY, nil
			case 'Z':
				return hotkey.KeyZ, nil
			}
		}
		if char >= '0' && char <= '9' {
			switch char {
			case '0':
				return hotkey.Key0, nil
			case '1':
				return hotkey.Key1, nil
			case '2':
				return hotkey.Key2, nil
			case '3':
				return hotkey.Key3, nil
			case '4':
				return hotkey.Key4, nil
			case '5':
				return hotkey.Key5, nil
			case '6':
				return hotkey.Key6, nil
			case '7':
				return hotkey.Key7, nil
			case '8':
				return hotkey.Key8, nil
			case '9':
				return hotkey.Key9, nil
			}
		}
	}

	// Handle function keys
	if strings.HasPrefix(m.binding.Key, "F") && len(m.binding.Key) <= 3 {
		switch m.binding.Key {
		case "F1":
			return hotkey.KeyF1, nil
		case "F2":
			return hotkey.KeyF2, nil
		case "F3":
			return hotkey.KeyF3, nil
		case "F4":
			return hotkey.KeyF4, nil
		case "F5":
			return hotkey.KeyF5, nil
		case "F6":
			return hotkey.KeyF6, nil
		case "F7":
			return hotkey.KeyF7, nil
		case "F8":
			return hotkey.KeyF8, nil
		case "F9":
			return hotkey.KeyF9, nil
		case "F10":
			return hotkey.KeyF10, nil
		case "F11":
			return hotkey.KeyF11, nil
		case "F12":
			return hotkey.KeyF12, nil
		}
	}

	// Handle special keys
	switch m.binding.Key {
	case "Space":
		return hotkey.KeySpace, nil
	case "Return", "Enter":
		return hotkey.KeyReturn, nil
	case "Escape", "Esc":
		return hotkey.KeyEscape, nil
	case "Delete", "Del":
		return hotkey.KeyDelete, nil
	case "Tab":
		return hotkey.KeyTab, nil
	case "Left":
		return hotkey.KeyLeft, nil
	case "Right":
		return hotkey.KeyRight, nil
	case "Up":
		return hotkey.KeyUp, nil
	case "Down":
		return hotkey.KeyDown, nil
	}

	m.log.Error("Unsupported key", "key", m.binding.Key,
		"supported_keys", "A-Z, 0-9, F1-F12, Space, Return, Escape, Delete, Tab, Arrow keys")
	return 0, fmt.Errorf("unsupported key: %s", m.binding.Key)
}

func (m *UnifiedManager) registerWithRetries() error {
	var err error
	delay := time.Duration(m.config.RetryDelayMs) * time.Millisecond
	retries := m.config.MaxRetries
	if delay == 0 {
		delay = 100 * time.Millisecond
	}
	if retries == 0 {
		retries = 3
	}

	for i := 0; i < retries; i++ {
		m.log.Debug("Attempting to register hotkey with system", "attempt", i+1)
		if err = m.hk.Register(); err == nil {
			break
		}
		m.log.Error("Failed to register hotkey",
			"error", err,
			"attempt", i+1)
		time.Sleep(delay)
	}

	if err != nil {
		return fmt.Errorf("failed to register hotkey after %d attempts: %w", retries, err)
	}

	return nil
}

// Unregister removes the hotkey registration
func (m *UnifiedManager) Unregister() error {
	m.log.Info("Unregistering hotkey")

	if m.hk != nil {
		m.log.Info("Hotkey instance exists, attempting to unregister")
		if err := m.hk.Unregister(); err != nil {
			m.log.Error("Failed to unregister hotkey", "error", err)
			return fmt.Errorf("failed to unregister hotkey: %w", err)
		}
		m.log.Info("Successfully unregistered hotkey")
		m.hk = nil
	} else {
		m.log.Info("No hotkey instance to unregister")
	}
	return nil
}

// Start begins listening for hotkey events
func (m *UnifiedManager) Start() error {
	m.log.Debug("Starting hotkey manager",
		"hotkey_nil", m.hk == nil,
		"quicknote_nil", m.quickNote == nil)

	if m.hk == nil {
		m.log.Error("Hotkey not registered")
		return fmt.Errorf("hotkey not registered")
	}

	if m.quickNote == nil {
		m.log.Error("Quick note service not set")
		return ErrQuickNoteNotSet
	}

	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		m.log.Warn("Hotkey manager already running")
		return nil
	}
	m.running = true
	m.mu.Unlock()

	m.log.Debug("Hotkey manager started successfully")

	// Start listening for hotkey events in a goroutine
	go func() {
		defer func() {
			m.mu.Lock()
			m.running = false
			m.mu.Unlock()
			m.log.Info("Hotkey manager stopped")
		}()

		for {
			select {
			case <-m.stopChan:
				m.log.Info("Received stop signal, ending hotkey listening")
				return
			case event := <-m.hk.Keydown():
				m.log.Info("Hotkey triggered", "event", event)
				if m.quickNote != nil {
					m.quickNote.Show()
				}
			}
		}
	}()

	return nil
}

// Stop ends the hotkey listening and unregisters the hotkey
func (m *UnifiedManager) Stop() error {
	m.log.Info("Stopping hotkey manager")

	m.mu.Lock()
	if !m.running {
		m.mu.Unlock()
		m.log.Info("Hotkey manager not running")
		return nil
	}
	m.mu.Unlock()

	// Signal the goroutine to stop
	close(m.stopChan)

	// Unregister the hotkey
	if err := m.Unregister(); err != nil {
		m.log.Error("Failed to unregister hotkey during stop", "error", err)
		return err
	}

	m.log.Info("Hotkey manager stopped successfully")
	return nil
}
