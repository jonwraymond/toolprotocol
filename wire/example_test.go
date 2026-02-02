package wire_test

import (
	"context"
	"fmt"

	"github.com/jonwraymond/toolprotocol/wire"
)

func ExampleNewMCP() {
	w := wire.NewMCP()

	fmt.Println("Name:", w.Name())
	fmt.Println("Version:", w.Version())
	// Output:
	// Name: mcp
	// Version: 2025-11-25
}

func ExampleNewA2A() {
	w := wire.NewA2A()

	fmt.Println("Name:", w.Name())
	fmt.Println("Version:", w.Version())
	// Output:
	// Name: a2a
	// Version: 0.2.1
}

func ExampleNewACP() {
	w := wire.NewACP()

	fmt.Println("Name:", w.Name())
	fmt.Println("Version:", w.Version())
	// Output:
	// Name: acp
	// Version: 1.0.0
}

func ExampleMCPWire_EncodeRequest() {
	w := wire.NewMCP()
	ctx := context.Background()

	req := &wire.Request{
		ID:     "req-1",
		Method: "tools/call",
		ToolID: "search",
		Arguments: map[string]any{
			"query": "golang tutorials",
		},
	}

	data, err := w.EncodeRequest(ctx, req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Output is JSON-RPC 2.0 format
	fmt.Println("Encoded successfully:", len(data) > 0)
	// Output:
	// Encoded successfully: true
}

func ExampleMCPWire_DecodeRequest() {
	w := wire.NewMCP()
	ctx := context.Background()

	// JSON-RPC 2.0 request
	data := []byte(`{
		"jsonrpc": "2.0",
		"id": "req-1",
		"method": "tools/call",
		"params": {
			"name": "search",
			"arguments": {"query": "test"}
		}
	}`)

	req, err := w.DecodeRequest(ctx, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("ID:", req.ID)
	fmt.Println("Method:", req.Method)
	fmt.Println("ToolID:", req.ToolID)
	fmt.Println("Query:", req.Arguments["query"])
	// Output:
	// ID: req-1
	// Method: tools/call
	// ToolID: search
	// Query: test
}

func ExampleMCPWire_EncodeResponse() {
	w := wire.NewMCP()
	ctx := context.Background()

	resp := &wire.Response{
		ID: "req-1",
		Content: []wire.Content{
			{Type: wire.ContentTypeText, Text: "Hello, world!"},
		},
	}

	data, err := w.EncodeResponse(ctx, resp)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Encoded successfully:", len(data) > 0)
	// Output:
	// Encoded successfully: true
}

func ExampleMCPWire_DecodeResponse() {
	w := wire.NewMCP()
	ctx := context.Background()

	// JSON-RPC 2.0 response
	data := []byte(`{
		"jsonrpc": "2.0",
		"id": "req-1",
		"result": {
			"content": [
				{"type": "text", "text": "Search results found"}
			]
		}
	}`)

	resp, err := w.DecodeResponse(ctx, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("ID:", resp.ID)
	fmt.Println("IsError:", resp.IsError)
	fmt.Println("Content count:", len(resp.Content))
	fmt.Println("Text:", resp.Content[0].Text)
	// Output:
	// ID: req-1
	// IsError: false
	// Content count: 1
	// Text: Search results found
}

func ExampleMCPWire_DecodeResponse_error() {
	w := wire.NewMCP()
	ctx := context.Background()

	// JSON-RPC 2.0 error response
	data := []byte(`{
		"jsonrpc": "2.0",
		"id": "req-1",
		"error": {
			"code": -32600,
			"message": "Invalid Request"
		}
	}`)

	resp, err := w.DecodeResponse(ctx, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("ID:", resp.ID)
	fmt.Println("IsError:", resp.IsError)
	fmt.Println("Error code:", resp.Error.Code)
	fmt.Println("Error message:", resp.Error.Message)
	// Output:
	// ID: req-1
	// IsError: true
	// Error code: -32600
	// Error message: Invalid Request
}

func ExampleMCPWire_EncodeToolList() {
	w := wire.NewMCP()
	ctx := context.Background()

	tools := []wire.Tool{
		{
			Name:        "search",
			Description: "Search the web",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{"type": "string"},
				},
			},
		},
	}

	data, err := w.EncodeToolList(ctx, tools)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Encoded successfully:", len(data) > 0)
	// Output:
	// Encoded successfully: true
}

func ExampleMCPWire_DecodeToolList() {
	w := wire.NewMCP()
	ctx := context.Background()

	data := []byte(`{
		"tools": [
			{
				"name": "search",
				"description": "Search the web",
				"inputSchema": {"type": "object"}
			}
		]
	}`)

	tools, err := w.DecodeToolList(ctx, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Tool count:", len(tools))
	fmt.Println("Tool name:", tools[0].Name)
	fmt.Println("Tool description:", tools[0].Description)
	// Output:
	// Tool count: 1
	// Tool name: search
	// Tool description: Search the web
}

func ExampleMCPWire_Capabilities() {
	w := wire.NewMCP()
	caps := w.Capabilities()

	fmt.Println("Streaming:", caps.Streaming)
	fmt.Println("BatchRequests:", caps.BatchRequests)
	fmt.Println("Progress:", caps.Progress)
	fmt.Println("Cancellation:", caps.Cancellation)
	// Output:
	// Streaming: true
	// BatchRequests: false
	// Progress: true
	// Cancellation: true
}

func ExampleNewRegistry() {
	reg := wire.NewRegistry()

	// Register a custom wire format
	reg.Register("mcp", wire.NewMCP())

	// List registered formats
	names := reg.List()
	fmt.Println("Registered:", len(names) > 0)
	// Output:
	// Registered: true
}

func ExampleRegistry_Register() {
	reg := wire.NewRegistry()

	// Register formats
	reg.Register("mcp", wire.NewMCP())
	reg.Register("a2a", wire.NewA2A())
	reg.Register("acp", wire.NewACP())

	fmt.Println("Registered count:", len(reg.List()))
	// Output:
	// Registered count: 3
}

func ExampleRegistry_Get() {
	reg := wire.NewRegistry()
	reg.Register("mcp", wire.NewMCP())

	// Get existing format
	w := reg.Get("mcp")
	fmt.Println("Found mcp:", w != nil)
	fmt.Println("Name:", w.Name())

	// Get non-existent format
	missing := reg.Get("unknown")
	fmt.Println("Found unknown:", missing != nil)
	// Output:
	// Found mcp: true
	// Name: mcp
	// Found unknown: false
}

func ExampleRegistry_List() {
	reg := wire.NewRegistry()
	reg.Register("mcp", wire.NewMCP())
	reg.Register("a2a", wire.NewA2A())

	names := reg.List()
	fmt.Println("Count:", len(names))
	// Output:
	// Count: 2
}

func ExampleDefaultRegistry() {
	reg := wire.DefaultRegistry()

	// Default registry has mcp, a2a, and acp pre-registered
	mcp := reg.Get("mcp")
	a2a := reg.Get("a2a")
	acp := reg.Get("acp")

	fmt.Println("Has mcp:", mcp != nil)
	fmt.Println("Has a2a:", a2a != nil)
	fmt.Println("Has acp:", acp != nil)
	// Output:
	// Has mcp: true
	// Has a2a: true
	// Has acp: true
}

func ExampleA2AWire_Capabilities() {
	w := wire.NewA2A()
	caps := w.Capabilities()

	fmt.Println("Streaming:", caps.Streaming)
	fmt.Println("BatchRequests:", caps.BatchRequests)
	fmt.Println("Progress:", caps.Progress)
	fmt.Println("Cancellation:", caps.Cancellation)
	// Output:
	// Streaming: true
	// BatchRequests: false
	// Progress: true
	// Cancellation: true
}

func ExampleACPWire_Capabilities() {
	w := wire.NewACP()
	caps := w.Capabilities()

	fmt.Println("Streaming:", caps.Streaming)
	fmt.Println("BatchRequests:", caps.BatchRequests)
	fmt.Println("Progress:", caps.Progress)
	fmt.Println("Cancellation:", caps.Cancellation)
	// Output:
	// Streaming: false
	// BatchRequests: true
	// Progress: false
	// Cancellation: true
}

func ExampleError_Error() {
	err := &wire.Error{
		Code:    -32600,
		Message: "Invalid Request",
	}

	fmt.Println(err.Error())
	// Output:
	// Invalid Request (code: -32600)
}

func ExampleError_Error_withData() {
	err := &wire.Error{
		Code:    -32602,
		Message: "Invalid params",
		Data:    "missing required field: query",
	}

	fmt.Println(err.Error())
	// Output:
	// Invalid params (code: -32602, data: missing required field: query)
}
