package wire

import (
	"context"
	"testing"
)

func TestRegistry_Register(t *testing.T) {
	reg := NewRegistry()
	w := NewMCP()

	reg.Register("custom", w)

	got := reg.Get("custom")
	if got == nil {
		t.Error("Get(custom) returned nil after Register")
	}
	if got.Name() != "mcp" {
		t.Errorf("Get(custom).Name() = %q, want %q", got.Name(), "mcp")
	}
}

func TestRegistry_Get(t *testing.T) {
	reg := NewRegistry()
	w := NewA2A()
	reg.Register("test", w)

	got := reg.Get("test")
	if got == nil {
		t.Error("Get(test) returned nil")
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	reg := NewRegistry()

	got := reg.Get("notfound")
	if got != nil {
		t.Errorf("Get(notfound) = %v, want nil", got)
	}
}

func TestRegistry_List(t *testing.T) {
	reg := NewRegistry()
	reg.Register("alpha", NewMCP())
	reg.Register("beta", NewA2A())

	list := reg.List()
	if len(list) != 2 {
		t.Errorf("len(List()) = %d, want 2", len(list))
	}

	found := make(map[string]bool)
	for _, name := range list {
		found[name] = true
	}
	if !found["alpha"] {
		t.Error("List() missing 'alpha'")
	}
	if !found["beta"] {
		t.Error("List() missing 'beta'")
	}
}

func TestDefaultRegistry(t *testing.T) {
	reg := DefaultRegistry()

	list := reg.List()
	found := make(map[string]bool)
	for _, name := range list {
		found[name] = true
	}

	if !found["mcp"] {
		t.Error("DefaultRegistry missing 'mcp'")
	}
	if !found["a2a"] {
		t.Error("DefaultRegistry missing 'a2a'")
	}
	if !found["acp"] {
		t.Error("DefaultRegistry missing 'acp'")
	}
}

func TestRegistry_Convert_MCPtoA2A(t *testing.T) {
	reg := DefaultRegistry()
	ctx := context.Background()

	// Create an MCP request
	mcpWire := reg.Get("mcp")
	req := &Request{
		ID:     "convert-1",
		Method: "tools/call",
		ToolID: "search",
		Arguments: map[string]any{
			"query": "test",
		},
	}

	// Encode with MCP
	mcpData, err := mcpWire.EncodeRequest(ctx, req)
	if err != nil {
		t.Fatalf("MCP EncodeRequest error = %v", err)
	}

	// Decode with MCP
	decodedReq, err := mcpWire.DecodeRequest(ctx, mcpData)
	if err != nil {
		t.Fatalf("MCP DecodeRequest error = %v", err)
	}

	// Re-encode with A2A
	a2aWire := reg.Get("a2a")
	a2aData, err := a2aWire.EncodeRequest(ctx, decodedReq)
	if err != nil {
		t.Fatalf("A2A EncodeRequest error = %v", err)
	}

	// Verify A2A data is valid
	if len(a2aData) == 0 {
		t.Error("A2A encoded data is empty")
	}
}

func TestRegistry_Convert_A2AtoMCP(t *testing.T) {
	reg := DefaultRegistry()
	ctx := context.Background()

	// Create an A2A request
	a2aWire := reg.Get("a2a")
	req := &Request{
		ID:     "convert-2",
		Method: "tasks/send",
		ToolID: "analyze",
		Arguments: map[string]any{
			"data": "input",
		},
	}

	// Encode with A2A
	a2aData, err := a2aWire.EncodeRequest(ctx, req)
	if err != nil {
		t.Fatalf("A2A EncodeRequest error = %v", err)
	}

	// Decode with A2A
	decodedReq, err := a2aWire.DecodeRequest(ctx, a2aData)
	if err != nil {
		t.Fatalf("A2A DecodeRequest error = %v", err)
	}

	// Re-encode with MCP
	mcpWire := reg.Get("mcp")
	mcpData, err := mcpWire.EncodeRequest(ctx, decodedReq)
	if err != nil {
		t.Fatalf("MCP EncodeRequest error = %v", err)
	}

	// Verify MCP data is valid
	if len(mcpData) == 0 {
		t.Error("MCP encoded data is empty")
	}
}

func TestRegistry_Convert_MCPtoACP(t *testing.T) {
	reg := DefaultRegistry()
	ctx := context.Background()

	// Create an MCP request
	mcpWire := reg.Get("mcp")
	req := &Request{
		ID:     "convert-3",
		Method: "tools/call",
		ToolID: "compute",
		Arguments: map[string]any{
			"x": 10,
		},
	}

	// Encode with MCP
	mcpData, err := mcpWire.EncodeRequest(ctx, req)
	if err != nil {
		t.Fatalf("MCP EncodeRequest error = %v", err)
	}

	// Decode with MCP
	decodedReq, err := mcpWire.DecodeRequest(ctx, mcpData)
	if err != nil {
		t.Fatalf("MCP DecodeRequest error = %v", err)
	}

	// Re-encode with ACP
	acpWire := reg.Get("acp")
	acpData, err := acpWire.EncodeRequest(ctx, decodedReq)
	if err != nil {
		t.Fatalf("ACP EncodeRequest error = %v", err)
	}

	// Verify ACP data is valid
	if len(acpData) == 0 {
		t.Error("ACP encoded data is empty")
	}
}

func TestErrUnsupportedFormat(t *testing.T) {
	if ErrUnsupportedFormat == nil {
		t.Fatal("ErrUnsupportedFormat is nil")
	}
	if ErrUnsupportedFormat.Error() != "wire: unsupported format" {
		t.Errorf("ErrUnsupportedFormat.Error() = %q, want %q",
			ErrUnsupportedFormat.Error(), "wire: unsupported format")
	}
}

func TestErrEncodeFailure(t *testing.T) {
	if ErrEncodeFailure == nil {
		t.Fatal("ErrEncodeFailure is nil")
	}
}

func TestErrDecodeFailure(t *testing.T) {
	if ErrDecodeFailure == nil {
		t.Fatal("ErrDecodeFailure is nil")
	}
}
