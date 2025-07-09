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

const (
	// retryDelay is the delay between retry attempts
	retryDelay = 100 * time.Millisecond
	// maxRetries is the maximum number of retry attempts
	maxRetries = 3
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
}

// NewUnifiedManager creates a new UnifiedManager instance
func NewUnifiedManager(log logger.Logger) (*UnifiedManager, error) {
	log.Debug("Creating unified hotkey manager",
		"os", runtime.GOOS,
		"arch", runtime.GOARCH,
		"pid", os.Getpid(),
		"display", os.Getenv("DISPLAY"))
	return &UnifiedManager{
		log:      log,
		stopChan: make(chan struct{}),
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
	switch m.binding.Key {
	case "G":
		m.log.Debug("Using key G", "key_code", hotkey.KeyG)
		return hotkey.KeyG, nil
	case "N":
		m.log.Debug("Using key N", "key_code", hotkey.KeyN)
		return hotkey.KeyN, nil
	default:
		m.log.Error("Unsupported key", "key", m.binding.Key,
			"supported_keys", []string{"G", "N"})
		return 0, fmt.Errorf("unsupported key: %s", m.binding.Key)
	}
}

func (m *UnifiedManager) registerWithRetries() error {
	var err error
	for i := 0; i < maxRetries; i++ {
		m.log.Debug("Attempting to register hotkey with system", "attempt", i+1)
		if err = m.hk.Register(); err == nil {
			break
		}
		m.log.Error("Failed to register hotkey",
			"error", err,
			"attempt", i+1)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return fmt.Errorf("failed to register hotkey after %d attempts: %w", maxRetries, err)
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
