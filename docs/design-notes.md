# Design Notes

## Architecture Decisions

1. **Layered protocols, not monoliths.** Transport, wire format, and content are
   isolated so that new protocols can reuse the same primitives.
2. **Protocol-agnostic primitives.** Packages like `task`, `stream`, `session`,
   `resource`, and `prompt` are independent of MCP/A2A/ACP specifics.
3. **Deterministic encoding.** `wire` implementations must encode/decode
   deterministically to preserve caching and reproducibility.
4. **Minimal dependencies.** The repo avoids heavy deps to keep transport
   implementations portable across environments.

## Contract Semantics

### transport
- **Concurrency:** implementations must be safe for concurrent use.
- **Cancellation:** `Serve` must respect context cancellation.
- **Idempotency:** `Close` is safe to call multiple times.

### wire
- **Lossless mapping:** encode/decode preserves request/response shape.
- **Capabilities:** `Capabilities()` must reflect actual encoder support.

### content
- **Immutability:** content instances are safe to share across goroutines.
- **MIME fidelity:** `MIMEType` must always match the payload.

### stream
- **Ordering:** events are delivered in send order.
- **Termination:** `Done()` closes after `Close()` completes.

### task
- **State machine:** only valid transitions are allowed.
- **Subscriptions:** updates are fan-out safe and non-blocking.

### session/resource/prompt/elicit
- **Thread safety:** registries are concurrency-safe.
- **Explicit errors:** missing keys or invalid args return typed errors.

## Trade-offs

- **Pure interfaces vs. convenience helpers.** We keep core interfaces small and
  provide optional helpers in each package to avoid bloated APIs.
- **Broad scope vs. cohesion.** This repo groups protocol primitives to simplify
  versioning, but avoids pulling in execution or schema concerns.
