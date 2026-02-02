package wire

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

// BenchmarkMCP_EncodeRequest measures MCP request encoding performance.
func BenchmarkMCP_EncodeRequest(b *testing.B) {
	w := NewMCP()
	ctx := context.Background()
	req := &Request{
		ID:     "req-1",
		Method: "tools/call",
		ToolID: "search",
		Arguments: map[string]any{
			"query": "test query",
			"limit": 10,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.EncodeRequest(ctx, req)
	}
}

// BenchmarkMCP_DecodeRequest measures MCP request decoding performance.
func BenchmarkMCP_DecodeRequest(b *testing.B) {
	w := NewMCP()
	ctx := context.Background()
	data := []byte(`{"jsonrpc":"2.0","id":"req-1","method":"tools/call","params":{"name":"search","arguments":{"query":"test"}}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.DecodeRequest(ctx, data)
	}
}

// BenchmarkMCP_EncodeResponse measures MCP response encoding performance.
func BenchmarkMCP_EncodeResponse(b *testing.B) {
	w := NewMCP()
	ctx := context.Background()
	resp := &Response{
		ID: "req-1",
		Content: []Content{
			{Type: ContentTypeText, Text: "Search results for your query"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.EncodeResponse(ctx, resp)
	}
}

// BenchmarkMCP_DecodeResponse measures MCP response decoding performance.
func BenchmarkMCP_DecodeResponse(b *testing.B) {
	w := NewMCP()
	ctx := context.Background()
	data := []byte(`{"jsonrpc":"2.0","id":"req-1","result":{"content":[{"type":"text","text":"Hello world"}]}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.DecodeResponse(ctx, data)
	}
}

// BenchmarkMCP_EncodeToolList measures MCP tool list encoding performance.
func BenchmarkMCP_EncodeToolList(b *testing.B) {
	w := NewMCP()
	ctx := context.Background()
	tools := []Tool{
		{Name: "search", Description: "Search the web", InputSchema: map[string]any{"type": "object"}},
		{Name: "fetch", Description: "Fetch a URL", InputSchema: map[string]any{"type": "object"}},
		{Name: "execute", Description: "Execute code", InputSchema: map[string]any{"type": "object"}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.EncodeToolList(ctx, tools)
	}
}

// BenchmarkMCP_DecodeToolList measures MCP tool list decoding performance.
func BenchmarkMCP_DecodeToolList(b *testing.B) {
	w := NewMCP()
	ctx := context.Background()
	data := []byte(`{"tools":[{"name":"search","description":"Search"},{"name":"fetch","description":"Fetch"}]}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.DecodeToolList(ctx, data)
	}
}

// BenchmarkA2A_EncodeRequest measures A2A request encoding performance.
func BenchmarkA2A_EncodeRequest(b *testing.B) {
	w := NewA2A()
	ctx := context.Background()
	req := &Request{
		ID:     "req-1",
		Method: "message/send",
		ToolID: "skill-1",
		Arguments: map[string]any{
			"input": "test input",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.EncodeRequest(ctx, req)
	}
}

// BenchmarkA2A_DecodeRequest measures A2A request decoding performance.
func BenchmarkA2A_DecodeRequest(b *testing.B) {
	w := NewA2A()
	ctx := context.Background()
	data := []byte(`{"jsonrpc":"2.0","id":"req-1","method":"message/send","params":{"skillId":"skill-1","arguments":{"input":"test"}}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.DecodeRequest(ctx, data)
	}
}

// BenchmarkA2A_EncodeResponse measures A2A response encoding performance.
func BenchmarkA2A_EncodeResponse(b *testing.B) {
	w := NewA2A()
	ctx := context.Background()
	resp := &Response{
		ID: "req-1",
		Content: []Content{
			{Type: ContentTypeText, Text: "Response text"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.EncodeResponse(ctx, resp)
	}
}

// BenchmarkACP_EncodeRequest measures ACP request encoding performance.
func BenchmarkACP_EncodeRequest(b *testing.B) {
	w := NewACP()
	ctx := context.Background()
	req := &Request{
		ID:     "req-1",
		Method: "agent/invoke",
		ToolID: "agent-1",
		Arguments: map[string]any{
			"input": "test input",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.EncodeRequest(ctx, req)
	}
}

// BenchmarkACP_DecodeRequest measures ACP request decoding performance.
func BenchmarkACP_DecodeRequest(b *testing.B) {
	w := NewACP()
	ctx := context.Background()
	data := []byte(`{"jsonrpc":"2.0","id":"req-1","method":"agent/invoke","params":{"agentId":"agent-1","input":{"data":"test"}}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.DecodeRequest(ctx, data)
	}
}

// BenchmarkACP_EncodeResponse measures ACP response encoding performance.
func BenchmarkACP_EncodeResponse(b *testing.B) {
	w := NewACP()
	ctx := context.Background()
	resp := &Response{
		ID: "req-1",
		Content: []Content{
			{Type: ContentTypeText, Text: "Response text"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.EncodeResponse(ctx, resp)
	}
}

// BenchmarkRegistry_Get measures registry lookup performance.
func BenchmarkRegistry_Get(b *testing.B) {
	reg := DefaultRegistry()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = reg.Get("mcp")
	}
}

// BenchmarkRegistry_Get_Miss measures registry miss performance.
func BenchmarkRegistry_Get_Miss(b *testing.B) {
	reg := DefaultRegistry()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = reg.Get("nonexistent")
	}
}

// BenchmarkRegistry_List measures registry list performance.
func BenchmarkRegistry_List(b *testing.B) {
	reg := DefaultRegistry()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = reg.List()
	}
}

// BenchmarkRegistry_Register measures registry registration performance.
func BenchmarkRegistry_Register(b *testing.B) {
	reg := NewRegistry()
	w := NewMCP()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reg.Register(fmt.Sprintf("wire-%d", i), w)
	}
}

// BenchmarkMCP_Concurrent measures concurrent encode/decode operations.
func BenchmarkMCP_Concurrent(b *testing.B) {
	w := NewMCP()
	ctx := context.Background()
	req := &Request{
		ID:        "req-1",
		Method:    "tools/call",
		ToolID:    "search",
		Arguments: map[string]any{"query": "test"},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data, _ := w.EncodeRequest(ctx, req)
			_, _ = w.DecodeRequest(ctx, data)
		}
	})
}

// BenchmarkRegistry_Concurrent measures concurrent registry access.
func BenchmarkRegistry_Concurrent(b *testing.B) {
	reg := DefaultRegistry()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 3 {
			case 0:
				_ = reg.Get("mcp")
			case 1:
				_ = reg.Get("a2a")
			case 2:
				_ = reg.Get("acp")
			}
			i++
		}
	})
}

// BenchmarkRegistry_ConcurrentReadWrite measures mixed concurrent operations.
func BenchmarkRegistry_ConcurrentReadWrite(b *testing.B) {
	reg := NewRegistry()
	reg.Register("mcp", NewMCP())
	reg.Register("a2a", NewA2A())
	reg.Register("acp", NewACP())

	var mu sync.Mutex
	counter := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			i := counter
			counter++
			mu.Unlock()

			if i%10 == 0 {
				// 10% writes
				reg.Register(fmt.Sprintf("custom-%d", i), NewMCP())
			} else {
				// 90% reads
				_ = reg.Get("mcp")
			}
		}
	})
}

// BenchmarkMCP_LargeRequest measures encoding of large requests.
func BenchmarkMCP_LargeRequest(b *testing.B) {
	w := NewMCP()
	ctx := context.Background()

	// Create large arguments
	args := make(map[string]any)
	for i := 0; i < 100; i++ {
		args[fmt.Sprintf("field_%d", i)] = fmt.Sprintf("value_%d with some additional text", i)
	}

	req := &Request{
		ID:        "req-1",
		Method:    "tools/call",
		ToolID:    "complex_tool",
		Arguments: args,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.EncodeRequest(ctx, req)
	}
}

// BenchmarkMCP_LargeResponse measures encoding of large responses.
func BenchmarkMCP_LargeResponse(b *testing.B) {
	w := NewMCP()
	ctx := context.Background()

	// Create response with multiple content items
	content := make([]Content, 50)
	for i := range content {
		content[i] = Content{
			Type: ContentTypeText,
			Text: fmt.Sprintf("Response item %d with additional content", i),
		}
	}

	resp := &Response{
		ID:      "req-1",
		Content: content,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.EncodeResponse(ctx, resp)
	}
}
