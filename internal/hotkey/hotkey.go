package hotkey

import (
	"context"
	"log"
	"runtime"
	"time"

	"github.com/micmonay/keybd_event"
)

type HotkeyManager struct {
	showCallback func()
	kb           keybd_event.KeyBonding
}

func New(showCallback func()) (*HotkeyManager, error) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return nil, err
	}

	// Initialize keyboard
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	// Set keys to listen for (Ctrl+Alt+T)
	kb.SetKeys(keybd_event.VK_T)
	kb.HasCTRL(true)
	kb.HasALT(true)

	return &HotkeyManager{
		showCallback: showCallback,
		kb:           kb,
	}, nil
}

func (h *HotkeyManager) Start(ctx context.Context) error {
	log.Println("Starting hotkey listener (Ctrl+Alt+T)...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping hotkey listener...")
			return nil
		default:
			// Launch a key press
			err := h.kb.Press()
			if err != nil {
				log.Printf("Error pressing key: %v", err)
				continue
			}

			time.Sleep(10 * time.Millisecond)

			err = h.kb.Release()
			if err != nil {
				log.Printf("Error releasing key: %v", err)
				continue
			}

			h.showCallback()

			// Small delay to prevent multiple triggers
			time.Sleep(300 * time.Millisecond)
		}
	}
}
