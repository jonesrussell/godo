//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/jonesrussell/godo/internal/common"
	"golang.design/x/hotkey"
)

const (
	// cleanupDelay is the delay to wait for hotkey cleanup
	cleanupDelay = 500 * time.Millisecond
)

type platformManager struct {
	hk        *hotkey.Hotkey
	quickNote QuickNoteService
	binding   *common.HotkeyBinding
}

func newPlatformManager(quickNote QuickNoteService, binding *common.HotkeyBinding) Manager {
	fmt.Printf("[DEBUG] Creating hotkey manager (OS: %s)\n", runtime.GOOS)
	if quickNote == nil {
		panic("quickNote service cannot be nil")
	}
	return &platformManager{
		quickNote: quickNote,
		binding:   binding,
	}
}

func (m *platformManager) Register() error {
	modStr := strings.Join(m.binding.Modifiers, "+")
	fmt.Printf("[DEBUG] Registering hotkey (%s+%s) on %s\n", modStr, m.binding.Key, runtime.GOOS)

	// Convert string modifiers to hotkey.Modifier
	var mods []hotkey.Modifier
	for _, mod := range m.binding.Modifiers {
		switch strings.ToLower(mod) {
		case "ctrl":
			mods = append(mods, hotkey.ModCtrl)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "alt":
			mods = append(mods, hotkey.ModAlt)
		}
	}

	// Convert key string to hotkey.Key
	var key hotkey.Key
	switch strings.ToUpper(m.binding.Key) {
	case "N":
		key = hotkey.KeyN
	case "G":
		key = hotkey.KeyG
	// Add more key mappings as needed
	default:
		return fmt.Errorf("unsupported key: %s", m.binding.Key)
	}

	// Create the hotkey
	fmt.Println("[DEBUG] Creating hotkey instance...")
	hk := hotkey.New(mods, key)

	// Try to unregister any existing hotkey first
	if m.hk != nil {
		fmt.Println("[DEBUG] Attempting to unregister existing hotkey...")
		if err := m.hk.Unregister(); err != nil {
			fmt.Printf("[WARN] Failed to unregister existing hotkey: %v\n", err)
		}
		m.hk = nil
		time.Sleep(cleanupDelay) // Increased delay for cleanup
	}

	// Try to register with retries and increasing delays
	var err error
	delays := []time.Duration{100 * time.Millisecond, 500 * time.Millisecond, 1 * time.Second}
	for i := 0; i < len(delays); i++ {
		fmt.Printf("[DEBUG] Attempting to register hotkey (attempt %d/%d)...\n", i+1, len(delays))

		// Try to unregister before each attempt
		if unregErr := hk.Unregister(); unregErr != nil {
			fmt.Printf("[DEBUG] Unregister before attempt returned: %v\n", unregErr)
		}

		err = hk.Register()
		if err == nil {
			fmt.Printf("[DEBUG] Successfully registered hotkey on attempt %d\n", i+1)
			break
		}

		fmt.Printf("[WARN] Failed to register hotkey (attempt %d/%d): %v\n", i+1, len(delays), err)
		if i < len(delays)-1 { // Don't sleep after the last attempt
			delay := delays[i]
			fmt.Printf("[DEBUG] Waiting %v before next attempt...\n", delay)
			time.Sleep(delay)
		}
	}

	if err != nil {
		fmt.Printf("[ERROR] All attempts to register hotkey failed: %v\n", err)
		return fmt.Errorf("failed to register hotkey after %d attempts: %w", len(delays), err)
	}

	m.hk = hk
	fmt.Println("[DEBUG] Hotkey registered successfully, starting listener...")

	// Start listening for hotkey in a goroutine
	go func() {
		fmt.Println("[DEBUG] Hotkey listener started")
		for range hk.Keydown() {
			fmt.Println("[DEBUG] Hotkey triggered!")
			if m.quickNote != nil {
				fmt.Println("[DEBUG] Showing quick note window...")
				m.quickNote.Show()
			} else {
				fmt.Println("[ERROR] QuickNote service is nil!")
			}
		}
		fmt.Println("[DEBUG] Hotkey listener stopped")
	}()

	fmt.Println("[DEBUG] Register function completed successfully")
	return nil
}

func (m *platformManager) Unregister() error {
	fmt.Println("[DEBUG] Unregistering hotkey...")
	if m.hk != nil {
		if err := m.hk.Unregister(); err != nil {
			fmt.Printf("[ERROR] Failed to unregister hotkey: %v\n", err)
			return fmt.Errorf("failed to unregister hotkey: %w", err)
		}
		fmt.Println("[DEBUG] Hotkey unregistered successfully")
	} else {
		fmt.Println("[DEBUG] No hotkey to unregister")
	}
	return nil
}

// cleanupExistingHotkeys attempts to clean up any stale hotkey registrations
func (m *platformManager) cleanupExistingHotkeys() {
	fmt.Println("[DEBUG] Cleaning up existing hotkeys...")

	// Create a temporary hotkey with our configuration
	var mods []hotkey.Modifier
	for _, mod := range m.binding.Modifiers {
		switch strings.ToLower(mod) {
		case "ctrl":
			mods = append(mods, hotkey.ModCtrl)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "alt":
			mods = append(mods, hotkey.ModAlt)
		}
	}

	var key hotkey.Key
	switch strings.ToUpper(m.binding.Key) {
	case "N":
		key = hotkey.KeyN
	case "G":
		key = hotkey.KeyG
	default:
		fmt.Printf("[WARN] Unsupported key for cleanup: %s\n", m.binding.Key)
		return
	}

	// Try to unregister the hotkey
	hk := hotkey.New(mods, key)
	if err := hk.Unregister(); err != nil {
		fmt.Printf("[DEBUG] Cleanup unregister returned: %v\n", err)
	}

	// Wait for cleanup
	time.Sleep(cleanupDelay)
}

func (m *platformManager) Start() error {
	fmt.Println("[DEBUG] Starting hotkey manager...")

	// Clean up any existing hotkeys first
	m.cleanupExistingHotkeys()

	if m.hk == nil {
		return fmt.Errorf("hotkey not registered")
	}

	// Start listening for hotkey in a goroutine
	go func() {
		fmt.Println("[DEBUG] Hotkey listener started")
		for range m.hk.Keydown() {
			fmt.Println("[DEBUG] Hotkey triggered!")
			if m.quickNote != nil {
				fmt.Println("[DEBUG] Showing quick note window...")
				m.quickNote.Show()
			} else {
				fmt.Println("[ERROR] QuickNote service is nil!")
			}
		}
		fmt.Println("[DEBUG] Hotkey listener stopped")
	}()

	return nil
}

func (m *platformManager) Stop() error {
	fmt.Println("[DEBUG] Stopping hotkey manager...")
	return m.Unregister()
}
