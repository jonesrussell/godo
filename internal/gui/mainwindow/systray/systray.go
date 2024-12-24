package systray

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/logger"
)

type Systray struct {
	app   fyne.App
	log   logger.Logger
	ready bool
	menu  *fyne.Menu
	icon  fyne.Resource
}

func New(app fyne.App, log logger.Logger) *Systray {
	return &Systray{
		app: app,
		log: log,
	}
}

func (s *Systray) Setup(menu *fyne.Menu) {
	s.menu = menu
	s.ready = true
}

func (s *Systray) SetIcon(resource fyne.Resource) {
	s.icon = resource
}

func (s *Systray) IsReady() bool {
	return s.ready
}
