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
	// Add more key mappings as needed
	default:
		return fmt.Errorf("unsupported key: %s", m.binding.Key)
	}

	// Create the hotkey
	fmt.Println("[DEBUG] Creating hotkey instance...")
	hk := hotkey.New(mods, key)

	// Try to register with retries
	var err error
	for i := 0; i < 3; i++ {
		fmt.Printf("[DEBUG] Attempting to register hotkey (attempt %d/3)...\n", i+1)
		err = hk.Register()
		if err == nil {
			break
		}
		fmt.Printf("[WARN] Failed to register hotkey (attempt %d/3): %v\n", i+1, err)
		time.Sleep(time.Second) // Wait before retrying
	}

	if err != nil {
		fmt.Printf("[ERROR] All attempts to register hotkey failed: %v\n", err)
		return fmt.Errorf("failed to register hotkey after 3 attempts: %w", err)
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
