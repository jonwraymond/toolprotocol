package discover

import (
	"testing"
)

func TestService_ID(t *testing.T) {
	svc := NewService("my-service", "http://localhost:8080")
	if svc.ID() != "my-service" {
		t.Errorf("ID() = %q, want %q", svc.ID(), "my-service")
	}
}

func TestService_Name(t *testing.T) {
	svc := NewService("test", "http://localhost:8080")
	svc.SetName("Test Service")

	if svc.Name() != "Test Service" {
		t.Errorf("Name() = %q, want %q", svc.Name(), "Test Service")
	}
}

func TestService_Name_Default(t *testing.T) {
	svc := NewService("my-id", "http://localhost:8080")

	// Name defaults to ID if not set
	if svc.Name() != "my-id" {
		t.Errorf("Name() = %q, want %q (default to ID)", svc.Name(), "my-id")
	}
}

func TestService_Description(t *testing.T) {
	svc := NewService("test", "http://localhost:8080")
	svc.SetDescription("A test service")

	if svc.Description() != "A test service" {
		t.Errorf("Description() = %q, want %q", svc.Description(), "A test service")
	}
}

func TestService_Version(t *testing.T) {
	svc := NewService("test", "http://localhost:8080")
	svc.SetVersion("1.0.0")

	if svc.Version() != "1.0.0" {
		t.Errorf("Version() = %q, want %q", svc.Version(), "1.0.0")
	}
}

func TestService_Endpoint(t *testing.T) {
	svc := NewService("test", "http://localhost:9000/mcp")
	if svc.Endpoint() != "http://localhost:9000/mcp" {
		t.Errorf("Endpoint() = %q, want %q", svc.Endpoint(), "http://localhost:9000/mcp")
	}
}

func TestService_Capabilities(t *testing.T) {
	svc := NewService("test", "http://localhost:8080")
	caps := &Capabilities{
		Tools:     true,
		Streaming: true,
	}
	svc.SetCapabilities(caps)

	got := svc.Capabilities()
	if got == nil {
		t.Fatal("Capabilities() returned nil")
	}
	if !got.Tools {
		t.Error("Capabilities().Tools = false, want true")
	}
}

func TestService_Capabilities_Default(t *testing.T) {
	svc := NewService("test", "http://localhost:8080")

	// Should return non-nil even if not set
	got := svc.Capabilities()
	if got == nil {
		t.Fatal("Capabilities() returned nil")
	}
}

func TestService_WithCapability(t *testing.T) {
	svc := NewService("test", "http://localhost:8080").
		WithCapability("tools").
		WithCapability("streaming")

	caps := svc.Capabilities()
	if !caps.Tools {
		t.Error("Tools = false, want true")
	}
	if !caps.Streaming {
		t.Error("Streaming = false, want true")
	}
}

func TestService_WithCapability_Unknown(t *testing.T) {
	svc := NewService("test", "http://localhost:8080").
		WithCapability("custom.extension")

	caps := svc.Capabilities()
	found := false
	for _, ext := range caps.Extensions {
		if ext == "custom.extension" {
			found = true
			break
		}
	}
	if !found {
		t.Error("custom.extension not found in Extensions")
	}
}

func TestService_Validate(t *testing.T) {
	svc := NewService("test", "http://localhost:8080")
	if err := svc.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestService_Validate_EmptyID(t *testing.T) {
	svc := NewService("", "http://localhost:8080")
	if err := svc.Validate(); err == nil {
		t.Error("Validate() with empty ID should return error")
	}
}

func TestService_Validate_EmptyEndpoint(t *testing.T) {
	svc := NewService("test", "")
	if err := svc.Validate(); err == nil {
		t.Error("Validate() with empty endpoint should return error")
	}
}

func TestService_ImplementsDiscoverable(t *testing.T) {
	var _ Discoverable = (*Service)(nil)
}
