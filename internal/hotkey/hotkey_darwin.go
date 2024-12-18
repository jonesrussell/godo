//go:build darwin

package hotkey

import (
	"context"
)

func (h *HotkeyManager) Start(ctx context.Context) error {
	// For now, let's implement a placeholder that satisfies the interface
	// TODO: Implement proper macOS hotkey support
	go func() {
		<-ctx.Done()
	}()
	return nil
}
