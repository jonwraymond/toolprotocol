package transport

import "context"

// StdioTransport implements Transport using standard input/output.
//
// This transport is suitable for local process communication where
// the MCP server is invoked as a subprocess.
//
// StdioTransport has no configuration options and is stateless.
type StdioTransport struct{}

// Name returns "stdio" as the transport identifier.
func (t *StdioTransport) Name() string {
	return "stdio"
}

// Info returns descriptive information about the transport.
// For stdio, Addr and Path are empty since it's not a network transport.
func (t *StdioTransport) Info() Info {
	return Info{Name: "stdio"}
}

// Serve starts the stdio transport and blocks until context is cancelled.
// The server is responsible for handling stdin/stdout communication.
func (t *StdioTransport) Serve(ctx context.Context, server Server) error {
	return server.ServeTransport(ctx, t)
}

// Close is a no-op for stdio transport.
// Stdio handles are managed by the process, not the transport.
func (t *StdioTransport) Close() error {
	return nil
}
