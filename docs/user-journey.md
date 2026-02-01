# User Journey

## Installation

```bash
go get github.com/jonwraymond/toolprotocol@latest
```

## Basic Usage: Minimal Protocol Server

```go
import (
  "context"

  "github.com/jonwraymond/toolprotocol/transport"
  "github.com/jonwraymond/toolprotocol/wire"
)

type server struct{
  codec wire.Wire
}

func (s *server) ServeTransport(ctx context.Context, t transport.Transport) error {
  // Read bytes from transport, decode with s.codec, route to handlers.
  return nil
}

ctx := context.Background()
srv := &server{codec: wire.NewMCP()}
tp, _ := transport.New("stdio", nil)
_ = tp.Serve(ctx, srv)
```

## Intermediate: Sessions + Resources

```go
import (
  "context"

  "github.com/jonwraymond/toolprotocol/resource"
  "github.com/jonwraymond/toolprotocol/session"
)

ctx := context.Background()

store := session.NewMemoryStore()
sess, _ := store.Create(ctx, "client-1")
ctx = session.WithSession(ctx, sess)

registry := resource.NewRegistry()
static := resource.NewStaticProvider()
static.Add(
  resource.Resource{URI: "file:///readme.md", Name: "README", MIMEType: "text/markdown"},
  resource.Contents{URI: "file:///readme.md", MIMEType: "text/markdown", Text: "# README"},
)
registry.Register("file", static)

_, _ = registry.Read(ctx, "file:///readme.md")
```

## Advanced: Task + Streaming

```go
import (
  "context"

  "github.com/jonwraymond/toolprotocol/stream"
  "github.com/jonwraymond/toolprotocol/task"
)

ctx := context.Background()
mgr := task.NewManager()
_ , _ = mgr.Create(ctx, "job-123")

source := stream.NewSource()
s := source.NewBufferedStream(ctx, 100)
_ = s.Send(ctx, stream.Event{Type: stream.EventProgress, Data: 0.3})
_ = s.Send(ctx, stream.Event{Type: stream.EventComplete, Data: map[string]any{"ok": true}})
_ = s.Close()
```
