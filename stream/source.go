package stream

import "context"

// DefaultSource is the default Source implementation.
type DefaultSource struct {
	backpressure BackpressureMode
}

// NewSource creates a new DefaultSource.
func NewSource(opts ...SourceOption) *DefaultSource {
	s := &DefaultSource{
		backpressure: BackpressureBlock,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// SourceOption configures a DefaultSource.
type SourceOption func(*DefaultSource)

// WithBackpressure configures the backpressure mode for buffered streams.
func WithBackpressure(mode BackpressureMode) SourceOption {
	return func(s *DefaultSource) {
		s.backpressure = mode
	}
}

// NewStream creates a new unbuffered stream.
func (s *DefaultSource) NewStream(ctx context.Context) Stream {
	return newDefaultStream()
}

// NewBufferedStream creates a new buffered stream with the given size.
func (s *DefaultSource) NewBufferedStream(ctx context.Context, size int) Stream {
	return newBufferedStream(size, s.backpressure)
}

// DefaultBufferSize is the default buffer size for buffered streams.
const DefaultBufferSize = 100
