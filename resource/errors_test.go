package resource

import (
	"errors"
	"testing"
)

func TestErrResourceNotFound(t *testing.T) {
	if ErrResourceNotFound.Error() != "resource: not found" {
		t.Errorf("Error() = %q, want %q", ErrResourceNotFound.Error(), "resource: not found")
	}
}

func TestErrProviderNotFound(t *testing.T) {
	if ErrProviderNotFound.Error() != "resource: provider not found" {
		t.Errorf("Error() = %q, want %q", ErrProviderNotFound.Error(), "resource: provider not found")
	}
}

func TestErrInvalidURI(t *testing.T) {
	if ErrInvalidURI.Error() != "resource: invalid URI" {
		t.Errorf("Error() = %q, want %q", ErrInvalidURI.Error(), "resource: invalid URI")
	}
}

func TestErrDuplicateProvider(t *testing.T) {
	if ErrDuplicateProvider.Error() != "resource: provider already registered" {
		t.Errorf("Error() = %q, want %q", ErrDuplicateProvider.Error(), "resource: provider already registered")
	}
}

func TestResourceError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *ResourceError
		want string
	}{
		{
			name: "with underlying error",
			err: &ResourceError{
				URI: "file:///test",
				Op:  "read",
				Err: ErrResourceNotFound,
			},
			want: "resource file:///test: read: resource: not found",
		},
		{
			name: "without underlying error",
			err: &ResourceError{
				URI: "file:///test",
				Op:  "list",
				Err: nil,
			},
			want: "resource file:///test: list",
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

func TestResourceError_Unwrap(t *testing.T) {
	underlying := ErrResourceNotFound
	err := &ResourceError{
		URI: "file:///test",
		Op:  "read",
		Err: underlying,
	}

	if !errors.Is(err, underlying) {
		t.Error("errors.Is should return true for underlying error")
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != underlying {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlying)
	}
}
