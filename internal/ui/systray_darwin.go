//go:build darwin
// +build darwin

package ui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void setUpMenuBar(const char* title, const char* tooltip) {
    NSStatusBar *statusBar = [NSStatusBar systemStatusBar];
    NSStatusItem *statusItem = [statusBar statusItemWithLength:NSVariableStatusItemLength];

    NSString *titleStr = [NSString stringWithUTF8String:title];
    NSString *tooltipStr = [NSString stringWithUTF8String:tooltip];

    [statusItem setTitle:titleStr];
    [statusItem setToolTip:tooltipStr];
}
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/logger"
)

func init() {
	// Override the default systray manager with Darwin-specific implementation
	newSystrayManager = func() SystrayManager {
		return &darwinSystray{}
	}
}

type darwinSystray struct{}

func (s *darwinSystray) Setup() error {
	_, err := assets.GetIcon()
	if err != nil {
		logger.Error("Failed to load icon: %v", err)
		return fmt.Errorf("failed to load icon: %w", err)
	}

	title := C.CString("Godo")
	tooltip := C.CString("Godo - Quick Note Todo App")
	defer C.free(unsafe.Pointer(title))
	defer C.free(unsafe.Pointer(tooltip))

	C.setUpMenuBar(title, tooltip)

	return nil
}
