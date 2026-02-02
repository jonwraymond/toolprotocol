package transport

import (
	"errors"
	"testing"
)

func TestErrTransportClosed(t *testing.T) {
	if ErrTransportClosed == nil {
		t.Fatal("ErrTransportClosed is nil")
	}
	if ErrTransportClosed.Error() != "transport: closed" {
		t.Errorf("ErrTransportClosed.Error() = %q, want %q",
			ErrTransportClosed.Error(), "transport: closed")
	}
}

func TestErrAlreadyServing(t *testing.T) {
	if ErrAlreadyServing == nil {
		t.Fatal("ErrAlreadyServing is nil")
	}
	if ErrAlreadyServing.Error() != "transport: already serving" {
		t.Errorf("ErrAlreadyServing.Error() = %q, want %q",
			ErrAlreadyServing.Error(), "transport: already serving")
	}
}

func TestErrInvalidConfig(t *testing.T) {
	if ErrInvalidConfig == nil {
		t.Fatal("ErrInvalidConfig is nil")
	}
	if ErrInvalidConfig.Error() != "transport: invalid configuration" {
		t.Errorf("ErrInvalidConfig.Error() = %q, want %q",
			ErrInvalidConfig.Error(), "transport: invalid configuration")
	}
}

func TestErrors_AreDistinct(t *testing.T) {
	if errors.Is(ErrTransportClosed, ErrAlreadyServing) {
		t.Error("ErrTransportClosed should not match ErrAlreadyServing")
	}
	if errors.Is(ErrTransportClosed, ErrInvalidConfig) {
		t.Error("ErrTransportClosed should not match ErrInvalidConfig")
	}
	if errors.Is(ErrAlreadyServing, ErrInvalidConfig) {
		t.Error("ErrAlreadyServing should not match ErrInvalidConfig")
	}
}
