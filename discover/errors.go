package discover

import "errors"

// Sentinel errors for discovery operations.
var (
	// ErrNotFound is returned when a service is not found.
	ErrNotFound = errors.New("service not found")

	// ErrDuplicate is returned when registering a service with an existing ID.
	ErrDuplicate = errors.New("service already registered")

	// ErrInvalidService is returned when a service fails validation.
	ErrInvalidService = errors.New("invalid service")
)
