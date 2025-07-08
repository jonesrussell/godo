// Package systray provides system tray functionality
package systray

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/gui/theme"
	"github.com/jonesrussell/godo/internal/utils"
)

// SetupSystray configures the system tray icon and menu
func SetupSystray(app fyne.App, mainWindow fyne.Window, quickNote gui.QuickNote, logPath, errorLogPath string) {
	if desk, ok := app.(desktop.App); ok {
		desk.SetSystemTrayIcon(theme.AppIcon())

		menu := fyne.NewMenu("Godo",
			fyne.NewMenuItem("Show", func() {
				mainWindow.Show()
				mainWindow.RequestFocus()
			}),
			fyne.NewMenuItem("Quick Note", func() {
				quickNote.Show()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("View Logs", func() {
				if err := utils.OpenLogFile(logPath); err != nil {
					// If main log fails, try error log
					if logErr := utils.OpenErrorLogFile(errorLogPath); logErr != nil {
						// Both log files failed to open - this is a non-critical error
						// since it's just a convenience feature for viewing logs
						_ = logErr // explicitly ignore the error
					}
				}
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				app.Quit()
			}),
		)
		desk.SetSystemTrayMenu(menu)
	}
}
