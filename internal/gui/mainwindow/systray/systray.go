package systray

import (
	"fmt"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/jonesrussell/godo/internal/logger"
)

type Systray struct {
	app   fyne.App
	log   logger.Logger
	ready bool
	menu  *fyne.Menu
	icon  fyne.Resource
	desk  desktop.App
}

func New(app fyne.App, log logger.Logger) *Systray {
	s := &Systray{
		app: app,
		log: log,
	}

	// Detailed app type inspection
	appType := reflect.TypeOf(app)
	s.log.Debug("App type details",
		"type", fmt.Sprintf("%T", app),
		"implements_desktop", fmt.Sprintf("%v", appType.Implements(reflect.TypeOf((*desktop.App)(nil)).Elem())),
		"kind", appType.Kind().String(),
		"pkg_path", appType.PkgPath(),
	)

	// Check if system tray is supported
	desk, ok := app.(desktop.App)
	if ok {
		s.desk = desk
		s.log.Debug("System tray is supported",
			"app_type", fmt.Sprintf("%T", app),
			"desk_type", fmt.Sprintf("%T", desk),
		)

		// Test system tray capabilities
		if dt, ok := desk.(interface{ SystemTraySupported() bool }); ok {
			supported := dt.SystemTraySupported()
			s.log.Debug("System tray support explicitly checked", "supported", supported)
		} else {
			s.log.Warn("Cannot check SystemTraySupported - method not available")
		}
	} else {
		s.log.Error("System tray is not supported",
			"app_type", fmt.Sprintf("%T", app),
			"app_implements_desktop", fmt.Sprintf("%v", appType.Implements(reflect.TypeOf((*desktop.App)(nil)).Elem())),
		)
	}

	return s
}

func (s *Systray) Setup(menu *fyne.Menu) {
	s.log.Debug("Setting up system tray menu",
		"menu_label", menu.Label,
		"menu_type", fmt.Sprintf("%T", menu),
		"items_count", len(menu.Items),
	)
	s.menu = menu

	if s.desk != nil {
		s.log.Debug("Setting system tray menu via desktop.App",
			"desk_type", fmt.Sprintf("%T", s.desk),
			"menu_items", len(menu.Items),
		)
		for i, item := range menu.Items {
			s.log.Debug("Menu item details",
				"index", i,
				"label", item.Label,
				"type", fmt.Sprintf("%T", item),
				"has_action", item.Action != nil,
			)
		}
		s.desk.SetSystemTrayMenu(menu)
		s.ready = true
		s.log.Debug("System tray menu setup completed")
	} else {
		s.log.Error("Failed to set system tray menu - desktop.App is nil",
			"app_type", fmt.Sprintf("%T", s.app),
			"has_menu", s.menu != nil,
		)
	}
}

func (s *Systray) SetIcon(resource fyne.Resource) {
	s.log.Debug("Setting system tray icon",
		"resource_name", resource.Name(),
		"content_length", len(resource.Content()),
		"resource_type", fmt.Sprintf("%T", resource),
	)
	s.icon = resource

	if s.desk != nil {
		s.log.Debug("Setting system tray icon via desktop.App",
			"desk_type", fmt.Sprintf("%T", s.desk),
			"icon_name", resource.Name(),
		)
		s.desk.SetSystemTrayIcon(resource)
		s.log.Debug("System tray icon set successfully")
	} else {
		s.log.Error("Failed to set system tray icon - desktop.App is nil",
			"app_type", fmt.Sprintf("%T", s.app),
			"has_icon", s.icon != nil,
		)
	}
}

func (s *Systray) IsReady() bool {
	s.log.Debug("System tray state",
		"ready", s.ready,
		"has_menu", s.menu != nil,
		"has_icon", s.icon != nil,
		"has_desk", s.desk != nil,
		"app_type", fmt.Sprintf("%T", s.app),
		"desk_type", fmt.Sprintf("%T", s.desk),
	)
	return s.ready
}
