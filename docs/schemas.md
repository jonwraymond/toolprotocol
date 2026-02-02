# Schemas and Contracts

This document defines the protocol schema contracts in `toolprotocol`.
These are **structural** schemas (wire formats, streaming events, sessions),
not execution schemas.

## wire

The `wire` package defines protocol envelopes for requests and responses.

Contract highlights:

- All envelopes carry a **version** and **type**.
- Payloads are JSON-serializable and stable across transports.
- Unknown fields must be ignored for forward compatibility.

## transport

`transport` defines interfaces for moving `wire` payloads across channels.

Contract highlights:

- Implementations must be concurrency-safe.
- Close operations are idempotent.
- Errors must wrap `context.Canceled` and `context.DeadlineExceeded` when applicable.

## stream

Streams represent ordered event sequences.

Contract highlights:

- Events are delivered in-order.
- Close is idempotent; no panic on double-close.
- Readers must not block indefinitely after Close.

## session

Sessions track ongoing interactions across transports.

Contract highlights:

- Session IDs are unique and immutable.
- Context cancellation must terminate session-bound operations.

## content

`content` defines structured content blocks.

Contract highlights:

- Content blocks are self-describing (`type` + `data`).
- Blocks must be valid JSON and stable for storage/transmission.

## task

`task` defines long-running work events.

Contract highlights:

- Task state transitions are monotonic.
- Errors are represented using a stable error envelope.

## discover

Discovery request/response shapes.

Contract highlights:

- Search requests are deterministic.
- Responses include stable identifiers; schemas are not required in discovery.

## prompt + elicit

Prompt templates and elicitation flows.

Contract highlights:

- Prompts are explicit about required variables.
- Elicitation responses map 1:1 to requested fields.
