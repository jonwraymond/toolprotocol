package transport

import "time"

// HTTPConfig holds common configuration for HTTP-based transports.
type HTTPConfig struct {
	// Host is the network interface to bind (default: "0.0.0.0").
	Host string

	// Port is the TCP port to listen on.
	Port int

	// Path is the HTTP endpoint path (default: "/mcp").
	Path string

	// ReadHeaderTimeout limits how long to wait for request headers.
	// Prevents slowloris attacks.
	ReadHeaderTimeout time.Duration
}

// TLSConfig holds TLS/HTTPS configuration for secure transport.
//
// When Enabled is true, the transport serves HTTPS using the specified
// certificate and key files. TLS 1.2 is the minimum supported version.
type TLSConfig struct {
	// Enabled activates TLS encryption for the transport.
	Enabled bool

	// CertFile is the path to the PEM-encoded certificate file.
	CertFile string

	// KeyFile is the path to the PEM-encoded private key file.
	KeyFile string
}

// SSEConfig holds configuration for the SSE transport.
type SSEConfig struct {
	HTTPConfig
}

// StreamableConfig holds configuration for the Streamable HTTP transport.
//
// Streamable HTTP is the recommended HTTP transport per MCP spec 2025-03-26,
// replacing the deprecated SSE transport.
type StreamableConfig struct {
	HTTPConfig

	// TLS enables HTTPS with certificate-based encryption.
	TLS TLSConfig

	// Stateless disables session management when true.
	// In stateless mode, no Mcp-Session-Id validation occurs.
	Stateless bool

	// JSONResponse causes responses to use application/json instead of
	// text/event-stream (SSE).
	JSONResponse bool

	// SessionTimeout configures idle session cleanup duration.
	// Sessions with no HTTP activity for this duration are closed.
	SessionTimeout time.Duration
}
