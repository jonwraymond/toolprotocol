package resource

import "errors"

// Sentinel errors for resource operations.
// All errors use the "resource: " prefix for consistent error identification.
var (
	// ErrResourceNotFound is returned when a resource cannot be found.
	ErrResourceNotFound = errors.New("resource: not found")

	// ErrProviderNotFound is returned when no provider handles a URI scheme.
	ErrProviderNotFound = errors.New("resource: provider not found")

	// ErrInvalidURI is returned when a URI is invalid.
	ErrInvalidURI = errors.New("resource: invalid URI")

	// ErrDuplicateProvider is returned when registering a duplicate provider.
	ErrDuplicateProvider = errors.New("resource: provider already registered")

	// ErrNotSubscribed is returned when unsubscribing from a non-subscribed resource.
	ErrNotSubscribed = errors.New("resource: not subscribed")
)

// ResourceError wraps an error with resource context.
type ResourceError struct {
	// URI is the resource URI that caused the error.
	URI string

	// Op is the operation that failed.
	Op string

	// Err is the underlying error.
	Err error
}

// Error returns the error message.
func (e *ResourceError) Error() string {
	if e.Err == nil {
		return "resource " + e.URI + ": " + e.Op
	}
	return "resource " + e.URI + ": " + e.Op + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *ResourceError) Unwrap() error {
	return e.Err
}
