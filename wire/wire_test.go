package wire

import (
	"context"
	"testing"
)

func TestRequest_Fields(t *testing.T) {
	req := Request{
		ID:     "req-123",
		Method: "tools/call",
		ToolID: "search",
		Arguments: map[string]any{
			"query": "test",
		},
		Meta: map[string]any{
			"requestId": "abc",
		},
	}

	if req.ID != "req-123" {
		t.Errorf("ID = %q, want %q", req.ID, "req-123")
	}
	if req.Method != "tools/call" {
		t.Errorf("Method = %q, want %q", req.Method, "tools/call")
	}
	if req.ToolID != "search" {
		t.Errorf("ToolID = %q, want %q", req.ToolID, "search")
	}
	if req.Arguments["query"] != "test" {
		t.Errorf("Arguments[query] = %v, want %q", req.Arguments["query"], "test")
	}
}

func TestRequest_Empty(t *testing.T) {
	var req Request
	if req.ID != "" {
		t.Errorf("empty Request.ID = %q, want empty", req.ID)
	}
	if req.Arguments != nil {
		t.Errorf("empty Request.Arguments = %v, want nil", req.Arguments)
	}
}

func TestResponse_Fields(t *testing.T) {
	resp := Response{
		ID: "req-123",
		Content: []Content{
			{Type: ContentTypeText, Text: "result"},
		},
		IsError: false,
		Meta: map[string]any{
			"duration": 100,
		},
	}

	if resp.ID != "req-123" {
		t.Errorf("ID = %q, want %q", resp.ID, "req-123")
	}
	if len(resp.Content) != 1 {
		t.Fatalf("len(Content) = %d, want 1", len(resp.Content))
	}
	if resp.Content[0].Type != ContentTypeText {
		t.Errorf("Content[0].Type = %q, want %q", resp.Content[0].Type, ContentTypeText)
	}
	if resp.IsError {
		t.Error("IsError = true, want false")
	}
}

func TestResponse_Error(t *testing.T) {
	resp := Response{
		ID:      "req-123",
		IsError: true,
		Error: &Error{
			Code:    -32600,
			Message: "Invalid request",
		},
	}

	if !resp.IsError {
		t.Error("IsError = false, want true")
	}
	if resp.Error == nil {
		t.Fatal("Error is nil")
	}
	if resp.Error.Code != -32600 {
		t.Errorf("Error.Code = %d, want %d", resp.Error.Code, -32600)
	}
	if resp.Error.Message != "Invalid request" {
		t.Errorf("Error.Message = %q, want %q", resp.Error.Message, "Invalid request")
	}
}

func TestError_Format(t *testing.T) {
	err := &Error{
		Code:    -32601,
		Message: "Method not found",
		Data:    map[string]any{"method": "unknown"},
	}

	errStr := err.Error()
	if errStr == "" {
		t.Error("Error() returned empty string")
	}
}

func TestError_Format_NoData(t *testing.T) {
	err := &Error{
		Code:    -32600,
		Message: "Invalid request",
	}

	errStr := err.Error()
	if errStr == "" {
		t.Error("Error() returned empty string")
	}
	// Should not contain "data:" when Data is nil
	expected := "Invalid request (code: -32600)"
	if errStr != expected {
		t.Errorf("Error() = %q, want %q", errStr, expected)
	}
}

func TestCapabilities_Defaults(t *testing.T) {
	var caps Capabilities

	if caps.Streaming {
		t.Error("default Streaming = true, want false")
	}
	if caps.BatchRequests {
		t.Error("default BatchRequests = true, want false")
	}
	if caps.Progress {
		t.Error("default Progress = true, want false")
	}
	if caps.Cancellation {
		t.Error("default Cancellation = true, want false")
	}
}

func TestCapabilities_AllEnabled(t *testing.T) {
	caps := Capabilities{
		Streaming:     true,
		BatchRequests: true,
		Progress:      true,
		Cancellation:  true,
	}

	if !caps.Streaming {
		t.Error("Streaming = false, want true")
	}
	if !caps.BatchRequests {
		t.Error("BatchRequests = false, want true")
	}
	if !caps.Progress {
		t.Error("Progress = false, want true")
	}
	if !caps.Cancellation {
		t.Error("Cancellation = false, want true")
	}
}

func TestContent_Text(t *testing.T) {
	c := Content{
		Type: ContentTypeText,
		Text: "Hello, world!",
	}

	if c.Type != ContentTypeText {
		t.Errorf("Type = %q, want %q", c.Type, ContentTypeText)
	}
	if c.Text != "Hello, world!" {
		t.Errorf("Text = %q, want %q", c.Text, "Hello, world!")
	}
}

func TestContent_Image(t *testing.T) {
	c := Content{
		Type:     ContentTypeImage,
		MIMEType: "image/png",
		Data:     []byte{0x89, 0x50, 0x4e, 0x47},
	}

	if c.Type != ContentTypeImage {
		t.Errorf("Type = %q, want %q", c.Type, ContentTypeImage)
	}
	if c.MIMEType != "image/png" {
		t.Errorf("MIMEType = %q, want %q", c.MIMEType, "image/png")
	}
}

func TestContent_Resource(t *testing.T) {
	c := Content{
		Type: ContentTypeResource,
		URI:  "file:///path/to/resource",
	}

	if c.Type != ContentTypeResource {
		t.Errorf("Type = %q, want %q", c.Type, ContentTypeResource)
	}
	if c.URI != "file:///path/to/resource" {
		t.Errorf("URI = %q, want %q", c.URI, "file:///path/to/resource")
	}
}

func TestTool_Fields(t *testing.T) {
	tool := Tool{
		Name:        "search",
		Description: "Search for items",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{"type": "string"},
			},
		},
	}

	if tool.Name != "search" {
		t.Errorf("Name = %q, want %q", tool.Name, "search")
	}
	if tool.Description != "Search for items" {
		t.Errorf("Description = %q, want %q", tool.Description, "Search for items")
	}
	if tool.InputSchema == nil {
		t.Error("InputSchema is nil")
	}
}

func TestContentType_Constants(t *testing.T) {
	tests := []struct {
		ct   ContentType
		want string
	}{
		{ContentTypeText, "text"},
		{ContentTypeImage, "image"},
		{ContentTypeResource, "resource"},
	}

	for _, tt := range tests {
		if string(tt.ct) != tt.want {
			t.Errorf("ContentType = %q, want %q", tt.ct, tt.want)
		}
	}
}

func TestWireInterface_Contract(t *testing.T) {
	// Verify Wire interface has expected methods via compile-time check
	var _ Wire = (*mockWire)(nil)
}

// mockWire implements Wire for testing interface contract
type mockWire struct{}

func (m *mockWire) Name() string                                                     { return "mock" }
func (m *mockWire) Version() string                                                  { return "1.0" }
func (m *mockWire) EncodeRequest(ctx context.Context, req *Request) ([]byte, error)  { return nil, nil }
func (m *mockWire) DecodeRequest(ctx context.Context, data []byte) (*Request, error) { return nil, nil }
func (m *mockWire) EncodeResponse(ctx context.Context, resp *Response) ([]byte, error) {
	return nil, nil
}
func (m *mockWire) DecodeResponse(ctx context.Context, data []byte) (*Response, error) {
	return nil, nil
}
func (m *mockWire) EncodeToolList(ctx context.Context, tools []Tool) ([]byte, error) { return nil, nil }
func (m *mockWire) DecodeToolList(ctx context.Context, data []byte) ([]Tool, error)  { return nil, nil }
func (m *mockWire) Capabilities() *Capabilities                                      { return nil }
