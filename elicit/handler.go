package elicit

import (
	"context"
	"sync"
	"time"
)

// DefaultHandler is a basic implementation of Handler that uses
// a callback function to process requests.
type DefaultHandler struct {
	mu      sync.RWMutex
	timeout time.Duration
	handler func(ctx context.Context, req *Request) (*Response, error)
}

// NewHandler creates a new DefaultHandler.
func NewHandler(opts ...HandlerOption) *DefaultHandler {
	h := &DefaultHandler{
		timeout: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// HandlerOption configures a DefaultHandler.
type HandlerOption func(*DefaultHandler)

// WithDefaultTimeout sets the default timeout for requests.
func WithDefaultTimeout(timeout time.Duration) HandlerOption {
	return func(h *DefaultHandler) {
		h.timeout = timeout
	}
}

// WithCallback sets the callback function for handling requests.
func WithCallback(fn func(ctx context.Context, req *Request) (*Response, error)) HandlerOption {
	return func(h *DefaultHandler) {
		h.handler = fn
	}
}

// Handle processes an elicitation request.
func (h *DefaultHandler) Handle(ctx context.Context, req *Request) (*Response, error) {
	if err := req.Validate(); err != nil {
		return nil, &ElicitError{
			RequestID: req.ID,
			Op:        "validate",
			Err:       err,
		}
	}

	h.mu.RLock()
	handler := h.handler
	timeout := h.timeout
	h.mu.RUnlock()

	if handler == nil {
		return nil, &ElicitError{
			RequestID: req.ID,
			Op:        "handle",
			Err:       ErrNoHandler,
		}
	}

	// Use request timeout if set, otherwise use default
	if req.Timeout > 0 {
		timeout = req.Timeout
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Channel for response
	type result struct {
		resp *Response
		err  error
	}
	done := make(chan result, 1)

	go func() {
		resp, err := handler(ctx, req)
		done <- result{resp, err}
	}()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return &Response{
				RequestID: req.ID,
				TimedOut:  true,
			}, nil
		}
		return &Response{
			RequestID: req.ID,
			Cancelled: true,
		}, nil
	case r := <-done:
		if r.err != nil {
			return nil, &ElicitError{
				RequestID: req.ID,
				Op:        "handle",
				Err:       r.err,
			}
		}
		return r.resp, nil
	}
}
