// Package systray provides system tray functionality
package systray

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// SetupSystray configures the system tray icon and menu
func SetupSystray(app fyne.App, mainWindow fyne.Window) {
	if desk, ok := app.(desktop.App); ok {
		menu := fyne.NewMenu("Godo",
			fyne.NewMenuItem("Show", func() {
				mainWindow.Show()
				mainWindow.RequestFocus()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				app.Quit()
			}),
		)
		desk.SetSystemTrayMenu(menu)
	}
}
