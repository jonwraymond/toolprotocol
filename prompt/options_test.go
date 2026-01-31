package prompt

import (
	"testing"
)

func TestWithValidator(t *testing.T) {
	validator := func(args map[string]string) error {
		return nil
	}

	opt := WithValidator(validator)
	if opt == nil {
		t.Error("WithValidator() returned nil")
	}

	// Call the option to ensure it doesn't panic
	opt(nil)
}

func TestNewRegistry_Defaults(t *testing.T) {
	r := NewRegistry()

	if r.prompts == nil {
		t.Error("prompts map not initialized")
	}
	if r.handlers == nil {
		t.Error("handlers map not initialized")
	}
}
