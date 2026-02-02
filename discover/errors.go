package discover

import "errors"

// Sentinel errors for discovery operations.
// All errors use the "discover: " prefix for consistent error identification.
var (
	// ErrNotFound is returned when a service is not found.
	ErrNotFound = errors.New("discover: service not found")

	// ErrDuplicate is returned when registering a service with an existing ID.
	ErrDuplicate = errors.New("discover: service already registered")

	// ErrInvalidService is returned when a service fails validation.
	ErrInvalidService = errors.New("discover: invalid service")
)
