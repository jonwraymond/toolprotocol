package transport

import (
	"context"
	"testing"
)

func TestInfo_Fields(t *testing.T) {
	info := Info{
		Name: "test-transport",
		Addr: "localhost:8080",
		Path: "/mcp",
	}

	if info.Name != "test-transport" {
		t.Errorf("Name = %q, want %q", info.Name, "test-transport")
	}
	if info.Addr != "localhost:8080" {
		t.Errorf("Addr = %q, want %q", info.Addr, "localhost:8080")
	}
	if info.Path != "/mcp" {
		t.Errorf("Path = %q, want %q", info.Path, "/mcp")
	}
}

func TestInfo_Empty(t *testing.T) {
	var info Info
	if info.Name != "" {
		t.Errorf("empty Info.Name = %q, want empty", info.Name)
	}
	if info.Addr != "" {
		t.Errorf("empty Info.Addr = %q, want empty", info.Addr)
	}
	if info.Path != "" {
		t.Errorf("empty Info.Path = %q, want empty", info.Path)
	}
}

func TestTransportInterface_Contract(t *testing.T) {
	// Verify that Transport interface has the expected methods
	// This is a compile-time check via type assertion
	var _ Transport = (*mockTransport)(nil)
}

func TestServerInterface_Contract(t *testing.T) {
	// Verify that Server interface has the expected methods
	var _ Server = (*mockServer)(nil)
}

// mockTransport implements Transport for testing interface contracts
type mockTransport struct {
	name   string
	info   Info
	closed bool
}

func (m *mockTransport) Name() string {
	return m.name
}

func (m *mockTransport) Info() Info {
	return m.info
}

func (m *mockTransport) Serve(ctx context.Context, server Server) error {
	<-ctx.Done()
	return ctx.Err()
}

func (m *mockTransport) Close() error {
	m.closed = true
	return nil
}

// mockServer implements Server for testing interface contracts
type mockServer struct{}

func (m *mockServer) ServeTransport(ctx context.Context, transport Transport) error {
	return nil
}
