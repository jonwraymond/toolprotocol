# toolprotocol

Protocol layer providing transport, wire format, and protocol primitives for MCP,
A2A, and ACP integrations. The packages here are transport-agnostic and composable.

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

## Installation

```bash
go get github.com/jonwraymond/toolprotocol@latest
```

## Quick Start: Wire + Transport

```go
import (
  "context"

  "github.com/jonwraymond/toolprotocol/transport"
  "github.com/jonwraymond/toolprotocol/wire"
)

type server struct{}

func (s *server) ServeTransport(ctx context.Context, t transport.Transport) error {
  // decode/route using wire codecs
  return nil
}

ctx := context.Background()
codec := wire.NewMCP()
payload, _ := codec.EncodeRequest(ctx, &wire.Request{
  ID:        "1",
  Method:    "tools/list",
  ToolID:    "",
  Arguments: nil,
})

tp, _ := transport.New("stdio", nil)
_ = payload
_ = tp.Serve(ctx, &server{})
```

## Quick Start: Tasks + Streaming

```go
import (
  "context"

  "github.com/jonwraymond/toolprotocol/stream"
  "github.com/jonwraymond/toolprotocol/task"
)

ctx := context.Background()
mgr := task.NewManager()
_ , _ = mgr.Create(ctx, "task-1")

source := stream.NewSource()
s := source.NewBufferedStream(ctx, 50)
_ = s.Send(ctx, stream.Event{Type: stream.EventProgress, Data: 0.5})
_ = s.Send(ctx, stream.Event{Type: stream.EventComplete, Data: map[string]any{"ok": true}})
_ = s.Close()
```

## Contracts (Summary)

- **Transport**: concurrent-safe; `Serve` honors context; `Close` is idempotent.
- **Wire**: encode/decode deterministic; `Capabilities` must match actual behavior.
- **Content**: immutable content instances; `Bytes` must be safe to call multiple times.
- **Stream**: event ordering preserved; `Done` closes after `Close`.
- **Task**: state machine enforces valid transitions; terminal states are final.
- **Session**: store guarantees TTL cleanup and thread safety.
- **Resource/Prompt/Elicit**: registries are concurrency-safe; errors are explicit.
