package transport

import "errors"

// Sentinel errors for transport operations.
var (
	// ErrTransportClosed is returned when operations are attempted on a closed transport.
	ErrTransportClosed = errors.New("transport closed")

	// ErrAlreadyServing is returned when Serve is called on a transport that is already serving.
	ErrAlreadyServing = errors.New("transport already serving")

	// ErrInvalidConfig is returned when transport configuration is invalid.
	ErrInvalidConfig = errors.New("invalid transport configuration")
)
