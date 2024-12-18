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
	keyName := strings.ToLower(hook.RawcodetoKeychar(ev.Rawcode))
	logger.Debug("Key event received", "key", keyName, "kind", ev.Kind)

	if ev.Kind == hook.KeyDown {
		h.active[keyName] = true
		logger.Debug("Active keys", "keys", h.active)
	} else if ev.Kind == hook.KeyUp {
		delete(h.active, keyName)
		logger.Debug("Active keys after release", "keys", h.active)
	}

	// Check if our hotkey combination is active
	if h.isHotkeyActive() {
		logger.Debug("Hotkey combination detected", "binding", h.binding)
		h.eventChan <- struct{}{}
		// Clear the active keys to prevent repeated triggers
		h.active = make(map[string]bool)
	}
}

func (h *windowsHotkeyManager) isHotkeyActive() bool {
	if h.binding.Key == "" {
		logger.Debug("No key binding set")
		return false
	}

	// Check if the main key is pressed
	mainKey := strings.ToLower(h.binding.Key)
	if !h.active[mainKey] {
		return false
	}

	// Check if all modifiers are pressed
	for _, mod := range h.binding.Modifiers {
		modKey := strings.ToLower(mod)
		if !h.active[modKey] {
			logger.Debug("Missing modifier", "modifier", modKey)
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
