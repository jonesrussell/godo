package utils

import (
	"os/exec"
	"runtime"
)

// OpenLogFile opens the given log file in the system's default text editor or viewer.
func OpenLogFile(path string) error {
	return openFile(path)
}

// OpenErrorLogFile opens the given error log file in the system's default text editor or viewer.
func OpenErrorLogFile(path string) error {
	return openFile(path)
}

// openFile opens a file using the system's default application for the file type.
func openFile(path string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", path).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", "", path).Start()
	default: // Linux and others
		return exec.Command("xdg-open", path).Start()
	}
}
