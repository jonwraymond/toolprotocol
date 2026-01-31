package wire

import "errors"

// Sentinel errors for wire operations.
var (
	// ErrUnsupportedFormat is returned when an unknown wire format is requested.
	ErrUnsupportedFormat = errors.New("unsupported wire format")

	// ErrEncodeFailure is returned when encoding fails.
	ErrEncodeFailure = errors.New("failed to encode wire format")

	// ErrDecodeFailure is returned when decoding fails.
	ErrDecodeFailure = errors.New("failed to decode wire format")
)
