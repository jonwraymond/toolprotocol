# toolprotocol

Protocol layer providing transport, wire format, and protocol primitives for MCP,
A2A, and ACP integrations. This repo is pure Go: no network I/O is performed
outside the transport implementations.

Part of the [ApertureStack](https://github.com/jonwraymond) ecosystem.

## Installation

```bash
go get github.com/jonwraymond/toolprotocol@latest
```

## Packages

| Package | Purpose |
|---------|---------|
| `content` | Unified content parts (text, image, audio, file, resource) |
| `discover` | Service discovery + capability negotiation |
| `transport` | Transport interfaces (stdio, SSE, streamable HTTP) |
| `wire` | Protocol wire encoding (MCP, A2A, ACP) |
| `stream` | Streaming events for progress/partial/complete |
| `session` | Client session store + context helpers |
| `task` | Long-running task lifecycle + subscriptions |
| `resource` | MCP resources registry + subscriptions |
| `prompt` | Prompt templates + registry |
| `elicit` | User input elicitation (text/confirm/choice/form) |

## Quick Start

```go
import (
  "context"

  "github.com/jonwraymond/toolprotocol/content"
  "github.com/jonwraymond/toolprotocol/transport"
  "github.com/jonwraymond/toolprotocol/wire"
)

type server struct{}

func (s *server) ServeTransport(ctx context.Context, t transport.Transport) error {
  return nil
}

ctx := context.Background()

// Build content parts
builder := content.NewBuilder()
parts := []content.Content{
  builder.Text("hello"),
}

// Encode a request using MCP wire format
codec := wire.NewMCP()
payload, err := codec.EncodeRequest(ctx, &wire.Request{
  ID:        "1",
  Method:    "tools/call",
  ToolID:    "echo",
  Arguments: map[string]any{"message": "hello"},
  Meta:      map[string]any{"content": parts},
})
_ = payload
_ = err

// Serve over a transport
tp, _ := transport.New("stdio", nil)
_ = tp.Serve(ctx, &server{})
```

## Docs

See the [docs](./docs/) directory and the aggregated docs site in
`ai-tools-stack`.

## License

MIT License - see [LICENSE](./LICENSE)
