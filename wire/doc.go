// Package wire provides multi-protocol wire format adapters for MCP tool communication.
//
// This package enables encoding and decoding of tool requests and responses
// across different protocols including MCP, A2A (Google), and ACP (IBM).
//
// # Wire Interface
//
// The Wire interface abstracts protocol-specific encoding:
//
//	type Wire interface {
//	    Name() string
//	    Version() string
//	    EncodeRequest(ctx context.Context, req *Request) ([]byte, error)
//	    DecodeRequest(ctx context.Context, data []byte) (*Request, error)
//	    EncodeResponse(ctx context.Context, resp *Response) ([]byte, error)
//	    DecodeResponse(ctx context.Context, data []byte) (*Response, error)
//	    Capabilities() *Capabilities
//	}
//
// # Available Formats
//
//   - MCPWire: Model Context Protocol (JSON-RPC based)
//   - A2AWire: Agent-to-Agent Protocol (Google)
//   - ACPWire: Agent Communication Protocol (IBM)
//
// # Usage
//
//	wire := wire.NewMCP()
//	encoded, err := wire.EncodeRequest(ctx, &wire.Request{
//	    Method: "tools/call",
//	    ToolID: "search",
//	    Arguments: map[string]any{"query": "test"},
//	})
package wire
