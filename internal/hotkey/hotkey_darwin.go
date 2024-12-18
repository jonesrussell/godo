package hotkey

import (
	"context"

	"github.com/robotn/gohook"
)

func (h *HotkeyManager) Start(ctx context.Context) error {
	keyCombo := []string{"cmd", "alt", "g"}

	go func() {
		gohook.Register(gohook.KeyDown, keyCombo, func(e gohook.Event) {
			select {
			case h.eventChan <- struct{}{}:
			default:
				// Channel is full, skip this event
			}
		})
		s := gohook.Start()
		<-ctx.Done()
		gohook.End()
		<-s
	}()
	return nil
}
