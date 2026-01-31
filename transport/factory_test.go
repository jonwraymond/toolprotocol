package transport

import (
	"testing"
)

func TestNewTransport_Stdio(t *testing.T) {
	transport, err := New("stdio", nil)
	if err != nil {
		t.Fatalf("New(stdio) error = %v, want nil", err)
	}
	if transport == nil {
		t.Fatal("New(stdio) returned nil transport")
	}
	if transport.Name() != "stdio" {
		t.Errorf("transport.Name() = %q, want %q", transport.Name(), "stdio")
	}
}

func TestNewTransport_SSE(t *testing.T) {
	cfg := &SSEConfig{
		HTTPConfig: HTTPConfig{
			Host: "localhost",
			Port: 8080,
		},
	}
	transport, err := New("sse", cfg)
	if err != nil {
		t.Fatalf("New(sse) error = %v, want nil", err)
	}
	if transport == nil {
		t.Fatal("New(sse) returned nil transport")
	}
	if transport.Name() != "sse" {
		t.Errorf("transport.Name() = %q, want %q", transport.Name(), "sse")
	}
}

func TestNewTransport_SSE_NilConfig(t *testing.T) {
	transport, err := New("sse", nil)
	if err != nil {
		t.Fatalf("New(sse, nil) error = %v, want nil", err)
	}
	if transport.Name() != "sse" {
		t.Errorf("transport.Name() = %q, want %q", transport.Name(), "sse")
	}
}

func TestNewTransport_Streamable(t *testing.T) {
	cfg := &StreamableConfig{
		HTTPConfig: HTTPConfig{
			Host: "localhost",
			Port: 9000,
		},
		Stateless: true,
	}
	transport, err := New("streamable", cfg)
	if err != nil {
		t.Fatalf("New(streamable) error = %v, want nil", err)
	}
	if transport == nil {
		t.Fatal("New(streamable) returned nil transport")
	}
	if transport.Name() != "streamable" {
		t.Errorf("transport.Name() = %q, want %q", transport.Name(), "streamable")
	}
}

func TestNewTransport_Streamable_NilConfig(t *testing.T) {
	transport, err := New("streamable", nil)
	if err != nil {
		t.Fatalf("New(streamable, nil) error = %v, want nil", err)
	}
	if transport.Name() != "streamable" {
		t.Errorf("transport.Name() = %q, want %q", transport.Name(), "streamable")
	}
}

func TestNewTransport_Unknown(t *testing.T) {
	transport, err := New("unknown", nil)
	if err == nil {
		t.Error("New(unknown) error = nil, want error")
	}
	if transport != nil {
		t.Errorf("New(unknown) returned non-nil transport: %v", transport)
	}
}

func TestNewTransport_WrongConfigType(t *testing.T) {
	// Passing SSEConfig to streamable should use defaults, not error
	cfg := &SSEConfig{
		HTTPConfig: HTTPConfig{
			Port: 8080,
		},
	}
	transport, err := New("streamable", cfg)
	if err != nil {
		t.Fatalf("New(streamable, SSEConfig) error = %v, want nil", err)
	}
	if transport.Name() != "streamable" {
		t.Errorf("transport.Name() = %q, want %q", transport.Name(), "streamable")
	}
}

func TestRegistry_Register(t *testing.T) {
	reg := NewRegistry()
	factory := func(cfg any) (Transport, error) {
		return &StdioTransport{}, nil
	}

	reg.Register("custom", factory)

	transport, err := reg.New("custom", nil)
	if err != nil {
		t.Fatalf("New(custom) error = %v, want nil", err)
	}
	if transport.Name() != "stdio" {
		t.Errorf("transport.Name() = %q, want %q", transport.Name(), "stdio")
	}
}

func TestRegistry_Register_Overwrite(t *testing.T) {
	reg := NewRegistry()
	called := ""

	reg.Register("test", func(cfg any) (Transport, error) {
		called = "first"
		return &StdioTransport{}, nil
	})
	reg.Register("test", func(cfg any) (Transport, error) {
		called = "second"
		return &StdioTransport{}, nil
	})

	_, _ = reg.New("test", nil)
	if called != "second" {
		t.Errorf("called = %q, want %q (second registration should overwrite)", called, "second")
	}
}

func TestRegistry_Get(t *testing.T) {
	reg := NewRegistry()
	factory := func(cfg any) (Transport, error) {
		return &StdioTransport{}, nil
	}

	reg.Register("custom", factory)

	got := reg.Get("custom")
	if got == nil {
		t.Error("Get(custom) returned nil, want factory")
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	reg := NewRegistry()

	got := reg.Get("notfound")
	if got != nil {
		t.Errorf("Get(notfound) = %v, want nil", got)
	}
}

func TestRegistry_List(t *testing.T) {
	reg := NewRegistry()
	reg.Register("alpha", func(cfg any) (Transport, error) {
		return &StdioTransport{}, nil
	})
	reg.Register("beta", func(cfg any) (Transport, error) {
		return &StdioTransport{}, nil
	})

	list := reg.List()
	if len(list) != 2 {
		t.Errorf("len(List()) = %d, want 2", len(list))
	}

	// Check both are present (order not guaranteed)
	found := make(map[string]bool)
	for _, name := range list {
		found[name] = true
	}
	if !found["alpha"] {
		t.Error("List() missing 'alpha'")
	}
	if !found["beta"] {
		t.Error("List() missing 'beta'")
	}
}

func TestRegistry_New_NotFound(t *testing.T) {
	reg := NewRegistry()

	transport, err := reg.New("notfound", nil)
	if err == nil {
		t.Error("New(notfound) error = nil, want error")
	}
	if transport != nil {
		t.Errorf("New(notfound) returned non-nil: %v", transport)
	}
}

func TestDefaultRegistry(t *testing.T) {
	// Verify default registry has standard transports
	list := DefaultRegistry().List()
	found := make(map[string]bool)
	for _, name := range list {
		found[name] = true
	}

	if !found["stdio"] {
		t.Error("DefaultRegistry missing 'stdio'")
	}
	if !found["sse"] {
		t.Error("DefaultRegistry missing 'sse'")
	}
	if !found["streamable"] {
		t.Error("DefaultRegistry missing 'streamable'")
	}
}
