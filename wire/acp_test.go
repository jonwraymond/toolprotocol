package wire

import (
	"context"
	"encoding/json"
	"testing"
)

func TestACPWire_Name(t *testing.T) {
	w := NewACP()
	if got := w.Name(); got != "acp" {
		t.Errorf("Name() = %q, want %q", got, "acp")
	}
}

func TestACPWire_Version(t *testing.T) {
	w := NewACP()
	if got := w.Version(); got == "" {
		t.Error("Version() is empty")
	}
}

func TestACPWire_EncodeRequest_Basic(t *testing.T) {
	w := NewACP()
	ctx := context.Background()

	req := &Request{
		ID:     "msg-1",
		Method: "agent/invoke",
		ToolID: "calculator",
		Arguments: map[string]any{
			"operation": "add",
			"a":         1,
			"b":         2,
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

	// ACP uses JSON-RPC 2.0
	if result["jsonrpc"] != "2.0" {
		t.Errorf("jsonrpc = %v, want %q", result["jsonrpc"], "2.0")
	}
}

func TestACPWire_DecodeRequest_Valid(t *testing.T) {
	w := NewACP()
	ctx := context.Background()

	jsonData := `{
		"jsonrpc": "2.0",
		"id": "msg-1",
		"method": "agent/invoke",
		"params": {
			"agentId": "calculator",
			"input": {"operation": "multiply", "x": 3, "y": 4}
		}
	}`

	req, err := w.DecodeRequest(ctx, []byte(jsonData))
	if err != nil {
		t.Fatalf("DecodeRequest error = %v", err)
	}

	if req.ID != "msg-1" {
		t.Errorf("ID = %q, want %q", req.ID, "msg-1")
	}
	if req.Method != "agent/invoke" {
		t.Errorf("Method = %q, want %q", req.Method, "agent/invoke")
	}
}

func TestACPWire_EncodeResponse_Success(t *testing.T) {
	w := NewACP()
	ctx := context.Background()

	resp := &Response{
		ID: "msg-1",
		Content: []Content{
			{Type: ContentTypeText, Text: "Result: 3"},
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
}

func TestACPWire_DecodeResponse(t *testing.T) {
	w := NewACP()
	ctx := context.Background()

	jsonData := `{
		"jsonrpc": "2.0",
		"id": "msg-1",
		"result": {
			"status": "success",
			"output": {"text": "Computation complete"}
		}
	}`

	resp, err := w.DecodeResponse(ctx, []byte(jsonData))
	if err != nil {
		t.Fatalf("DecodeResponse error = %v", err)
	}

	if resp.ID != "msg-1" {
		t.Errorf("ID = %q, want %q", resp.ID, "msg-1")
	}
	if resp.IsError {
		t.Error("IsError = true, want false")
	}
}

func TestACPWire_EncodeToolList(t *testing.T) {
	w := NewACP()
	ctx := context.Background()

	tools := []Tool{
		{
			Name:        "calculator",
			Description: "Performs mathematical operations",
			InputSchema: map[string]any{
				"type": "object",
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

	// ACP uses "agents" terminology
	if result["agents"] == nil && result["tools"] == nil {
		t.Error("expected agents or tools in result")
	}
}

func TestACPWire_DecodeToolList(t *testing.T) {
	w := NewACP()
	ctx := context.Background()

	jsonData := `{
		"agents": [
			{
				"id": "calculator",
				"name": "Calculator Agent",
				"description": "Performs math"
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
}

func TestACPWire_Capabilities(t *testing.T) {
	w := NewACP()
	caps := w.Capabilities()

	if caps == nil {
		t.Fatal("Capabilities() returned nil")
	}
}

func TestACPWire_RoundTrip(t *testing.T) {
	w := NewACP()
	ctx := context.Background()

	origReq := &Request{
		ID:     "acp-rt",
		Method: "agent/invoke",
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
}

func TestACPWire_ImplementsInterface(t *testing.T) {
	var _ Wire = (*ACPWire)(nil)
}
