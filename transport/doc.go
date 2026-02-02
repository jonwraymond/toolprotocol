// Package transport provides abstractions for MCP protocol transports.
//
// This package defines the Transport interface and common transport
// implementations including stdio, SSE, and streamable HTTP transports.
//
// # Ecosystem Position
//
// transport sits at the network boundary, providing the communication layer
// between external clients and internal protocol handlers:
//
//	┌─────────────────────────────────────────────────────────────────┐
//	│                    Transport Layer Architecture                 │
//	├─────────────────────────────────────────────────────────────────┤
//	│                                                                 │
//	│   External                transport              Server         │
//	│   ┌─────────┐          ┌───────────┐         ┌─────────┐       │
//	│   │ Client  │──────────│ Stdio     │─────────│ Server  │       │
//	│   │(process)│          │ Transport │         │  Impl   │       │
//	│   └─────────┘          └───────────┘         └─────────┘       │
//	│                                                    │            │
//	│   ┌─────────┐          ┌───────────┐              │            │
//	│   │ Client  │──────────│ Streamable│──────────────┘            │
//	│   │ (HTTP)  │          │ HTTP      │                           │
//	│   └─────────┘          └───────────┘                           │
//	│                                                                 │
//	│   ┌─────────┐          ┌───────────┐                           │
//	│   │ Client  │──────────│ SSE       │──────────────────────────┘│
//	│   │ (HTTP)  │          │ (legacy)  │                            │
//	│   └─────────┘          └───────────┘                           │
//	│                              │                                  │
//	│                         ┌────┴────┐                            │
//	│                         │Registry │                            │
//	│                         │(factory)│                            │
//	│                         └─────────┘                            │
//	│                                                                 │
//	└─────────────────────────────────────────────────────────────────┘
//
// # Core Components
//
//   - [Transport]: Interface for protocol communication mechanisms
//   - [Server]: Interface for handling transport requests
//   - [StdioTransport]: Standard I/O for subprocess communication
//   - [StreamableHTTPTransport]: Modern HTTP per MCP spec 2025-03-26
//   - [SSETransport]: Server-Sent Events (legacy, prefer Streamable)
//   - [Registry]: Thread-safe factory registry for transport creation
//   - [DefaultRegistry]: Pre-configured registry with all standard transports
//
// # Quick Start
//
//	// Create a transport using the factory
//	t, err := transport.New("stdio", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Start serving
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	err = t.Serve(ctx, myServer)
//
// # Available Transports
//
// Stdio Transport:
//   - Use case: Local subprocess communication
//   - Config: None required
//   - Concurrency: Stateless, safe for concurrent use
//
// Streamable HTTP Transport (recommended for HTTP):
//   - Use case: Network-based MCP servers
//   - Config: [StreamableConfig] with host, port, path, TLS
//   - Features: Session management, bidirectional, optional stateless mode
//   - Concurrency: Safe for concurrent use via sync.Mutex
//
// SSE Transport (legacy):
//   - Use case: Legacy HTTP clients
//   - Config: [SSEConfig] with host, port, path
//   - Note: Prefer StreamableHTTPTransport for new implementations
//   - Concurrency: Safe for concurrent use via sync.Mutex
//
// # Thread Safety
//
// All exported types are safe for concurrent use:
//
//   - [StdioTransport]: Stateless, concurrent-safe
//   - [StreamableHTTPTransport]: sync.Mutex protects listener/server state
//   - [SSETransport]: sync.Mutex protects listener/server state
//   - [Registry]: sync.RWMutex protects all operations
//
// # Error Handling
//
// Sentinel errors (use errors.Is for checking):
//
//   - [ErrTransportClosed]: Operations on closed transport
//   - [ErrAlreadyServing]: Serve called on active transport
//   - [ErrInvalidConfig]: Invalid configuration provided
//
// Transport operations wrap underlying errors with context:
//
//	err := transport.Serve(ctx, server)
//	// err may contain: "listen localhost:8080: <underlying error>"
//
// # Graceful Shutdown
//
// HTTP transports support graceful shutdown with a 5-second timeout:
//
//	ctx, cancel := context.WithCancel(context.Background())
//
//	go func() {
//	    time.Sleep(10 * time.Second)
//	    cancel() // Triggers graceful shutdown
//	}()
//
//	err := transport.Serve(ctx, server) // Blocks until ctx cancelled
//
// Close() is idempotent and safe to call multiple times.
//
// # TLS Configuration
//
// StreamableHTTPTransport supports TLS for secure communication:
//
//	cfg := &transport.StreamableConfig{
//	    HTTPConfig: transport.HTTPConfig{Port: 443},
//	    TLS: transport.TLSConfig{
//	        Enabled:  true,
//	        CertFile: "/path/to/cert.pem",
//	        KeyFile:  "/path/to/key.pem",
//	    },
//	}
//	t, _ := transport.New("streamable", cfg)
//
// TLS 1.2 is the minimum supported version.
//
// # Integration with ApertureStack
//
// transport integrates with other ApertureStack packages:
//
//   - wire: Encodes/decodes messages for the transport layer
//   - session: Manages client sessions for stateful transports
//   - stream: Provides streaming event delivery over transports
package transport
