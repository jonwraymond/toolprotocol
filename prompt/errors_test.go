package prompt

import (
	"errors"
	"testing"
)

func TestErrPromptNotFound(t *testing.T) {
	if ErrPromptNotFound.Error() != "prompt: not found" {
		t.Errorf("Error() = %q, want %q", ErrPromptNotFound.Error(), "prompt: not found")
	}
}

func TestErrMissingArgument(t *testing.T) {
	if ErrMissingArgument.Error() != "prompt: missing required argument" {
		t.Errorf("Error() = %q, want %q", ErrMissingArgument.Error(), "prompt: missing required argument")
	}
}

func TestErrInvalidPrompt(t *testing.T) {
	if ErrInvalidPrompt.Error() != "prompt: invalid" {
		t.Errorf("Error() = %q, want %q", ErrInvalidPrompt.Error(), "prompt: invalid")
	}
}

func TestErrHandlerFailed(t *testing.T) {
	if ErrHandlerFailed.Error() != "prompt: handler failed" {
		t.Errorf("Error() = %q, want %q", ErrHandlerFailed.Error(), "prompt: handler failed")
	}
}

func TestPromptError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *PromptError
		want string
	}{
		{
			name: "with underlying error",
			err: &PromptError{
				PromptName: "greeting",
				Op:         "get",
				Err:        ErrPromptNotFound,
			},
			want: "prompt greeting: get: prompt: not found",
		},
		{
			name: "without underlying error",
			err: &PromptError{
				PromptName: "greeting",
				Op:         "register",
				Err:        nil,
			},
			want: "prompt greeting: register",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPromptError_Unwrap(t *testing.T) {
	underlying := ErrPromptNotFound
	err := &PromptError{
		PromptName: "greeting",
		Op:         "get",
		Err:        underlying,
	}

	if !errors.Is(err, underlying) {
		t.Error("errors.Is should return true for underlying error")
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != underlying {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlying)
	}
}
