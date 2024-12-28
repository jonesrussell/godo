//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"fmt"
	"runtime"
	"time"

	"golang.design/x/hotkey"
)

type platformManager struct {
	hk        *hotkey.Hotkey
	quickNote QuickNoteService
}

func newPlatformManager(quickNote QuickNoteService) Manager {
	fmt.Printf("[DEBUG] Creating hotkey manager (OS: %s)\n", runtime.GOOS)
	if quickNote == nil {
		panic("quickNote service cannot be nil")
	}
	return &platformManager{
		quickNote: quickNote,
	}
}

func (m *platformManager) Register() error {
	fmt.Printf("[DEBUG] Registering hotkey (Ctrl+Shift+N) on %s\n", runtime.GOOS)

	// Create the hotkey
	fmt.Println("[DEBUG] Creating hotkey instance...")
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyN)

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
