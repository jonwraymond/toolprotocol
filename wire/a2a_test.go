package wire

import (
	"context"
	"encoding/json"
	"testing"
)

func TestA2AWire_Name(t *testing.T) {
	w := NewA2A()
	if got := w.Name(); got != "a2a" {
		t.Errorf("Name() = %q, want %q", got, "a2a")
	}
}

func TestA2AWire_Version(t *testing.T) {
	w := NewA2A()
	// A2A version should be non-empty
	if got := w.Version(); got == "" {
		t.Error("Version() is empty")
	}
}

func TestA2AWire_EncodeRequest_Basic(t *testing.T) {
	w := NewA2A()
	ctx := context.Background()

	req := &Request{
		ID:     "task-1",
		Method: "tasks/send",
		ToolID: "search",
		Arguments: map[string]any{
			"query": "test",
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

	// A2A uses JSON-RPC 2.0
	if result["jsonrpc"] != "2.0" {
		t.Errorf("jsonrpc = %v, want %q", result["jsonrpc"], "2.0")
	}
}

func TestA2AWire_EncodeRequest_AgentCard(t *testing.T) {
	w := NewA2A()
	ctx := context.Background()

	req := &Request{
		ID:     "task-2",
		Method: "tasks/send",
		ToolID: "analyze",
		Arguments: map[string]any{
			"data": "input",
		},
		Meta: map[string]any{
			"agentCard": map[string]any{
				"name":    "TestAgent",
				"version": "1.0",
			},
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

	// Verify agent card is preserved in params
	if params, ok := result["params"].(map[string]any); ok {
		if meta, ok := params["_meta"].(map[string]any); ok {
			if meta["agentCard"] == nil {
				t.Error("agentCard not preserved in metadata")
			}
		}
	}
}

func TestA2AWire_DecodeRequest_Valid(t *testing.T) {
	w := NewA2A()
	ctx := context.Background()

	jsonData := `{
		"jsonrpc": "2.0",
		"id": "task-1",
		"method": "tasks/send",
		"params": {
			"id": "task-1",
			"message": {
				"role": "user",
				"parts": [{"kind": "text", "text": "search for test"}]
			}
		}
	}`

	req, err := w.DecodeRequest(ctx, []byte(jsonData))
	if err != nil {
		t.Fatalf("DecodeRequest error = %v", err)
	}

	if req.ID != "task-1" {
		t.Errorf("ID = %q, want %q", req.ID, "task-1")
	}
	if req.Method != "tasks/send" {
		t.Errorf("Method = %q, want %q", req.Method, "tasks/send")
	}
}

func TestA2AWire_EncodeResponse_Success(t *testing.T) {
	w := NewA2A()
	ctx := context.Background()

	resp := &Response{
		ID: "task-1",
		Content: []Content{
			{Type: ContentTypeText, Text: "Result found"},
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

func TestA2AWire_EncodeResponse_WithArtifacts(t *testing.T) {
	w := NewA2A()
	ctx := context.Background()

	resp := &Response{
		ID: "task-1",
		Content: []Content{
			{Type: ContentTypeText, Text: "Generated code"},
			{Type: ContentTypeResource, URI: "file:///output.py", MIMEType: "text/x-python"},
		},
		Meta: map[string]any{
			"artifacts": []map[string]any{
				{"type": "file", "uri": "file:///output.py"},
			},
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

	// Should have result with content
	if result["result"] == nil {
		t.Error("result is nil")
	}
}

func TestA2AWire_DecodeResponse(t *testing.T) {
	w := NewA2A()
	ctx := context.Background()

	jsonData := `{
		"jsonrpc": "2.0",
		"id": "task-1",
		"result": {
			"status": {
				"state": "completed"
			},
			"artifacts": [
				{"parts": [{"kind": "text", "text": "output"}]}
			]
		}
	}`

	resp, err := w.DecodeResponse(ctx, []byte(jsonData))
	if err != nil {
		t.Fatalf("DecodeResponse error = %v", err)
	}

	if resp.ID != "task-1" {
		t.Errorf("ID = %q, want %q", resp.ID, "task-1")
	}
	if resp.IsError {
		t.Error("IsError = true, want false")
	}
}

func TestA2AWire_EncodeToolList_AsActions(t *testing.T) {
	w := NewA2A()
	ctx := context.Background()

	tools := []Tool{
		{
			Name:        "search",
			Description: "Search capability",
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

	// A2A uses "skills" or "capabilities" terminology
	if result["skills"] == nil && result["capabilities"] == nil && result["tools"] == nil {
		t.Error("expected skills, capabilities, or tools in result")
	}
}

func TestA2AWire_DecodeToolList(t *testing.T) {
	w := NewA2A()
	ctx := context.Background()

	jsonData := `{
		"skills": [
			{
				"id": "search",
				"name": "Search",
				"description": "Search capability"
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
	if tools[0].Name != "search" && tools[0].Name != "Search" {
		t.Errorf("tools[0].Name = %q, want search or Search", tools[0].Name)
	}
}

func TestA2AWire_Capabilities(t *testing.T) {
	w := NewA2A()
	caps := w.Capabilities()

	if caps == nil {
		t.Fatal("Capabilities() returned nil")
	}

	// A2A supports streaming
	if !caps.Streaming {
		t.Error("Streaming = false, want true")
	}
}

func TestA2AWire_RoundTrip(t *testing.T) {
	w := NewA2A()
	ctx := context.Background()

	origReq := &Request{
		ID:     "a2a-rt",
		Method: "tasks/send",
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

func TestA2AWire_ImplementsInterface(t *testing.T) {
	var _ Wire = (*A2AWire)(nil)
}
