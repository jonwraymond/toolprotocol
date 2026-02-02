package resource

import (
	"context"
	"fmt"
	"testing"
)

// BenchmarkNewRegistry measures registry creation.
func BenchmarkNewRegistry(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewRegistry()
	}
}

// BenchmarkRegistry_Register measures provider registration.
func BenchmarkRegistry_Register(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		registry := NewRegistry()
		provider := NewStaticProvider()
		b.StartTimer()

		_ = registry.Register("file", provider)
	}
}

// BenchmarkRegistry_Read measures resource read performance.
func BenchmarkRegistry_Read(b *testing.B) {
	registry := NewRegistry()
	provider := NewStaticProvider()
	provider.Add(
		Resource{URI: "file:///doc.txt"},
		Contents{URI: "file:///doc.txt", Text: "content"},
	)
	_ = registry.Register("file", provider)
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = registry.Read(ctx, "file:///doc.txt")
	}
}

// BenchmarkRegistry_List measures list performance.
func BenchmarkRegistry_List(b *testing.B) {
	registry := NewRegistry()
	provider := NewStaticProvider()
	for i := range 100 {
		provider.Add(
			Resource{URI: fmt.Sprintf("file:///doc%d.txt", i)},
			Contents{URI: fmt.Sprintf("file:///doc%d.txt", i), Text: "content"},
		)
	}
	_ = registry.Register("file", provider)
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = registry.List(ctx)
	}
}

// BenchmarkNewStaticProvider measures static provider creation.
func BenchmarkNewStaticProvider(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewStaticProvider()
	}
}

// BenchmarkStaticProvider_Add measures resource addition.
func BenchmarkStaticProvider_Add(b *testing.B) {
	provider := NewStaticProvider()
	res := Resource{URI: "file:///doc.txt"}
	contents := Contents{URI: "file:///doc.txt", Text: "content"}

	b.ResetTimer()
	for b.Loop() {
		provider.Add(res, contents)
	}
}

// BenchmarkStaticProvider_Read measures read performance.
func BenchmarkStaticProvider_Read(b *testing.B) {
	provider := NewStaticProvider()
	provider.Add(
		Resource{URI: "file:///doc.txt"},
		Contents{URI: "file:///doc.txt", Text: "content"},
	)
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = provider.Read(ctx, "file:///doc.txt")
	}
}

// BenchmarkStaticProvider_List measures list performance.
func BenchmarkStaticProvider_List(b *testing.B) {
	provider := NewStaticProvider()
	for i := range 100 {
		provider.Add(
			Resource{URI: fmt.Sprintf("file:///doc%d.txt", i)},
			Contents{URI: fmt.Sprintf("file:///doc%d.txt", i), Text: "content"},
		)
	}
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = provider.List(ctx)
	}
}

// BenchmarkResource_Clone measures clone performance.
func BenchmarkResource_Clone(b *testing.B) {
	res := &Resource{
		URI:         "file:///doc.txt",
		Name:        "Document",
		Description: "A document",
		MIMEType:    "text/plain",
		Annotations: map[string]any{"key": "value"},
	}

	b.ResetTimer()
	for b.Loop() {
		_ = res.Clone()
	}
}

// BenchmarkTemplate_Expand measures template expansion.
func BenchmarkTemplate_Expand(b *testing.B) {
	tmpl := &Template{
		URITemplate: "file:///users/{userId}/docs/{docId}",
	}
	values := map[string]string{"userId": "123", "docId": "456"}

	b.ResetTimer()
	for b.Loop() {
		_ = tmpl.Expand(values)
	}
}

// BenchmarkNewSubscriptionManager measures manager creation.
func BenchmarkNewSubscriptionManager(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewSubscriptionManager()
	}
}

// BenchmarkSubscriptionManager_Subscribe measures subscription.
func BenchmarkSubscriptionManager_Subscribe(b *testing.B) {
	mgr := NewSubscriptionManager()
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = mgr.Subscribe(ctx, "file:///config.json")
	}
}

// BenchmarkSubscriptionManager_Notify measures notification.
func BenchmarkSubscriptionManager_Notify(b *testing.B) {
	mgr := NewSubscriptionManager()
	ctx := context.Background()
	_, _ = mgr.Subscribe(ctx, "file:///config.json")
	contents := &Contents{URI: "file:///config.json", Text: "updated"}

	b.ResetTimer()
	for b.Loop() {
		mgr.Notify("file:///config.json", contents)
	}
}

// BenchmarkRegistry_Concurrent measures concurrent access.
func BenchmarkRegistry_Concurrent(b *testing.B) {
	registry := NewRegistry()
	provider := NewStaticProvider()
	for i := range 100 {
		provider.Add(
			Resource{URI: fmt.Sprintf("file:///doc%d.txt", i)},
			Contents{URI: fmt.Sprintf("file:///doc%d.txt", i), Text: "content"},
		)
	}
	_ = registry.Register("file", provider)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 3 {
			case 0:
				_, _ = registry.Read(ctx, fmt.Sprintf("file:///doc%d.txt", i%100))
			case 1:
				_, _ = registry.List(ctx)
			case 2:
				_, _ = registry.Templates(ctx)
			}
			i++
		}
	})
}
