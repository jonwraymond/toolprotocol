// Package wire provides multi-protocol wire format adapters for tool communication.
//
// It enables encoding and decoding of tool requests and responses across
// different protocols including MCP (Anthropic), A2A (Google), and ACP (IBM).
//
// # Ecosystem Position
//
// wire sits at the protocol boundary, translating between internal representations
// and protocol-specific wire formats:
//
//	┌─────────────────────────────────────────────────────────────────┐
//	│                    Protocol Translation Flow                    │
//	├─────────────────────────────────────────────────────────────────┤
//	│                                                                 │
//	│   External Client          wire               Internal          │
//	│   ┌───────────┐         ┌─────────┐         ┌─────────┐        │
//	│   │ MCP/A2A/  │─────────│ Decode  │─────────│ Request │        │
//	│   │ ACP JSON  │         │         │         │ struct  │        │
//	│   └───────────┘         │ ┌─────┐ │         └─────────┘        │
//	│        ▲                │ │Wire │ │              │              │
//	│        │                │ │ Impl│ │              ▼              │
//	│        │                │ └─────┘ │         ┌─────────┐        │
//	│   ┌───────────┐         │         │         │ Execute │        │
//	│   │ MCP/A2A/  │◀────────│ Encode  │◀────────│  Tool   │        │
//	│   │ ACP JSON  │         │         │         │         │        │
//	│   └───────────┘         └─────────┘         └─────────┘        │
//	│                              │                                  │
//	│                         ┌────┴────┐                            │
//	│                         │Registry │                            │
//	│                         │ (lookup)│                            │
//	│                         └─────────┘                            │
//	│                                                                 │
//	└─────────────────────────────────────────────────────────────────┘
//
// # Core Components
//
//   - [Wire]: Interface for protocol-specific encoding/decoding
//   - [MCPWire]: Model Context Protocol (Anthropic) - JSON-RPC 2.0 based
//   - [A2AWire]: Agent-to-Agent Protocol (Google) - JSON-RPC with artifacts
//   - [ACPWire]: Agent Communication Protocol (IBM) - JSON-RPC with agents
//   - [Registry]: Thread-safe registry of wire format handlers
//   - [DefaultRegistry]: Pre-configured registry with all standard formats
//
// # Quick Start
//
//	// Use default registry for standard protocols
//	reg := wire.DefaultRegistry()
//	w := reg.Get("mcp")
//
//	// Encode a request
//	req := &wire.Request{
//	    ID:     "1",
//	    Method: "tools/call",
//	    ToolID: "search",
//	    Arguments: map[string]any{"query": "golang"},
//	}
//	data, err := w.EncodeRequest(ctx, req)
//
//	// Decode a response
//	resp, err := w.DecodeResponse(ctx, responseData)
//
// # Available Formats
//
// MCP (Model Context Protocol):
//   - Version: 2025-11-25
//   - Streaming: Yes
//   - Batch requests: No
//   - Progress notifications: Yes
//   - Cancellation: Yes
//
// A2A (Agent-to-Agent Protocol):
//   - Version: 0.2.1
//   - Streaming: Yes
//   - Batch requests: No
//   - Progress notifications: Yes
//   - Cancellation: Yes
//
// ACP (Agent Communication Protocol):
//   - Version: 1.0.0
//   - Streaming: No
//   - Batch requests: Yes
//   - Progress notifications: No
//   - Cancellation: Yes
//
// # Thread Safety
//
// All exported types are safe for concurrent use:
//
//   - [MCPWire], [A2AWire], [ACPWire]: Stateless, concurrent-safe
//   - [Registry]: sync.RWMutex protects all operations
//   - [DefaultRegistry]: Returns shared instance, safe to use concurrently
//
// # Error Handling
//
// Sentinel errors (use errors.Is for checking):
//
//   - [ErrUnsupportedFormat]: Unknown wire format requested
//   - [ErrEncodeFailure]: Encoding to wire format failed
//   - [ErrDecodeFailure]: Decoding from wire format failed
//
// Encode/Decode methods wrap underlying errors with context:
//
//	resp, err := w.DecodeResponse(ctx, data)
//	if err != nil {
//	    // err contains: "decode response: <underlying error>"
//	}
//
// # Integration with ApertureStack
//
// wire integrates with other ApertureStack packages:
//
//   - transport: Uses wire for protocol-specific message encoding
//   - stream: Streaming responses use wire for event encoding
//   - discover: Tool lists encoded via EncodeToolList/DecodeToolList
//   - content: Response content types map to wire.Content
package wire
