package wire

import "errors"

// Sentinel errors for wire operations.
// All errors use the "wire: " prefix for consistent error identification.
var (
	// ErrUnsupportedFormat is returned when an unknown wire format is requested.
	ErrUnsupportedFormat = errors.New("wire: unsupported format")

	// ErrEncodeFailure is returned when encoding fails.
	ErrEncodeFailure = errors.New("wire: encode failed")

	// ErrDecodeFailure is returned when decoding fails.
	ErrDecodeFailure = errors.New("wire: decode failed")
)
