package transport

import (
	"fmt"
	"testing"
)

// BenchmarkStdio_Name measures Name() call performance.
func BenchmarkStdio_Name(b *testing.B) {
	t := &StdioTransport{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.Name()
	}
}

// BenchmarkStdio_Info measures Info() call performance.
func BenchmarkStdio_Info(b *testing.B) {
	t := &StdioTransport{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.Info()
	}
}

// BenchmarkStdio_Close measures Close() call performance.
func BenchmarkStdio_Close(b *testing.B) {
	t := &StdioTransport{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.Close()
	}
}

// BenchmarkStreamable_Name measures Name() call performance.
func BenchmarkStreamable_Name(b *testing.B) {
	t := &StreamableHTTPTransport{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.Name()
	}
}

// BenchmarkStreamable_Info measures Info() call performance.
func BenchmarkStreamable_Info(b *testing.B) {
	t := &StreamableHTTPTransport{
		Config: StreamableConfig{
			HTTPConfig: HTTPConfig{
				Host: "localhost",
				Port: 8080,
				Path: "/mcp",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.Info()
	}
}

// BenchmarkStreamable_Info_DefaultPath measures Info() with default path.
func BenchmarkStreamable_Info_DefaultPath(b *testing.B) {
	t := &StreamableHTTPTransport{
		Config: StreamableConfig{
			HTTPConfig: HTTPConfig{
				Port: 8080,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.Info()
	}
}

// BenchmarkSSE_Name measures Name() call performance.
func BenchmarkSSE_Name(b *testing.B) {
	t := &SSETransport{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.Name()
	}
}

// BenchmarkSSE_Info measures Info() call performance.
func BenchmarkSSE_Info(b *testing.B) {
	t := &SSETransport{
		Config: SSEConfig{
			HTTPConfig: HTTPConfig{
				Host: "localhost",
				Port: 8081,
				Path: "/events",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.Info()
	}
}

// BenchmarkRegistry_Get measures registry lookup performance.
func BenchmarkRegistry_Get(b *testing.B) {
	reg := DefaultRegistry()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = reg.Get("stdio")
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

// BenchmarkRegistry_Register measures registration performance.
func BenchmarkRegistry_Register(b *testing.B) {
	reg := NewRegistry()
	factory := func(cfg any) (Transport, error) {
		return &StdioTransport{}, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reg.Register(fmt.Sprintf("transport-%d", i), factory)
	}
}

// BenchmarkRegistry_New measures transport creation performance.
func BenchmarkRegistry_New(b *testing.B) {
	reg := DefaultRegistry()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = reg.New("stdio", nil)
	}
}

// BenchmarkNew_Stdio measures factory creation for stdio.
func BenchmarkNew_Stdio(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = New("stdio", nil)
	}
}

// BenchmarkNew_Streamable measures factory creation for streamable.
func BenchmarkNew_Streamable(b *testing.B) {
	cfg := &StreamableConfig{
		HTTPConfig: HTTPConfig{
			Port: 8080,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = New("streamable", cfg)
	}
}

// BenchmarkNew_SSE measures factory creation for SSE.
func BenchmarkNew_SSE(b *testing.B) {
	cfg := &SSEConfig{
		HTTPConfig: HTTPConfig{
			Port: 8081,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = New("sse", cfg)
	}
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
				_ = reg.Get("stdio")
			case 1:
				_ = reg.Get("sse")
			case 2:
				_ = reg.Get("streamable")
			}
			i++
		}
	})
}

// BenchmarkStreamable_Info_Concurrent measures concurrent Info() calls.
func BenchmarkStreamable_Info_Concurrent(b *testing.B) {
	t := &StreamableHTTPTransport{
		Config: StreamableConfig{
			HTTPConfig: HTTPConfig{
				Host: "localhost",
				Port: 8080,
				Path: "/mcp",
			},
		},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = t.Info()
		}
	})
}
