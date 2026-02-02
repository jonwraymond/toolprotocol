# Examples

## Wire Envelope

```go
import "github.com/jonwraymond/toolprotocol/wire"

msg := wire.Envelope{
  Version: "1.0",
  Type:    "discover.request",
  Payload: map[string]any{
    "query": "create issue",
  },
}
```

## Stream

```go
import "github.com/jonwraymond/toolprotocol/stream"

s := stream.NewDefaultStream()

go func() {
  _ = s.Send(stream.Event{Type: "progress", Data: map[string]any{"pct": 50}})
  _ = s.Close()
}()

for ev := range s.Events() {
  _ = ev
}
```

## Session

```go
import "github.com/jonwraymond/toolprotocol/session"

sess := session.New(session.Config{ID: "session-1"})
ctx := sess.Context()
_ = ctx
```

## Content Blocks

```go
import "github.com/jonwraymond/toolprotocol/content"

block := content.Block{
  Type: "text",
  Data: map[string]any{"text": "hello"},
}
```
