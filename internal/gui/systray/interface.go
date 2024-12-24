package systray

import "fyne.io/fyne/v2"

type Interface interface {
	Setup(menu *fyne.Menu)
	SetIcon(resource fyne.Resource)
	IsReady() bool
}
