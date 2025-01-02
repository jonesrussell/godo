//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/logger"
	"golang.design/x/hotkey"
)

const (
	// retryDelay is the delay between retry attempts
	retryDelay = 100 * time.Millisecond
	// maxRetries is the maximum number of retry attempts
	maxRetries = 3
)

// WindowsManager manages hotkeys for Windows
type WindowsManager struct {
	log       logger.Logger
	quickNote gui.QuickNoteManager
	binding   *common.HotkeyBinding
	hk        hotkeyInterface
	stopChan  chan struct{}
}

// NewWindowsManager creates a new WindowsManager instance with the provided logger.
// It initializes the manager but does not register any hotkeys until Register is called.
func NewWindowsManager(log logger.Logger) (*WindowsManager, error) {
	log.Info("Creating Windows hotkey manager",
		"os", runtime.GOOS,
		"arch", runtime.GOARCH,
		"pid", os.Getpid())
	return &WindowsManager{
		log:      log,
		stopChan: make(chan struct{}),
	}, nil
}

// SetQuickNote sets the quick note service and hotkey binding for this manager.
// Both the service and binding are required for the hotkey to function.
func (m *WindowsManager) SetQuickNote(quickNote gui.QuickNoteManager, binding *common.HotkeyBinding) {
	m.log.Info("Setting quick note service and binding",
		"binding", fmt.Sprintf("%+v", binding),
		"quicknote_nil", quickNote == nil)
	m.quickNote = quickNote
	m.binding = binding
}

// SetHotkey sets the hotkey instance (used for testing)
func (m *WindowsManager) SetHotkey(hk hotkeyInterface) {
	m.hk = hk
}

// Register registers the configured hotkey with the Windows system.
// It will attempt to register the hotkey multiple times in case of failure.
// Returns an error if registration fails after all attempts.
func (m *WindowsManager) Register() error {
	if m.binding == nil {
		m.log.Error("Hotkey binding not set")
		return fmt.Errorf("hotkey binding not set")
	}

	// Check if already registered
	if m.hk != nil && m.hk.IsRegistered() {
		m.log.Error("Hotkey already registered")
		return fmt.Errorf("hotkey already registered")
	}

	m.log.Info("Starting hotkey registration",
		"modifiers", strings.Join(m.binding.Modifiers, "+"),
		"key", m.binding.Key,
		"os", runtime.GOOS,
		"pid", os.Getpid())

	// Convert string modifiers to hotkey.Modifier
	var mods []hotkey.Modifier
	m.log.Info("Converting modifiers", "raw_modifiers", m.binding.Modifiers)
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
			mods = append(mods, hotkey.ModAlt)
		default:
			m.log.Error("Unknown modifier", "modifier", mod)
			return fmt.Errorf("unknown modifier: %s", mod)
		}
	}
	m.log.Info("Converted modifiers", "count", len(mods))

	// Convert key string to hotkey.Key
	var key hotkey.Key
	m.log.Info("Converting key", "raw_key", m.binding.Key)
	switch m.binding.Key {
	case "G":
		key = hotkey.KeyG
		m.log.Info("Using key G", "key_code", hotkey.KeyG)
	case "N":
		key = hotkey.KeyN
		m.log.Info("Using key N", "key_code", hotkey.KeyN)
	default:
		m.log.Error("Unsupported key", "key", m.binding.Key,
			"supported_keys", []string{"G", "N"})
		return fmt.Errorf("unsupported key: %s", m.binding.Key)
	}

	// Create and register the hotkey
	m.log.Info("Creating hotkey instance",
		"modifiers_count", len(mods),
		"key", key)

	// If we don't have a hotkey instance yet, create one
	if m.hk == nil {
		m.hk = newHotkeyWrapper(mods, key)
	}

	// Try registration with retries
	var err error
	for i := 0; i < maxRetries; i++ {
		m.log.Info("Attempting to register hotkey with system", "attempt", i+1)
		if err = m.hk.Register(); err == nil {
			break
		}
		m.log.Error("Failed to register hotkey",
			"error", err,
			"attempt", i+1,
			"modifiers", mods,
			"key", key)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return fmt.Errorf("failed to register hotkey after %d attempts: %w", maxRetries, err)
	}

	m.log.Info("Successfully registered hotkey",
		"modifiers", strings.Join(m.binding.Modifiers, "+"),
		"key", m.binding.Key,
		"os", runtime.GOOS,
		"pid", os.Getpid())

	return nil
}

// Unregister removes the hotkey registration from the Windows system.
// It's safe to call this method multiple times, even if no hotkey is registered.
func (m *WindowsManager) Unregister() error {
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

// Start begins listening for hotkey events and shows the quick note window when triggered.
// Returns an error if the hotkey is not registered or the quick note service is not set.
func (m *WindowsManager) Start() error {
	m.log.Info("Starting hotkey manager",
		"hotkey_nil", m.hk == nil,
		"quicknote_nil", m.quickNote == nil)

	if m.hk == nil {
		m.log.Error("Cannot start - hotkey not registered")
		return fmt.Errorf("hotkey not registered")
	}

	if m.quickNote == nil {
		m.log.Error("Cannot start - quick note service not set")
		return fmt.Errorf("quick note service not set")
	}

	go func() {
		m.log.Info("Starting hotkey listener goroutine")
		for {
			select {
			case <-m.stopChan:
				m.log.Info("Hotkey manager received quit signal")
				return
			case <-m.hk.Keydown():
				m.log.Info("Hotkey triggered - showing quick note window")
				m.quickNote.Show()
			}
		}
	}()

	m.log.Info("Hotkey manager started successfully")
	return nil
}

// Stop ends the hotkey listener and unregisters the hotkey.
// It's safe to call this method multiple times.
func (m *WindowsManager) Stop() error {
	m.log.Info("Stopping hotkey manager")
	close(m.stopChan)
	return m.Unregister()
}
