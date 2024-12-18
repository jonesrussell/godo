//go:build windows
// +build windows

package hotkey

import (
	"context"
	"strings"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/types"
	hook "github.com/robotn/gohook"
)

type windowsHotkeyManager struct {
	eventChan chan struct{}
	binding   types.HotkeyBinding
	active    map[string]bool
}

func init() {
	newPlatformHotkeyManager = func() (HotkeyManager, error) {
		return &windowsHotkeyManager{
			eventChan: make(chan struct{}),
			active:    make(map[string]bool),
		}, nil
	}
}

func (h *windowsHotkeyManager) RegisterHotkey(binding types.HotkeyBinding) error {
	h.binding = binding
	logger.Debug("Registered hotkey binding: %v", binding)
	return nil
}

func (h *windowsHotkeyManager) Start(ctx context.Context) error {
	logger.Debug("Starting Windows hotkey manager")
	go func() {
		evChan := hook.Start()
		defer hook.End()

		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-evChan:
				if ev.Kind == hook.KeyHold || ev.Kind == hook.KeyDown {
					h.handleKeyEvent(ev)
				}
			}
		}
	}()
	return nil
}

func (h *windowsHotkeyManager) handleKeyEvent(ev hook.Event) {
	// Convert key to proper name
	keyName := strings.ToLower(hook.RawcodetoKeychar(ev.Rawcode))

	// Map special keys
	switch keyName {
	case "؛":
		keyName = "ctrl"
	case "$":
		keyName = "alt"
	}

	logger.Debug("Key event received",
		"key", keyName,
		"kind", ev.Kind,
		"active_keys", h.active)

	if ev.Kind == hook.KeyDown {
		h.active[keyName] = true
	} else if ev.Kind == hook.KeyUp {
		delete(h.active, keyName)
	}

	// Check if hotkey combination is active
	if h.isHotkeyActive() {
		logger.Info("Hotkey combination detected", "binding", h.binding)
		h.eventChan <- struct{}{}
		// Reset active keys
		h.active = make(map[string]bool)
	}
}

func (h *windowsHotkeyManager) isHotkeyActive() bool {
	if h.binding.Key == "" {
		return false
	}

	// Check main key
	if !h.active[strings.ToLower(h.binding.Key)] {
		return false
	}

	// Check all modifiers
	for _, mod := range h.binding.Modifiers {
		if !h.active[strings.ToLower(mod)] {
			return false
		}
	}

	return true
}

func (h *windowsHotkeyManager) Stop() error {
	hook.End()
	return nil
}

func (h *windowsHotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}
