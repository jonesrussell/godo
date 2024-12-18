//go:build windows
// +build windows

package systray

import (
	"fmt"
	"sync"
	"time"

	"github.com/getlantern/systray"
	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/logger"
)

func init() {
	newManager = func() Manager {
		return &windowsSystray{
			ready: make(chan struct{}),
		}
	}
}

type windowsSystray struct {
	defaultManager
	ready chan struct{}
	once  sync.Once
}

func (w *windowsSystray) Setup() error {
	icon, err := assets.GetIcon()
	if err != nil {
		logger.Error("Failed to load icon: %v", err)
		return err
	}

	go systray.Run(func() {
		systray.SetIcon(icon)
		systray.SetTooltip("Godo - Quick Note Todo App")

		// Signal that systray is ready
		w.once.Do(func() {
			close(w.ready)
		})
	}, func() {
		logger.Debug("Cleaning up systray")
	})

	// Wait for systray to be ready
	select {
	case <-w.ready:
		logger.Info("Systray initialized successfully")
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for systray initialization")
	}
}

func (w *windowsSystray) AddMenuItem(label, tooltip string) *systray.MenuItem {
	// Wait for systray to be ready before adding menu items
	select {
	case <-w.ready:
		return systray.AddMenuItem(label, tooltip)
	}
}
