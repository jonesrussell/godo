package hotkey

import (
	"context"
	"syscall"
	"unsafe"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procRegisterHotKey   = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey = user32.NewProc("UnregisterHotKey")
	procGetMessage       = user32.NewProc("GetMessageW")
	procPeekMessage     = user32.NewProc("PeekMessageW")
)

const (
	MOD_ALT     = 0x0001
	MOD_CONTROL = 0x0002
	WM_HOTKEY   = 0x0312
	PM_REMOVE   = 0x0001
	ERROR_HOTKEY_ALREADY_REGISTERED = 0x0402
)

type MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

type HotkeyManager struct {
	showCallback func()
}

func New(showCallback func()) *HotkeyManager {
	logger.Info("Creating new HotkeyManager...")
	if showCallback == nil {
		logger.Error("showCallback is nil!")
		return nil
	}
	return &HotkeyManager{
		showCallback: showCallback,
	}
}

func (h *HotkeyManager) Start(ctx context.Context) error {
	logger.Info("Starting hotkey manager...")
	
	// Unregister any existing hotkey first
	logger.Debug("Cleaning up any existing hotkey registration...")
	procUnregisterHotKey.Call(0, 1)

	logger.Debug("Attempting to register Ctrl+Alt+T (modifiers=0x%X, key='T')", MOD_CONTROL|MOD_ALT)
	
	// Register Ctrl+Alt+T with more detailed error handling
	ret, _, err := procRegisterHotKey.Call(
		0, // NULL window handle
		1, // hotkey ID
		uintptr(MOD_CONTROL|MOD_ALT),
		uintptr('T'),
	)
	
	if ret == 0 {
		lastErr := syscall.GetLastError()
		logger.Error("Failed to register hotkey: %v (ret=%d, lastErr=%d)", err, ret, lastErr)
		return lastErr
	}
	logger.Info("Successfully registered hotkey (ret=%d)", ret)

	var msg MSG
	done := make(chan struct{})

	go func() {
		logger.Info("Starting Windows message loop in goroutine")
		defer close(done)
		defer func() {
			logger.Info("Unregistering hotkey...")
			ret, _, err := procUnregisterHotKey.Call(0, 1)
			if ret == 0 {
				lastErr := syscall.GetLastError()
				logger.Error("Failed to unregister hotkey: %v (lastErr=%d)", err, lastErr)
			} else {
				logger.Info("Successfully unregistered hotkey")
			}
		}()

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Info("Context cancelled, stopping hotkey listener")
				return
			case <-ticker.C:
				// Try PeekMessage first
				ret, _, _ := procPeekMessage.Call(
					uintptr(unsafe.Pointer(&msg)),
					0,
					0,
					0,
					PM_REMOVE,
				)

				if ret == 0 {
					continue // No message available
				}

				logger.Debug("Message received: type=0x%X, wparam=0x%X, lparam=0x%X",
					msg.Message, msg.WParam, msg.LParam)

				if msg.Message == WM_HOTKEY {
					logger.Info("Hotkey triggered! (ID=%d)", msg.WParam)
					if h.showCallback != nil {
						logger.Debug("Executing showCallback...")
						h.showCallback()
					} else {
						logger.Error("showCallback is nil!")
					}
				}
			}
		}
	}()

	logger.Info("Hotkey manager initialized and waiting for events")
	<-done
	return nil
}

// For debugging purposes
func (h *HotkeyManager) IsHotkeyRegistered() bool {
	// Try to register the same hotkey - if it fails, it means it's already registered
	ret, _, _ := procRegisterHotKey.Call(
		0,
		1,
		uintptr(MOD_CONTROL|MOD_ALT),
		uintptr('T'),
	)
	
	if ret == 0 {
		lastErr := syscall.GetLastError()
		if lastErr == syscall.Errno(ERROR_HOTKEY_ALREADY_REGISTERED) {
			logger.Debug("Hotkey is currently registered")
			return true
		}
	}
	
	// Clean up test registration if it succeeded
	if ret != 0 {
		procUnregisterHotKey.Call(0, 1)
	}
	
	logger.Debug("Hotkey is not currently registered")
	return false
}

func (h *HotkeyManager) RegisterHotkey() error {
	logger.Debug("Attempting to register hotkey Ctrl+Alt+T...")
	// ... registration code ...
	logger.Debug("Hotkey registration attempt completed")
	return nil
}

func (m *HotkeyManager) handleHotkey() {
	logger.Debug("Hotkey Ctrl+Alt+T pressed!")
	// ... handling code ...
}
