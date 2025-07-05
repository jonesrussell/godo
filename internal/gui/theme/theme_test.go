package theme

import (
	"testing"

	"fyne.io/fyne/v2"
)

func TestAppIcon(t *testing.T) {
	icon := AppIcon()

	if icon == nil {
		t.Error("AppIcon() returned nil")
	}

	if icon.Name() != "favicon.ico" {
		t.Errorf("Expected icon name to be 'favicon.ico', got '%s'", icon.Name())
	}

	if len(icon.Content()) == 0 {
		t.Error("Icon content is empty")
	}

	// Verify it implements the fyne.Resource interface
	var _ = icon
}
