package transport

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

func TestStreamableTransport_Name(t *testing.T) {
	transport := &StreamableHTTPTransport{}
	if got := transport.Name(); got != "streamable" {
		t.Errorf("Name() = %q, want %q", got, "streamable")
	}
}

func TestStreamableTransport_Info_BeforeServe(t *testing.T) {
	transport := &StreamableHTTPTransport{
		Config: StreamableConfig{
			HTTPConfig: HTTPConfig{
				Host: "localhost",
				Port: 9000,
				Path: "/api/mcp",
			},
		},
	}

	info := transport.Info()
	if info.Name != "streamable" {
		t.Errorf("Info().Name = %q, want %q", info.Name, "streamable")
	}
	if info.Addr != "localhost:9000" {
		t.Errorf("Info().Addr = %q, want %q", info.Addr, "localhost:9000")
	}
	if info.Path != "/api/mcp" {
		t.Errorf("Info().Path = %q, want %q", info.Path, "/api/mcp")
	}
}

func TestStreamableTransport_Info_Defaults(t *testing.T) {
	transport := &StreamableHTTPTransport{
		Config: StreamableConfig{
			HTTPConfig: HTTPConfig{
				Port: 8080,
			},
		},
	}

	info := transport.Info()
	if info.Path != "/mcp" {
		t.Errorf("Info().Path = %q, want default %q", info.Path, "/mcp")
	}
	if info.Addr != "0.0.0.0:8080" {
		t.Errorf("Info().Addr = %q, want %q", info.Addr, "0.0.0.0:8080")
	}
}

func TestStreamableTransport_Serve_BindsPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatalf("failed to close listener: %v", err)
	}

	transport := &StreamableHTTPTransport{
		Config: StreamableConfig{
			HTTPConfig: HTTPConfig{
				Host:              "127.0.0.1",
				Port:              port,
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
		if err := conn.Close(); err != nil {
			t.Errorf("failed to close connection: %v", err)
		}
	}

	cancel()
	<-errCh
}

func TestStreamableTransport_Serve_ContextCancellation(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatalf("failed to close listener: %v", err)
	}

	transport := &StreamableHTTPTransport{
		Config: StreamableConfig{
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

func TestStreamableTransport_Serve_StatelessMode(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatalf("failed to close listener: %v", err)
	}

	transport := &StreamableHTTPTransport{
		Config: StreamableConfig{
			HTTPConfig: HTTPConfig{
				Host: "127.0.0.1",
				Port: port,
			},
			Stateless: true,
		},
	}

	server := &httpTestServer{}
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- transport.Serve(ctx, server)
	}()

	time.Sleep(50 * time.Millisecond)

	// Just verify it starts successfully
	cancel()
	<-errCh
}

func TestStreamableTransport_Serve_JSONMode(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatalf("failed to close listener: %v", err)
	}

	transport := &StreamableHTTPTransport{
		Config: StreamableConfig{
			HTTPConfig: HTTPConfig{
				Host: "127.0.0.1",
				Port: port,
			},
			JSONResponse: true,
		},
	}

	server := &httpTestServer{}
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- transport.Serve(ctx, server)
	}()

	time.Sleep(50 * time.Millisecond)

	// Just verify it starts successfully
	cancel()
	<-errCh
}

func TestStreamableTransport_SessionTimeout(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatalf("failed to close listener: %v", err)
	}

	transport := &StreamableHTTPTransport{
		Config: StreamableConfig{
			HTTPConfig: HTTPConfig{
				Host: "127.0.0.1",
				Port: port,
			},
			SessionTimeout: 30 * time.Second,
		},
	}

	server := &httpTestServer{}
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- transport.Serve(ctx, server)
	}()

	time.Sleep(50 * time.Millisecond)

	// Just verify it starts successfully with session timeout configured
	cancel()
	<-errCh
}

func TestStreamableTransport_Close_GracefulShutdown(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatalf("failed to close listener: %v", err)
	}

	transport := &StreamableHTTPTransport{
		Config: StreamableConfig{
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

	time.Sleep(50 * time.Millisecond)

	// Close should shut down gracefully
	if err := transport.Close(); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}

	select {
	case <-errCh:
		// Success
	case <-time.After(time.Second):
		t.Error("Serve did not return after Close")
	}
}

func TestStreamableTransport_Close_BeforeServe(t *testing.T) {
	transport := &StreamableHTTPTransport{}

	// Close before Serve should be safe
	if err := transport.Close(); err != nil {
		t.Errorf("Close() before Serve error = %v, want nil", err)
	}
}

func TestStreamableTransport_Close_Idempotent(t *testing.T) {
	transport := &StreamableHTTPTransport{}

	if err := transport.Close(); err != nil {
		t.Errorf("first Close() error = %v, want nil", err)
	}
	if err := transport.Close(); err != nil {
		t.Errorf("second Close() error = %v, want nil", err)
	}
}

func TestStreamableTransport_ConcurrentSafety(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatalf("failed to close listener: %v", err)
	}

	transport := &StreamableHTTPTransport{
		Config: StreamableConfig{
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

	// Concurrent Info calls
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = transport.Info()
		}()
	}
	wg.Wait()

	if err := transport.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestStreamableTransport_ImplementsInterface(t *testing.T) {
	var _ Transport = (*StreamableHTTPTransport)(nil)
}
