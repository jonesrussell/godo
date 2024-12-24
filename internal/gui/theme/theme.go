package theme

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed favicon.ico
var iconData []byte

// AppIcon returns the application icon
func AppIcon() fyne.Resource {
	return fyne.NewStaticResource("favicon.ico", iconData)
}
