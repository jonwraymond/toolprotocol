package wire

import "context"

// Wire encodes/decodes protocol-specific wire formats.
type Wire interface {
	// Name returns the protocol name (e.g., "mcp", "a2a", "acp").
	Name() string

	// Version returns the protocol version.
	Version() string

	// EncodeRequest encodes a request to wire format.
	EncodeRequest(ctx context.Context, req *Request) ([]byte, error)

	// DecodeRequest decodes a request from wire format.
	DecodeRequest(ctx context.Context, data []byte) (*Request, error)

	// EncodeResponse encodes a response to wire format.
	EncodeResponse(ctx context.Context, resp *Response) ([]byte, error)

	// DecodeResponse decodes a response from wire format.
	DecodeResponse(ctx context.Context, data []byte) (*Response, error)

	// EncodeToolList encodes a list of tools to wire format.
	EncodeToolList(ctx context.Context, tools []Tool) ([]byte, error)

	// DecodeToolList decodes a list of tools from wire format.
	DecodeToolList(ctx context.Context, data []byte) ([]Tool, error)

	// Capabilities returns the protocol capabilities.
	Capabilities() *Capabilities
}

// Request represents a tool invocation request.
type Request struct {
	// ID is the request identifier.
	ID string

	// Method is the RPC method (e.g., "tools/call", "tools/list").
	Method string

	// ToolID is the tool to invoke.
	ToolID string

	// Arguments are the tool input parameters.
	Arguments map[string]any

	// Meta contains protocol-specific metadata.
	Meta map[string]any
}

// Response represents a tool invocation response.
type Response struct {
	// ID is the request identifier this responds to.
	ID string

	// Content is the response payload.
	Content []Content

	// IsError indicates if this is an error response.
	IsError bool

	// Error contains error details when IsError is true.
	Error *Error

	// Meta contains protocol-specific metadata.
	Meta map[string]any
}

// Capabilities describes protocol features.
type Capabilities struct {
	// Streaming indicates support for streaming responses.
	Streaming bool

	// BatchRequests indicates support for batched requests.
	BatchRequests bool

	// Progress indicates support for progress notifications.
	Progress bool

	// Cancellation indicates support for request cancellation.
	Cancellation bool
}
