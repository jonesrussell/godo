//go:build darwin
// +build darwin

package hotkey

import (
	"context"

	"github.com/jonesrussell/godo/internal/types"
)

func init() {
	newPlatformHotkeyManager = func() (HotkeyManager, error) {
		return &darwinHotkeyManager{
			eventChan: make(chan struct{}),
		}, nil
	}
}

type darwinHotkeyManager struct {
	eventChan chan struct{}
}

func (h *darwinHotkeyManager) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (h *darwinHotkeyManager) Stop() error {
	return nil
}

func (h *darwinHotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}

// RegisterHotkey implements the HotkeyManager interface
func (h *darwinHotkeyManager) RegisterHotkey(binding types.HotkeyBinding) error {
	// TODO: Implement Darwin-specific hotkey registration using CGEventTap
	// For now, return nil as a placeholder
	return nil
}
