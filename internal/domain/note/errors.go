package note

import "fmt"

// ErrorKind represents the type of note error
type ErrorKind int

const (
	// Unknown represents an unknown error
	Unknown ErrorKind = iota
	// ValidationFailed indicates a validation error
	ValidationFailed
	// NotFound indicates a note was not found
	NotFound
	// StorageError indicates a storage-related error
	StorageError
)

// Error represents a note-related error
type Error struct {
	// Op is the operation that failed
	Op string
	// Kind is the type of error
	Kind ErrorKind
	// Msg contains the error details
	Msg string
	// Err is the underlying error if any
	Err error
}

// Error returns the error message
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Msg, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Msg)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// Is reports whether the target error matches this error
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Kind == t.Kind
}
