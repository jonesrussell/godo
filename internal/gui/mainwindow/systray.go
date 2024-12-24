package systray

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/jonesrussell/godo/internal/logger"
)

type Service struct {
	desk  desktop.App
	log   logger.Logger
	ready bool
}

func New(app fyne.App, log logger.Logger) *Service {
	desk, _ := app.(desktop.App)
	return &Service{
		desk: desk,
		log:  log,
	}
}

func (s *Service) Setup(menu *fyne.Menu) {
	if s.desk != nil {
		s.desk.SetSystemTrayMenu(menu)
		time.Sleep(100 * time.Millisecond)
		s.ready = true
		s.log.Debug("System tray menu setup complete")
	}
}

func (s *Service) SetIcon(resource fyne.Resource) {
	if s.desk != nil {
		if !s.ready {
			s.log.Warn("Attempted to set icon before systray was ready")
			return
		}
		s.desk.SetSystemTrayIcon(resource)
		s.log.Debug("System tray icon set successfully")
	}
}

func (s *Service) IsReady() bool {
	return s.ready
}
