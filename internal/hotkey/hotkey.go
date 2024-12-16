package hotkey

import (
	"context"
	"log"
	"syscall"
	"unsafe"
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
	log.Println("Creating new HotkeyManager...")
	return &HotkeyManager{
		showCallback: showCallback,
	}
}

func (h *HotkeyManager) Start(ctx context.Context) error {
	log.Println("Registering hotkey Ctrl+Alt+T...")

	// Register Ctrl+Alt+T
	ret, _, err := procRegisterHotKey.Call(
		0, // NULL window handle
		1, // hotkey ID
		uintptr(MOD_CONTROL|MOD_ALT),
		uintptr('T'),
	)
	if ret == 0 {
		return err
	}

	var msg MSG
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				procUnregisterHotKey.Call(0, 1)
				log.Println("Unregistered hotkey Ctrl+Alt+T")
				return
			default:
				ret, _, _ := procGetMessage.Call(
					uintptr(unsafe.Pointer(&msg)),
					0,
					0,
					0,
				)
				if ret != 0 && msg.Message == WM_HOTKEY {
					log.Println("Hotkey triggered!")
					h.showCallback()
				}
			}
		}
	}()

	<-done
	return nil
}
