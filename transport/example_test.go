package transport_test

import (
	"context"
	"fmt"

	"github.com/jonwraymond/toolprotocol/transport"
)

func ExampleNew() {
	// Create a stdio transport using the default registry
	t, err := transport.New("stdio", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Transport name:", t.Name())
	// Output:
	// Transport name: stdio
}

func ExampleNew_streamable() {
	// Create a streamable HTTP transport with configuration
	cfg := &transport.StreamableConfig{
		HTTPConfig: transport.HTTPConfig{
			Host: "localhost",
			Port: 8080,
			Path: "/mcp",
		},
	}

	t, err := transport.New("streamable", cfg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Transport name:", t.Name())
	info := t.Info()
	fmt.Println("Path:", info.Path)
	// Output:
	// Transport name: streamable
	// Path: /mcp
}

func ExampleNew_sse() {
	// Create an SSE transport with configuration
	cfg := &transport.SSEConfig{
		HTTPConfig: transport.HTTPConfig{
			Host: "localhost",
			Port: 8081,
			Path: "/events",
		},
	}

	t, err := transport.New("sse", cfg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Transport name:", t.Name())
	info := t.Info()
	fmt.Println("Path:", info.Path)
	// Output:
	// Transport name: sse
	// Path: /events
}

func ExampleStdioTransport() {
	t := &transport.StdioTransport{}

	fmt.Println("Name:", t.Name())
	info := t.Info()
	fmt.Println("Info.Name:", info.Name)
	fmt.Println("Info.Addr:", info.Addr) // Empty for stdio
	// Output:
	// Name: stdio
	// Info.Name: stdio
	// Info.Addr:
}

func ExampleStdioTransport_Info() {
	t := &transport.StdioTransport{}
	info := t.Info()

	fmt.Println("Name:", info.Name)
	fmt.Println("Addr is empty:", info.Addr == "")
	fmt.Println("Path is empty:", info.Path == "")
	// Output:
	// Name: stdio
	// Addr is empty: true
	// Path is empty: true
}

func ExampleStdioTransport_Close() {
	t := &transport.StdioTransport{}

	// Close is idempotent and safe to call multiple times
	err := t.Close()
	fmt.Println("First close error:", err)

	err = t.Close()
	fmt.Println("Second close error:", err)
	// Output:
	// First close error: <nil>
	// Second close error: <nil>
}

func ExampleStreamableHTTPTransport() {
	t := &transport.StreamableHTTPTransport{
		Config: transport.StreamableConfig{
			HTTPConfig: transport.HTTPConfig{
				Host: "localhost",
				Port: 8080,
				Path: "/mcp",
			},
		},
	}

	fmt.Println("Name:", t.Name())
	info := t.Info()
	fmt.Println("Path:", info.Path)
	// Output:
	// Name: streamable
	// Path: /mcp
}

func ExampleStreamableHTTPTransport_Info() {
	t := &transport.StreamableHTTPTransport{
		Config: transport.StreamableConfig{
			HTTPConfig: transport.HTTPConfig{
				Host: "localhost",
				Port: 8080,
				Path: "/api/mcp",
			},
		},
	}

	info := t.Info()
	fmt.Println("Name:", info.Name)
	fmt.Println("Addr:", info.Addr)
	fmt.Println("Path:", info.Path)
	// Output:
	// Name: streamable
	// Addr: localhost:8080
	// Path: /api/mcp
}

func ExampleStreamableHTTPTransport_Info_defaults() {
	t := &transport.StreamableHTTPTransport{
		Config: transport.StreamableConfig{
			HTTPConfig: transport.HTTPConfig{
				Port: 3000,
			},
		},
	}

	info := t.Info()
	fmt.Println("Name:", info.Name)
	fmt.Println("Addr:", info.Addr)        // Uses 0.0.0.0 default
	fmt.Println("Path:", info.Path)        // Uses /mcp default
	// Output:
	// Name: streamable
	// Addr: 0.0.0.0:3000
	// Path: /mcp
}

func ExampleSSETransport() {
	t := &transport.SSETransport{
		Config: transport.SSEConfig{
			HTTPConfig: transport.HTTPConfig{
				Host: "localhost",
				Port: 8081,
				Path: "/events",
			},
		},
	}

	fmt.Println("Name:", t.Name())
	info := t.Info()
	fmt.Println("Path:", info.Path)
	// Output:
	// Name: sse
	// Path: /events
}

func ExampleSSETransport_Info() {
	t := &transport.SSETransport{
		Config: transport.SSEConfig{
			HTTPConfig: transport.HTTPConfig{
				Host: "127.0.0.1",
				Port: 9000,
				Path: "/sse",
			},
		},
	}

	info := t.Info()
	fmt.Println("Name:", info.Name)
	fmt.Println("Addr:", info.Addr)
	fmt.Println("Path:", info.Path)
	// Output:
	// Name: sse
	// Addr: 127.0.0.1:9000
	// Path: /sse
}

func ExampleNewRegistry() {
	reg := transport.NewRegistry()

	// Register a custom transport factory
	reg.Register("stdio", func(cfg any) (transport.Transport, error) {
		return &transport.StdioTransport{}, nil
	})

	// List registered transports
	names := reg.List()
	fmt.Println("Registered count:", len(names))
	// Output:
	// Registered count: 1
}

func ExampleRegistry_Register() {
	reg := transport.NewRegistry()

	// Register multiple transport factories
	reg.Register("stdio", func(cfg any) (transport.Transport, error) {
		return &transport.StdioTransport{}, nil
	})
	reg.Register("custom", func(cfg any) (transport.Transport, error) {
		return &transport.StdioTransport{}, nil // Example
	})

	fmt.Println("Registered count:", len(reg.List()))
	// Output:
	// Registered count: 2
}

func ExampleRegistry_Get() {
	reg := transport.NewRegistry()
	reg.Register("stdio", func(cfg any) (transport.Transport, error) {
		return &transport.StdioTransport{}, nil
	})

	// Get existing factory
	factory := reg.Get("stdio")
	fmt.Println("Found stdio:", factory != nil)

	// Get non-existent factory
	missing := reg.Get("unknown")
	fmt.Println("Found unknown:", missing != nil)
	// Output:
	// Found stdio: true
	// Found unknown: false
}

func ExampleRegistry_New() {
	reg := transport.NewRegistry()
	reg.Register("stdio", func(cfg any) (transport.Transport, error) {
		return &transport.StdioTransport{}, nil
	})

	// Create transport using registry
	t, err := reg.New("stdio", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Transport name:", t.Name())
	// Output:
	// Transport name: stdio
}

func ExampleRegistry_New_error() {
	reg := transport.NewRegistry()

	// Attempt to create unknown transport
	_, err := reg.New("unknown", nil)
	fmt.Println("Error:", err)
	// Output:
	// Error: unknown transport: unknown
}

func ExampleDefaultRegistry() {
	reg := transport.DefaultRegistry()

	// Default registry has stdio, sse, and streamable pre-registered
	names := reg.List()
	fmt.Println("Has transports:", len(names) >= 3)

	// Verify specific transports exist
	fmt.Println("Has stdio:", reg.Get("stdio") != nil)
	fmt.Println("Has sse:", reg.Get("sse") != nil)
	fmt.Println("Has streamable:", reg.Get("streamable") != nil)
	// Output:
	// Has transports: true
	// Has stdio: true
	// Has sse: true
	// Has streamable: true
}

func ExampleInfo() {
	info := transport.Info{
		Name: "streamable",
		Addr: "localhost:8080",
		Path: "/mcp",
	}

	fmt.Println("Name:", info.Name)
	fmt.Println("Addr:", info.Addr)
	fmt.Println("Path:", info.Path)
	// Output:
	// Name: streamable
	// Addr: localhost:8080
	// Path: /mcp
}

func ExampleHTTPConfig() {
	cfg := transport.HTTPConfig{
		Host: "0.0.0.0",
		Port: 8080,
		Path: "/mcp",
	}

	fmt.Println("Host:", cfg.Host)
	fmt.Println("Port:", cfg.Port)
	fmt.Println("Path:", cfg.Path)
	// Output:
	// Host: 0.0.0.0
	// Port: 8080
	// Path: /mcp
}

func ExampleStreamableConfig() {
	cfg := transport.StreamableConfig{
		HTTPConfig: transport.HTTPConfig{
			Host: "localhost",
			Port: 8080,
			Path: "/mcp",
		},
		Stateless:    false,
		JSONResponse: false,
	}

	fmt.Println("Stateless:", cfg.Stateless)
	fmt.Println("JSONResponse:", cfg.JSONResponse)
	// Output:
	// Stateless: false
	// JSONResponse: false
}

func ExampleStreamableConfig_withTLS() {
	cfg := transport.StreamableConfig{
		HTTPConfig: transport.HTTPConfig{
			Host: "localhost",
			Port: 443,
			Path: "/mcp",
		},
		TLS: transport.TLSConfig{
			Enabled:  true,
			CertFile: "/path/to/cert.pem",
			KeyFile:  "/path/to/key.pem",
		},
	}

	fmt.Println("TLS Enabled:", cfg.TLS.Enabled)
	fmt.Println("Port:", cfg.Port)
	// Output:
	// TLS Enabled: true
	// Port: 443
}

// mockServer implements transport.Server for examples
type mockServer struct{}

func (m *mockServer) ServeTransport(ctx context.Context, t transport.Transport) error {
	<-ctx.Done()
	return nil
}

func ExampleTransport_Serve() {
	t := &transport.StdioTransport{}
	server := &mockServer{}

	// Create a context that we cancel immediately for the example
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Serve would normally block until context is cancelled
	err := t.Serve(ctx, server)
	fmt.Println("Serve returned:", err)
	// Output:
	// Serve returned: <nil>
}
