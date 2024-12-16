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
}

func New(showCallback func()) *HotkeyManager {
	log.Println("Creating new HotkeyManager...")
	return &HotkeyManager{
		showCallback: showCallback,
	}
}

func (h *HotkeyManager) Start(ctx context.Context) error {
	log.Println("Registering hotkey Ctrl+Alt+T...")

	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	// Set keys
	kb.SetKeys(keybd_event.VK_T)
	// Set modifiers
	kb.HasCTRL(true)
	kb.HasALT(true)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping hotkey listener")
				return
			default:
				if err := kb.Launching(); err != nil {
					log.Printf("Error checking hotkey: %v\n", err)
				} else {
					h.showCallback()
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return nil
}
