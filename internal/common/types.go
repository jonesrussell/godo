package common

// LogConfig holds logging configuration
type LogConfig struct {
	Level       string   `mapstructure:"level" yaml:"level"`
	Console     bool     `mapstructure:"console" yaml:"console"`
	File        bool     `mapstructure:"file" yaml:"file"`
	FilePath    string   `mapstructure:"file_path" yaml:"file_path"`
	Output      []string `mapstructure:"output" yaml:"output"`
	ErrorOutput []string `mapstructure:"error_output" yaml:"error_output"`
}

// HotkeyBinding represents a keyboard shortcut configuration
type HotkeyBinding struct {
	Modifiers []string `yaml:"modifiers"`
	Key       string   `yaml:"key"`
}

// String implements the Stringer interface for HotkeyBinding
func (h HotkeyBinding) String() string {
	return h.Key
}

// Error represents a domain error
type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Code + ": " + e.Message + ": " + e.Err.Error()
	}
	return e.Code + ": " + e.Message
}

// NewError creates a new domain error
func NewError(code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
