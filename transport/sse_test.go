package transport

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestSSETransport_Name(t *testing.T) {
	transport := &SSETransport{}
	if got := transport.Name(); got != "sse" {
		t.Errorf("Name() = %q, want %q", got, "sse")
	}
}

func TestSSETransport_Info_BeforeServe(t *testing.T) {
	transport := &SSETransport{
		Config: SSEConfig{
			HTTPConfig: HTTPConfig{
				Host: "localhost",
				Port: 8080,
				Path: "/events",
			},
		},
	}

	info := transport.Info()
	if info.Name != "sse" {
		t.Errorf("Info().Name = %q, want %q", info.Name, "sse")
	}
	if info.Addr != "localhost:8080" {
		t.Errorf("Info().Addr = %q, want %q", info.Addr, "localhost:8080")
	}
	if info.Path != "/events" {
		t.Errorf("Info().Path = %q, want %q", info.Path, "/events")
	}
}

func TestSSETransport_Info_DefaultPath(t *testing.T) {
	transport := &SSETransport{
		Config: SSEConfig{
			HTTPConfig: HTTPConfig{
				Port: 8080,
			},
		},
	}

	info := transport.Info()
	if info.Path != "/mcp" {
		t.Errorf("Info().Path = %q, want default %q", info.Path, "/mcp")
	}
}

func TestSSETransport_Info_DefaultHost(t *testing.T) {
	transport := &SSETransport{
		Config: SSEConfig{
			HTTPConfig: HTTPConfig{
				Port: 8080,
			},
		},
	}

	info := transport.Info()
	if info.Addr != "0.0.0.0:8080" {
		t.Errorf("Info().Addr = %q, want %q", info.Addr, "0.0.0.0:8080")
	}
}

func TestSSETransport_Serve_BindsPort(t *testing.T) {
	// Find an available port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	transport := &SSETransport{
		Config: SSEConfig{
			HTTPConfig: HTTPConfig{
				Host:              "127.0.0.1",
				Port:              port,
				Path:              "/mcp",
				ReadHeaderTimeout: 5 * time.Second,
			},
		},
	}

	server := &httpTestServer{}
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- transport.Serve(ctx, server)
	}()

	// Wait for server to start
	time.Sleep(50 * time.Millisecond)

	// Verify port is bound
	info := transport.Info()
	if info.Addr == "" {
		t.Error("Info().Addr is empty after Serve started")
	}

	// Verify we can connect
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), time.Second)
	if err != nil {
		t.Errorf("failed to connect to transport: %v", err)
	} else {
		conn.Close()
	}

	cancel()
	<-errCh
}

func TestSSETransport_Serve_ContextCancellation(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	transport := &SSETransport{
		Config: SSEConfig{
			HTTPConfig: HTTPConfig{
				Host: "127.0.0.1",
				Port: port,
			},
		},
	}

	server := &httpTestServer{}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = transport.Serve(ctx, server)
	if err != nil {
		t.Errorf("Serve() error = %v, want nil on context cancellation", err)
	}
}

func TestSSETransport_Close_GracefulShutdown(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	transport := &SSETransport{
		Config: SSEConfig{
			HTTPConfig: HTTPConfig{
				Host: "127.0.0.1",
				Port: port,
			},
		},
	}

	server := &httpTestServer{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- transport.Serve(ctx, server)
	}()

	// Wait for server to start
	time.Sleep(50 * time.Millisecond)

	// Close should shut down gracefully
	if err := transport.Close(); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}

	// Serve should return
	select {
	case <-errCh:
		// Success
	case <-time.After(time.Second):
		t.Error("Serve did not return after Close")
	}
}

func TestSSETransport_Close_Idempotent(t *testing.T) {
	transport := &SSETransport{}

	// First close should succeed (even without Serve)
	if err := transport.Close(); err != nil {
		t.Errorf("first Close() error = %v, want nil", err)
	}

	// Second close should also succeed
	if err := transport.Close(); err != nil {
		t.Errorf("second Close() error = %v, want nil", err)
	}
}

func TestSSETransport_ConcurrentSafety(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	transport := &SSETransport{
		Config: SSEConfig{
			HTTPConfig: HTTPConfig{
				Host: "127.0.0.1",
				Port: port,
			},
		},
	}

	server := &httpTestServer{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = transport.Serve(ctx, server)
	}()

	time.Sleep(50 * time.Millisecond)

	// Concurrent Info and Close calls
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = transport.Info()
		}()
	}
	wg.Wait()

	transport.Close()
}

func TestSSETransport_ImplementsInterface(t *testing.T) {
	var _ Transport = (*SSETransport)(nil)
}

// httpTestServer is a test server for HTTP transports
type httpTestServer struct {
	handler http.Handler
}

func (s *httpTestServer) ServeTransport(ctx context.Context, transport Transport) error {
	<-ctx.Done()
	return nil
}

func (s *httpTestServer) Handler() http.Handler {
	if s.handler != nil {
		return s.handler
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
