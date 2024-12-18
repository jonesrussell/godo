package assets

import (
	_ "embed"

	"github.com/jonesrussell/godo/internal/logger"
)

//go:embed favicon.ico
var IconBytes []byte

// GetIcon loads and validates the application icon
func GetIcon() ([]byte, error) {
	logger.Debug("Loading application icon")

	if len(IconBytes) == 0 {
		logger.Error("Icon file is empty")
		return nil, ErrEmptyIcon
	}

	// Basic ICO file validation (check for ICO header magic numbers)
	if len(IconBytes) < 4 || IconBytes[0] != 0x00 || IconBytes[1] != 0x00 ||
		IconBytes[2] != 0x01 || IconBytes[3] != 0x00 {
		logger.Error("Invalid ICO file format")
		return nil, ErrInvalidIconFormat
	}

	logger.Debug("Successfully loaded icon",
		"size", len(IconBytes))
	return IconBytes, nil
}

// Define error types
var (
	ErrEmptyIcon         = NewAssetError("empty icon file")
	ErrInvalidIconFormat = NewAssetError("invalid ICO file format")
)

// AssetError represents an asset-related error
type AssetError struct {
	msg string
}

func (e AssetError) Error() string {
	return e.msg
}

// NewAssetError creates a new asset error
func NewAssetError(msg string) AssetError {
	return AssetError{msg: msg}
}
