// Package transport provides abstractions for MCP protocol transports.
//
// This package defines the Transport interface and common transport
// implementations including stdio, SSE, and streamable HTTP transports.
//
// # Transport Interface
//
// The Transport interface abstracts the underlying communication mechanism:
//
//	type Transport interface {
//	    Name() string
//	    Info() Info
//	    Serve(ctx context.Context, server Server) error
//	    Close() error
//	}
//
// # Available Transports
//
//   - StdioTransport: Standard input/output for local process communication
//   - SSETransport: Server-Sent Events over HTTP (legacy)
//   - StreamableHTTPTransport: Modern HTTP transport per MCP spec 2025-03-26
//
// # Usage
//
//	// Create a transport
//	transport := &transport.StdioTransport{}
//
//	// Or use the factory
//	transport, err := transport.New("stdio", nil)
//
//	// Start serving
//	err := transport.Serve(ctx, server)
package transport
