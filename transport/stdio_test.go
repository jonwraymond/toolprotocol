package transport

import (
	"context"
	"testing"
	"time"
)

func TestStdioTransport_Name(t *testing.T) {
	transport := &StdioTransport{}
	if got := transport.Name(); got != "stdio" {
		t.Errorf("Name() = %q, want %q", got, "stdio")
	}
}

func TestStdioTransport_Info(t *testing.T) {
	transport := &StdioTransport{}
	info := transport.Info()

	if info.Name != "stdio" {
		t.Errorf("Info().Name = %q, want %q", info.Name, "stdio")
	}
	if info.Addr != "" {
		t.Errorf("Info().Addr = %q, want empty for stdio", info.Addr)
	}
	if info.Path != "" {
		t.Errorf("Info().Path = %q, want empty for stdio", info.Path)
	}
}

func TestStdioTransport_Close_Idempotent(t *testing.T) {
	transport := &StdioTransport{}

	// First close should succeed
	if err := transport.Close(); err != nil {
		t.Errorf("first Close() error = %v, want nil", err)
	}

	// Second close should also succeed (idempotent)
	if err := transport.Close(); err != nil {
		t.Errorf("second Close() error = %v, want nil", err)
	}
}

func TestStdioTransport_Serve_ContextCancellation(t *testing.T) {
	transport := &StdioTransport{}
	server := &testServer{
		serveFunc: func(ctx context.Context, _ Transport) error {
			<-ctx.Done()
			return ctx.Err()
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := transport.Serve(ctx, server)
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("Serve() error = %v, want nil or DeadlineExceeded", err)
	}
}

func TestStdioTransport_ImplementsInterface(t *testing.T) {
	var _ Transport = (*StdioTransport)(nil)
}

// testServer is a helper for testing transports
type testServer struct {
	serveFunc func(ctx context.Context, transport Transport) error
}

func (s *testServer) ServeTransport(ctx context.Context, transport Transport) error {
	if s.serveFunc != nil {
		return s.serveFunc(ctx, transport)
	}
	<-ctx.Done()
	return ctx.Err()
}
