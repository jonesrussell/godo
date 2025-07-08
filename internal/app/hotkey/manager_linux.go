//go:build linux
// +build linux

package hotkey

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.design/x/hotkey"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
)

const (
	// retryDelay is the delay between retry attempts
	retryDelay = 100 * time.Millisecond
	// maxRetries is the maximum number of retry attempts
	maxRetries = 3
)

// LinuxManager manages hotkeys for Linux
type LinuxManager struct {
	log       logger.Logger
	quickNote QuickNoteService
	binding   *common.HotkeyBinding
	hk        hotkeyInterface
	stopChan  chan struct{}
	running   bool
	mu        sync.Mutex
}

// NewLinuxManager creates a new LinuxManager instance with the provided logger.
// It initializes the manager but does not register any hotkeys until Register is called.
func NewLinuxManager(log logger.Logger) (*LinuxManager, error) {
	log.Info("Creating Linux hotkey manager",
		"os", runtime.GOOS,
		"arch", runtime.GOARCH,
		"pid", os.Getpid(),
		"display", os.Getenv("DISPLAY"))
	return &LinuxManager{
		log:      log,
		stopChan: make(chan struct{}),
	}, nil
}

// newPlatformManager creates a platform-specific manager for Linux
func newPlatformManager(quickNote QuickNoteService, binding *common.HotkeyBinding) Manager {
	manager := &LinuxManager{
		log:       logger.NewNoopLogger(), // Default logger, should be set via injection
		quickNote: quickNote,
		binding:   binding,
		stopChan:  make(chan struct{}),
	}
	return manager
}

// SetLogger sets the logger for this manager
func (m *LinuxManager) SetLogger(log logger.Logger) {
	m.log = log
}

// SetQuickNote sets the quick note service and hotkey binding for this manager.
// Both the service and binding are required for the hotkey to function.
func (m *LinuxManager) SetQuickNote(quickNote QuickNoteService, binding *common.HotkeyBinding) {
	m.log.Info("Setting quick note service and binding",
		"binding", fmt.Sprintf("%+v", binding),
		"quicknote_nil", quickNote == nil)
	m.quickNote = quickNote
	m.binding = binding
}

// SetHotkey sets the hotkey instance (used for testing)
func (m *LinuxManager) SetHotkey(hk hotkeyInterface) {
	m.hk = hk
}

// checkX11Availability checks if X11 server is available
func (m *LinuxManager) checkX11Availability() error {
	display := os.Getenv("DISPLAY")
	if display == "" {
		return fmt.Errorf("DISPLAY environment variable not set - X11 server not available")
	}

	// Additional check for WSL2 or headless environments
	if strings.Contains(display, "WSL") || display == ":0" {
		m.log.Warn("Running in WSL2 or headless environment - hotkeys may not work",
			"display", display)
	}

	return nil
}

// Register registers the configured hotkey with the Linux system.
// It will attempt to register the hotkey multiple times in case of failure.
// Returns an error if registration fails after all attempts or if X11 is unavailable.
func (m *LinuxManager) Register() error {
	if bindErr := m.validateBindingAndX11(); bindErr != nil {
		return bindErr
	}
	if regCheckErr := m.checkAlreadyRegistered(); regCheckErr != nil {
		return regCheckErr
	}

	m.log.Info("Starting hotkey registration",
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

	m.createHotkeyIfNeeded(mods, key)

	if regErr := m.registerWithRetries(mods, key); regErr != nil {
		return regErr
	}

	m.log.Info("Successfully registered hotkey",
		"modifiers", strings.Join(m.binding.Modifiers, "+"),
		"key", m.binding.Key,
		"os", runtime.GOOS,
		"arch", runtime.GOARCH,
		"pid", os.Getpid(),
		"display", os.Getenv("DISPLAY"))

	return nil
}

func (m *LinuxManager) validateBindingAndX11() error {
	if m.binding == nil {
		m.log.Error("Hotkey binding not set")
		return fmt.Errorf("hotkey binding not set")
	}
	if x11Err := m.checkX11Availability(); x11Err != nil {
		m.log.Error("X11 server not available", "error", x11Err)
		return fmt.Errorf("X11 server not available: %w", x11Err)
	}
	return nil
}

func (m *LinuxManager) checkAlreadyRegistered() error {
	if m.hk != nil && m.hk.IsRegistered() {
		m.log.Error("Hotkey already registered")
		return fmt.Errorf("hotkey already registered")
	}
	return nil
}

func (m *LinuxManager) convertModifiers() ([]hotkey.Modifier, error) {
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
			mods = append(mods, hotkey.Mod1)
		default:
			m.log.Error("Unknown modifier", "modifier", mod)
			return nil, fmt.Errorf("unknown modifier: %s", mod)
		}
	}
	m.log.Info("Converted modifiers", "count", len(mods))
	return mods, nil
}

func (m *LinuxManager) convertKey() (hotkey.Key, error) {
	m.log.Info("Converting key", "raw_key", m.binding.Key)
	switch m.binding.Key {
	case "G":
		m.log.Info("Using key G", "key_code", hotkey.KeyG)
		return hotkey.KeyG, nil
	case "N":
		m.log.Info("Using key N", "key_code", hotkey.KeyN)
		return hotkey.KeyN, nil
	default:
		m.log.Error("Unsupported key", "key", m.binding.Key,
			"supported_keys", []string{"G", "N"})
		return 0, fmt.Errorf("unsupported key: %s", m.binding.Key)
	}
}

func (m *LinuxManager) createHotkeyIfNeeded(mods []hotkey.Modifier, key hotkey.Key) {
	m.log.Info("Creating hotkey instance",
		"modifiers_count", len(mods),
		"key", key)
	if m.hk == nil {
		m.hk = newHotkeyWrapper(mods, key)
	}
}

func (m *LinuxManager) registerWithRetries(mods []hotkey.Modifier, key hotkey.Key) error {
	var regErr error
	for i := 0; i < maxRetries; i++ {
		m.log.Info("Attempting to register hotkey with system", "attempt", i+1)
		if regErr = m.hk.Register(); regErr == nil {
			return nil
		}
		m.log.Error("Failed to register hotkey",
			"error", regErr,
			"attempt", i+1,
			"modifiers", mods,
			"key", key)
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("failed to register hotkey after %d attempts: %w", maxRetries, regErr)
}

// Unregister removes the hotkey registration from the Linux system.
// It's safe to call this method multiple times, even if no hotkey is registered.
func (m *LinuxManager) Unregister() error {
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
func (m *LinuxManager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		m.log.Info("Hotkey manager already running")
		return nil
	}

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

	m.running = true

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
func (m *LinuxManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		m.log.Info("Hotkey manager already stopped")
		return nil
	}

	m.log.Info("Stopping hotkey manager")
	m.running = false
	close(m.stopChan)
	return m.Unregister()
}
