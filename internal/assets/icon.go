package assets

import (
	_ "embed"
	"fmt"

	"github.com/jonesrussell/godo/internal/logger"
)

//go:embed favicon.ico
var IconBytes []byte

// GetIcon returns the application icon bytes and validates the icon data
func GetIcon() ([]byte, error) {
	if len(IconBytes) == 0 {
		return nil, fmt.Errorf("icon file is empty")
	}

	// Basic ICO file validation (check for ICO header magic numbers)
	if len(IconBytes) < 4 || IconBytes[0] != 0x00 || IconBytes[1] != 0x00 ||
		IconBytes[2] != 0x01 || IconBytes[3] != 0x00 {
		return nil, fmt.Errorf("invalid ICO file format")
	}

	logger.Debug("Icon loaded successfully, size: %d bytes", len(IconBytes))
	return IconBytes, nil
}
