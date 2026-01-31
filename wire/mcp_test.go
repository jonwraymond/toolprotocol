package wire

import (
	"context"
	"encoding/json"
	"testing"
)

func TestMCPWire_Name(t *testing.T) {
	w := NewMCP()
	if got := w.Name(); got != "mcp" {
		t.Errorf("Name() = %q, want %q", got, "mcp")
	}
}

func TestMCPWire_Version(t *testing.T) {
	w := NewMCP()
	if got := w.Version(); got != "2025-11-25" {
		t.Errorf("Version() = %q, want %q", got, "2025-11-25")
	}
}

func TestMCPWire_EncodeRequest_Basic(t *testing.T) {
	w := NewMCP()
	ctx := context.Background()

	req := &Request{
		ID:     "1",
		Method: "tools/call",
		ToolID: "search",
		Arguments: map[string]any{
			"query": "test",
		},
	}

	data, err := w.EncodeRequest(ctx, req)
	if err != nil {
		t.Fatalf("EncodeRequest error = %v", err)
	}

	// Verify it's valid JSON
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	// Verify JSON-RPC structure
	if result["jsonrpc"] != "2.0" {
		t.Errorf("jsonrpc = %v, want %q", result["jsonrpc"], "2.0")
	}
	if result["id"] != "1" {
		t.Errorf("id = %v, want %q", result["id"], "1")
	}
	if result["method"] != "tools/call" {
		t.Errorf("method = %v, want %q", result["method"], "tools/call")
	}
}

func TestMCPWire_EncodeRequest_WithMeta(t *testing.T) {
	w := NewMCP()
	ctx := context.Background()

	req := &Request{
		ID:     "2",
		Method: "tools/call",
		ToolID: "run",
		Arguments: map[string]any{
			"code": "print('hello')",
		},
		Meta: map[string]any{
			"progressToken": "progress-123",
		},
	}

	data, err := w.EncodeRequest(ctx, req)
	if err != nil {
		t.Fatalf("EncodeRequest error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	params, ok := result["params"].(map[string]any)
	if !ok {
		t.Fatal("params is not a map")
	}

	meta, ok := params["_meta"].(map[string]any)
	if !ok {
		t.Fatal("_meta is not present or not a map")
	}

	if meta["progressToken"] != "progress-123" {
		t.Errorf("progressToken = %v, want %q", meta["progressToken"], "progress-123")
	}
}

func TestMCPWire_DecodeRequest_Valid(t *testing.T) {
	w := NewMCP()
	ctx := context.Background()

	jsonData := `{
		"jsonrpc": "2.0",
		"id": "req-1",
		"method": "tools/call",
		"params": {
			"name": "search",
			"arguments": {"query": "test"}
		}
	}`

	req, err := w.DecodeRequest(ctx, []byte(jsonData))
	if err != nil {
		t.Fatalf("DecodeRequest error = %v", err)
	}

	if req.ID != "req-1" {
		t.Errorf("ID = %q, want %q", req.ID, "req-1")
	}
	if req.Method != "tools/call" {
		t.Errorf("Method = %q, want %q", req.Method, "tools/call")
	}
	if req.ToolID != "search" {
		t.Errorf("ToolID = %q, want %q", req.ToolID, "search")
	}
}

func TestMCPWire_DecodeRequest_Invalid(t *testing.T) {
	w := NewMCP()
	ctx := context.Background()

	_, err := w.DecodeRequest(ctx, []byte("not json"))
	if err == nil {
		t.Error("DecodeRequest(invalid) error = nil, want error")
	}
}

func TestMCPWire_EncodeResponse_Success(t *testing.T) {
	w := NewMCP()
	ctx := context.Background()

	resp := &Response{
		ID: "1",
		Content: []Content{
			{Type: ContentTypeText, Text: "Hello, world!"},
		},
	}

	data, err := w.EncodeResponse(ctx, resp)
	if err != nil {
		t.Fatalf("EncodeResponse error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	if result["jsonrpc"] != "2.0" {
		t.Errorf("jsonrpc = %v, want %q", result["jsonrpc"], "2.0")
	}
	if result["id"] != "1" {
		t.Errorf("id = %v, want %q", result["id"], "1")
	}
	if result["error"] != nil {
		t.Errorf("error = %v, want nil", result["error"])
	}
}

func TestMCPWire_EncodeResponse_Error(t *testing.T) {
	w := NewMCP()
	ctx := context.Background()

	resp := &Response{
		ID:      "1",
		IsError: true,
		Error: &Error{
			Code:    -32600,
			Message: "Invalid request",
		},
	}

	data, err := w.EncodeResponse(ctx, resp)
	if err != nil {
		t.Fatalf("EncodeResponse error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	errObj, ok := result["error"].(map[string]any)
	if !ok {
		t.Fatal("error is not present or not a map")
	}

	if errObj["code"].(float64) != -32600 {
		t.Errorf("error.code = %v, want %d", errObj["code"], -32600)
	}
}

func TestMCPWire_DecodeResponse_Success(t *testing.T) {
	w := NewMCP()
	ctx := context.Background()

	jsonData := `{
		"jsonrpc": "2.0",
		"id": "1",
		"result": {
			"content": [{"type": "text", "text": "result"}]
		}
	}`

	resp, err := w.DecodeResponse(ctx, []byte(jsonData))
	if err != nil {
		t.Fatalf("DecodeResponse error = %v", err)
	}

	if resp.ID != "1" {
		t.Errorf("ID = %q, want %q", resp.ID, "1")
	}
	if resp.IsError {
		t.Error("IsError = true, want false")
	}
}

func TestMCPWire_DecodeResponse_Error(t *testing.T) {
	w := NewMCP()
	ctx := context.Background()

	jsonData := `{
		"jsonrpc": "2.0",
		"id": "1",
		"error": {
			"code": -32601,
			"message": "Method not found"
		}
	}`

	resp, err := w.DecodeResponse(ctx, []byte(jsonData))
	if err != nil {
		t.Fatalf("DecodeResponse error = %v", err)
	}

	if !resp.IsError {
		t.Error("IsError = false, want true")
	}
	if resp.Error == nil {
		t.Fatal("Error is nil")
	}
	if resp.Error.Code != -32601 {
		t.Errorf("Error.Code = %d, want %d", resp.Error.Code, -32601)
	}
}

func TestMCPWire_EncodeToolList(t *testing.T) {
	w := NewMCP()
	ctx := context.Background()

	tools := []Tool{
		{
			Name:        "search",
			Description: "Search for items",
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
		t.Fatalf("EncodeToolList error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	toolsArr, ok := result["tools"].([]any)
	if !ok {
		t.Fatal("tools is not an array")
	}
	if len(toolsArr) != 1 {
		t.Errorf("len(tools) = %d, want 1", len(toolsArr))
	}
}

func TestMCPWire_DecodeToolList(t *testing.T) {
	w := NewMCP()
	ctx := context.Background()

	jsonData := `{
		"tools": [
			{
				"name": "search",
				"description": "Search for items",
				"inputSchema": {"type": "object"}
			}
		]
	}`

	tools, err := w.DecodeToolList(ctx, []byte(jsonData))
	if err != nil {
		t.Fatalf("DecodeToolList error = %v", err)
	}

	if len(tools) != 1 {
		t.Fatalf("len(tools) = %d, want 1", len(tools))
	}
	if tools[0].Name != "search" {
		t.Errorf("tools[0].Name = %q, want %q", tools[0].Name, "search")
	}
}

func TestMCPWire_Capabilities(t *testing.T) {
	w := NewMCP()
	caps := w.Capabilities()

	if caps == nil {
		t.Fatal("Capabilities() returned nil")
	}

	// MCP supports these features
	if !caps.Streaming {
		t.Error("Streaming = false, want true")
	}
	if !caps.Progress {
		t.Error("Progress = false, want true")
	}
	if !caps.Cancellation {
		t.Error("Cancellation = false, want true")
	}
}

func TestMCPWire_RoundTrip(t *testing.T) {
	w := NewMCP()
	ctx := context.Background()

	// Request round-trip
	origReq := &Request{
		ID:     "rt-1",
		Method: "tools/call",
		ToolID: "test",
		Arguments: map[string]any{
			"input": "value",
		},
	}

	encoded, err := w.EncodeRequest(ctx, origReq)
	if err != nil {
		t.Fatalf("EncodeRequest error = %v", err)
	}

	decoded, err := w.DecodeRequest(ctx, encoded)
	if err != nil {
		t.Fatalf("DecodeRequest error = %v", err)
	}

	if decoded.ID != origReq.ID {
		t.Errorf("round-trip ID = %q, want %q", decoded.ID, origReq.ID)
	}
	if decoded.Method != origReq.Method {
		t.Errorf("round-trip Method = %q, want %q", decoded.Method, origReq.Method)
	}
}

func TestMCPWire_ImplementsInterface(t *testing.T) {
	var _ Wire = (*MCPWire)(nil)
}
