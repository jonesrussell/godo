package hotkey

import (
	"context"
	"log"

	hook "github.com/robotn/gohook"
)

type HotkeyManager struct {
	showCallback func()
}

func New(showCallback func()) *HotkeyManager {
	return &HotkeyManager{
		showCallback: showCallback,
	}
}

func (h *HotkeyManager) Start(ctx context.Context) error {
	log.Println("Registering hotkey Ctrl+Alt+T...")
	// Register Ctrl+Alt+T hotkey
	hook.Register(hook.KeyDown, []string{"t", "ctrl", "alt"}, func(e hook.Event) {
		log.Println("Hotkey triggered")
		h.showCallback()
	})

	log.Println("Starting hook process...")
	s := hook.Start()

	go func() {
		<-ctx.Done()
		log.Println("Context cancelled, ending hook")
		hook.End()
	}()

	log.Println("Entering hook process loop")
	<-hook.Process(s)
	return nil
}
