// Package systray provides system tray functionality
package systray

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/jonesrussell/godo/internal/infrastructure/gui"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/theme"
	"github.com/jonesrussell/godo/internal/infrastructure/platform"
	"github.com/jonesrussell/godo/internal/shared/utils"
)

// SetupSystray configures the system tray icon and menu
func SetupSystray(app fyne.App, mainWindow fyne.Window, quickNote gui.QuickNote, logPath, errorLogPath string) error {
	// Check if systray is supported in this environment
	if !platform.SupportsGUI() {
		if platform.IsWSL2() {
			fmt.Println("Systray not supported in WSL2 environment - skipping setup")
		} else if platform.IsHeadless() {
			fmt.Println("Systray not supported in headless environment - skipping setup")
		}
		return nil // This is expected, not an error
	}

	desktopApp, ok := app.(desktop.App)
	if !ok {
		return fmt.Errorf("desktop features not available")
	}

	// Create and set the menu
	menu := createSystrayMenu(app, mainWindow, quickNote, logPath, errorLogPath)
	desktopApp.SetSystemTrayMenu(menu)

	// Set the icon
	icon := theme.AppIcon()
	if icon == nil {
		fmt.Println("Warning: AppIcon() returned nil - systray will have no icon")
		return nil // Continue without icon
	}

	desktopApp.SetSystemTrayIcon(icon)
	return nil
}

// createSystrayMenu creates the systray menu with enhanced View Logs functionality
func createSystrayMenu(
	app fyne.App,
	mainWindow fyne.Window,
	quickNote gui.QuickNote,
	logPath, errorLogPath string,
) *fyne.Menu {
	return fyne.NewMenu("Godo",
		fyne.NewMenuItem("Show", func() {
			mainWindow.Show()
			mainWindow.RequestFocus()
		}),
		fyne.NewMenuItem("Quick Note", func() {
			quickNote.Show()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("View Logs", func() {
			handleViewLogs(mainWindow, logPath, errorLogPath)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			app.Quit()
		}),
	)
}

// handleViewLogs handles the View Logs menu item with proper file existence checks
func handleViewLogs(mainWindow fyne.Window, logPath, errorLogPath string) {
	// Check if log files exist
	mainLogExists := fileExists(logPath)
	errorLogExists := fileExists(errorLogPath)

	if !mainLogExists && !errorLogExists {
		// No log files exist - show helpful message
		showLogFileInfo(mainWindow, logPath, errorLogPath)
		return
	}

	// Try to open existing log files
	var openedFiles []string
	var errors []string

	if mainLogExists {
		if err := utils.OpenLogFile(logPath); err != nil {
			errors = append(errors, fmt.Sprintf("Main log (%s): %v", logPath, err))
		} else {
			openedFiles = append(openedFiles, "Main log")
		}
	}

	if errorLogExists {
		if err := utils.OpenErrorLogFile(errorLogPath); err != nil {
			errors = append(errors, fmt.Sprintf("Error log (%s): %v", errorLogPath, err))
		} else {
			openedFiles = append(openedFiles, "Error log")
		}
	}

	// Show results to user
	if len(errors) > 0 {
		errorMsg := fmt.Sprintf("Some log files could not be opened:\n\n%s",
			formatErrorList(errors))
		dialog.ShowError(fmt.Errorf("log viewing issues: %s", errorMsg), mainWindow)
	}

	if len(openedFiles) > 0 {
		// Show success message briefly
		successMsg := fmt.Sprintf("Opened: %s", formatFileList(openedFiles))
		dialog.ShowInformation("Log Files", successMsg, mainWindow)
	}
}

// showLogFileInfo shows information about log file locations when they don't exist
func showLogFileInfo(mainWindow fyne.Window, logPath, errorLogPath string) {
	infoMsg := fmt.Sprintf(`No log files found yet.

Log files will be created at:
• Main log: %s
• Error log: %s

Log files are created when the application starts and logging occurs.`,
		logPath, errorLogPath)

	dialog.ShowInformation("Log Files", infoMsg, mainWindow)
}

// fileExists checks if a file exists and is readable
func fileExists(path string) bool {
	if path == "" {
		return false
	}

	// Resolve relative paths
	if !filepath.IsAbs(path) {
		// Try to resolve relative to current working directory
		if absPath, err := filepath.Abs(path); err == nil {
			path = absPath
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir() && info.Size() > 0
}

// formatErrorList formats a list of errors for display
func formatErrorList(errors []string) string {
	result := ""
	for i, err := range errors {
		if i > 0 {
			result += "\n"
		}
		result += "• " + err
	}
	return result
}

// formatFileList formats a list of files for display
func formatFileList(files []string) string {
	result := ""
	for i, file := range files {
		if i > 0 {
			result += ", "
		}
		result += file
	}
	return result
}
