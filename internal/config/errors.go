package config

import "fmt"

type (
	// LoadError represents a general error in processing the configuration file,
	// but not that the file was not found, which is handled by `NotFoundError`.
	LoadError struct {
		file    string
		message string
		err     error
	}

	// NotFoundError represents that either the default file could not be found in
	// the search paths, or the configuration file requested (which overrides the
	// default search, and so those files are not checked) could not be found.
	NotFoundError struct {
		file    string
		message string
		err     error
	}
)

// Error returns the error message for this error.
func (e *LoadError) Error() string {
	return fmt.Sprintf("%s: %s", e.message, e.err)
}

// NewLoadError creates a new `LoadError` error type with the provided `file`
// and `message` about the error, and the `err` from the upstream library.
func NewLoadError(file, message string, err error) error {
	return &LoadError{
		file:    file,
		message: message,
		err:     err,
	}
}

// Error returns the error message for this error.
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s: %s", e.message, e.err)
}

// NewNotFoundError creates a new `NotFoundError` error type with the provided
// `file` and `message` about the error, and the `err` from the upstream
// library.
func NewNotFoundError(file, message string, err error) error {
	return &NotFoundError{
		file:    file,
		message: message,
		err:     err,
	}
}
