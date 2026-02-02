package prompt

import "errors"

// Sentinel errors for prompt operations.
// All errors use the "prompt: " prefix for consistent error identification.
var (
	// ErrPromptNotFound is returned when a prompt cannot be found.
	ErrPromptNotFound = errors.New("prompt: not found")

	// ErrMissingArgument is returned when a required argument is missing.
	ErrMissingArgument = errors.New("prompt: missing required argument")

	// ErrInvalidPrompt is returned when a prompt is invalid.
	ErrInvalidPrompt = errors.New("prompt: invalid")

	// ErrHandlerFailed is returned when a prompt handler fails.
	ErrHandlerFailed = errors.New("prompt: handler failed")

	// ErrDuplicatePrompt is returned when registering a duplicate prompt.
	ErrDuplicatePrompt = errors.New("prompt: already registered")
)

// PromptError wraps an error with prompt context.
type PromptError struct {
	// PromptName is the name of the prompt that caused the error.
	PromptName string

	// Op is the operation that failed.
	Op string

	// Err is the underlying error.
	Err error
}

// Error returns the error message.
func (e *PromptError) Error() string {
	if e.Err == nil {
		return "prompt " + e.PromptName + ": " + e.Op
	}
	return "prompt " + e.PromptName + ": " + e.Op + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *PromptError) Unwrap() error {
	return e.Err
}
