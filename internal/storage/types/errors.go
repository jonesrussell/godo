package types

import "errors"

// Common storage errors
var (
	ErrEmptyID                 = errors.New("empty note ID")
	ErrDuplicateID             = errors.New("duplicate note ID")
	ErrNoteNotFound            = errors.New("note not found")
	ErrStoreClosed             = errors.New("store is closed")
	ErrTransactionNotSupported = errors.New("transactions not supported")
)
