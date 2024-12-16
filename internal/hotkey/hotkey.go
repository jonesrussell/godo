package hotkey

import (
	"context"
	"syscall"
	"unsafe"

	"github.com/jonesrussell/godo/internal/logger"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procRegisterHotKey   = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey = user32.NewProc("UnregisterHotKey")
	procGetMessage       = user32.NewProc("GetMessageW")
)

const (
	MOD_ALT     = 0x0001
	MOD_CONTROL = 0x0002
	WM_HOTKEY   = 0x0312
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
	return &HotkeyManager{
		showCallback: showCallback,
	}
}

func (h *HotkeyManager) Start(ctx context.Context) error {
	logger.Info("Registering hotkey Ctrl+Alt+T...")

	// Register Ctrl+Alt+T
	ret, _, err := procRegisterHotKey.Call(
		0, // NULL window handle
		1, // hotkey ID
		uintptr(MOD_CONTROL|MOD_ALT),
		uintptr('T'),
	)
	if ret == 0 {
		logger.Error("Failed to register hotkey: %v (ret=%d)", err, ret)
		return err
	}
	logger.Info("Successfully registered hotkey (ret=%d)", ret)

	var msg MSG
	done := make(chan struct{})

	go func() {
		logger.Info("Starting Windows message loop")
		defer close(done)
		defer func() {
			ret, _, err := procUnregisterHotKey.Call(0, 1)
			if ret == 0 {
				logger.Error("Failed to unregister hotkey: %v", err)
			} else {
				logger.Info("Unregistered hotkey Ctrl+Alt+T")
			}
		}()

		for {
			select {
			case <-ctx.Done():
				logger.Info("Context cancelled, stopping hotkey listener")
				return
			default:
				logger.Debug("Waiting for message...")
				ret, _, err := procGetMessage.Call(
					uintptr(unsafe.Pointer(&msg)),
					0,
					0,
					0,
				)
				if ret == 0 {
					logger.Info("Message loop ended (ret=0)")
					return
				}
				if int32(ret) < 0 {
					logger.Error("Error getting message: %v (ret=%d)", err, ret)
					continue
				}
				logger.Debug("Got message: type=%d, wparam=%d, lparam=%d",
					msg.Message, msg.WParam, msg.LParam)
				if msg.Message == WM_HOTKEY {
					logger.Info("Hotkey triggered! ID=%d", msg.WParam)
					h.showCallback()
				}
			}
		}
	}()

	logger.Info("Hotkey manager initialized and waiting for events")
	<-done
	return nil
}
