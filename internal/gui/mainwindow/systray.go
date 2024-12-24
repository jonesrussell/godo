package systray

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/logger"
)

type Service struct {
	app   fyne.App
	log   logger.Logger
	ready bool
	menu  *fyne.Menu
	icon  fyne.Resource
}

func New(app fyne.App, log logger.Logger) *Service {
	return &Service{
		app: app,
		log: log,
	}
}

func (s *Service) Setup(menu *fyne.Menu) {
	s.menu = menu
	s.ready = true
}

func (s *Service) SetIcon(resource fyne.Resource) {
	s.icon = resource
}

func (s *Service) IsReady() bool {
	return s.ready
}
