package prompt

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

// BenchmarkRegistry_Register measures prompt registration.
func BenchmarkRegistry_Register(b *testing.B) {
	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return nil, nil
	}

	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		registry := NewRegistry()
		p := Prompt{Name: "greeting"}
		b.StartTimer()

		_ = registry.Register(p, handler)
	}
}

// BenchmarkRegistry_List measures prompt listing.
func BenchmarkRegistry_List(b *testing.B) {
	registry := NewRegistry()
	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return nil, nil
	}

	for i := range 100 {
		_ = registry.Register(Prompt{Name: fmt.Sprintf("prompt-%d", i)}, handler)
	}

	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = registry.List(ctx)
	}
}

// BenchmarkRegistry_Get measures prompt retrieval.
func BenchmarkRegistry_Get(b *testing.B) {
	registry := NewRegistry()
	p := Prompt{
		Name: "greeting",
		Arguments: []Argument{
			{Name: "name", Required: true},
		},
	}
	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return []Message{NewUserMessage(TextContent("Hello, " + args["name"]))}, nil
	}
	_ = registry.Register(p, handler)

	ctx := context.Background()
	args := map[string]string{"name": "Alice"}

	b.ResetTimer()
	for b.Loop() {
		_, _ = registry.Get(ctx, "greeting", args)
	}
}

// BenchmarkNewPromptBuilder measures builder creation.
func BenchmarkNewPromptBuilder(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewPromptBuilder("greeting")
	}
}

// BenchmarkPromptBuilder_Build measures building a prompt.
func BenchmarkPromptBuilder_Build(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewPromptBuilder("greeting").
			WithDescription("A greeting").
			WithRequiredArgument("name", "Name").
			WithArgument("title", "Title", false).
			Build()
	}
}

// BenchmarkNewUserMessage measures message creation.
func BenchmarkNewUserMessage(b *testing.B) {
	content := TextContent("Hello!")

	b.ResetTimer()
	for b.Loop() {
		_ = NewUserMessage(content)
	}
}

// BenchmarkTextContent measures text content creation.
func BenchmarkTextContent(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = TextContent("Hello, world!")
	}
}

// BenchmarkImageContent measures image content creation.
func BenchmarkImageContent(b *testing.B) {
	data := make([]byte, 1024)

	b.ResetTimer()
	for b.Loop() {
		_ = ImageContent("image/png", data)
	}
}

// BenchmarkResourceContent measures resource content creation.
func BenchmarkResourceContent(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = ResourceContent("file:///doc.txt")
	}
}

// BenchmarkPrompt_RequiredArgs measures required args extraction.
func BenchmarkPrompt_RequiredArgs(b *testing.B) {
	p := Prompt{
		Name: "test",
		Arguments: []Argument{
			{Name: "a", Required: true},
			{Name: "b", Required: false},
			{Name: "c", Required: true},
			{Name: "d", Required: false},
			{Name: "e", Required: true},
		},
	}

	b.ResetTimer()
	for b.Loop() {
		_ = p.RequiredArgs()
	}
}

// BenchmarkPrompt_Validate measures validation.
func BenchmarkPrompt_Validate(b *testing.B) {
	p := Prompt{Name: "greeting"}

	b.ResetTimer()
	for b.Loop() {
		_ = p.Validate()
	}
}

// BenchmarkRole_Valid measures role validation.
func BenchmarkRole_Valid(b *testing.B) {
	roles := []Role{RoleUser, RoleAssistant}

	b.ResetTimer()
	for b.Loop() {
		for _, r := range roles {
			_ = r.Valid()
		}
	}
}

// BenchmarkContentType_Valid measures content type validation.
func BenchmarkContentType_Valid(b *testing.B) {
	types := []ContentType{ContentText, ContentImage, ContentResource}

	b.ResetTimer()
	for b.Loop() {
		for _, ct := range types {
			_ = ct.Valid()
		}
	}
}

// BenchmarkRegistry_Concurrent measures concurrent access.
func BenchmarkRegistry_Concurrent(b *testing.B) {
	registry := NewRegistry()
	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return []Message{NewUserMessage(TextContent("Hello"))}, nil
	}

	for i := range 100 {
		p := Prompt{
			Name:      fmt.Sprintf("prompt-%d", i),
			Arguments: []Argument{{Name: "name", Required: true}},
		}
		_ = registry.Register(p, handler)
	}

	ctx := context.Background()
	args := map[string]string{"name": "test"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 2 {
			case 0:
				_, _ = registry.List(ctx)
			case 1:
				_, _ = registry.Get(ctx, fmt.Sprintf("prompt-%d", i%100), args)
			}
			i++
		}
	})
}
