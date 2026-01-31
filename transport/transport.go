package transport

import "context"

// Info describes a transport instance.
type Info struct {
	// Name is the transport type identifier (e.g., "stdio", "sse", "streamable").
	Name string

	// Addr is the network address for HTTP transports (e.g., "localhost:8080").
	// Empty for non-network transports like stdio.
	Addr string

	// Path is the HTTP endpoint path (e.g., "/mcp").
	// Empty for non-HTTP transports.
	Path string
}

// Transport defines the interface for MCP protocol transports.
//
// Contract:
//   - Concurrency: implementations must be safe for concurrent use.
//   - Context: Serve must honor cancellation/deadlines.
//   - Errors: Close must be idempotent.
type Transport interface {
	// Name returns the transport type identifier.
	Name() string

	// Info returns descriptive information about the transport.
	Info() Info

	// Serve starts the transport and blocks until context is cancelled
	// or an unrecoverable error occurs.
	Serve(ctx context.Context, server Server) error

	// Close gracefully shuts down the transport.
	// Close is idempotent and safe to call multiple times.
	Close() error
}

// Server provides the contract needed by transports to serve requests.
//
// Contract:
//   - Concurrency: implementations must be safe for concurrent use.
//   - Context: ServeTransport must honor cancellation/deadlines.
type Server interface {
	// ServeTransport handles requests from the given transport.
	ServeTransport(ctx context.Context, transport Transport) error
}
