package prompt

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()

	p := Prompt{Name: "greeting"}
	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return nil, nil
	}

	err := r.Register(p, handler)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
}

func TestRegistry_Register_Duplicate(t *testing.T) {
	r := NewRegistry()

	p := Prompt{Name: "greeting"}
	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return nil, nil
	}

	_ = r.Register(p, handler)
	err := r.Register(p, handler)

	if !errors.Is(err, ErrDuplicatePrompt) {
		t.Errorf("Register() error = %v, want ErrDuplicatePrompt", err)
	}
}

func TestRegistry_Register_InvalidPrompt(t *testing.T) {
	r := NewRegistry()

	p := Prompt{Name: ""} // Invalid - empty name
	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return nil, nil
	}

	err := r.Register(p, handler)
	if !errors.Is(err, ErrInvalidPrompt) {
		t.Errorf("Register() error = %v, want ErrInvalidPrompt", err)
	}
}

func TestRegistry_List(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return nil, nil
	}

	_ = r.Register(Prompt{Name: "prompt1"}, handler)
	_ = r.Register(Prompt{Name: "prompt2"}, handler)

	prompts, err := r.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(prompts) != 2 {
		t.Errorf("List() length = %d, want 2", len(prompts))
	}
}

func TestRegistry_List_Empty(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	prompts, err := r.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(prompts) != 0 {
		t.Errorf("List() length = %d, want 0", len(prompts))
	}
}

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p := Prompt{
		Name: "greeting",
		Arguments: []Argument{
			{Name: "name", Required: true},
		},
	}
	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return []Message{
			NewUserMessage(TextContent("Hello, " + args["name"] + "!")),
		}, nil
	}

	_ = r.Register(p, handler)

	msgs, err := r.Get(ctx, "greeting", map[string]string{"name": "Alice"})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("Get() messages length = %d, want 1", len(msgs))
	}
	if msgs[0].Content[0].Text != "Hello, Alice!" {
		t.Errorf("message text = %q, want %q", msgs[0].Content[0].Text, "Hello, Alice!")
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	_, err := r.Get(ctx, "nonexistent", nil)
	if !errors.Is(err, ErrPromptNotFound) {
		t.Errorf("Get() error = %v, want ErrPromptNotFound", err)
	}
}

func TestRegistry_Get_WithArgs(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p := Prompt{
		Name: "format",
		Arguments: []Argument{
			{Name: "first", Required: true},
			{Name: "last", Required: true},
			{Name: "title", Required: false},
		},
	}
	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		name := args["first"] + " " + args["last"]
		if title, ok := args["title"]; ok {
			name = title + " " + name
		}
		return []Message{NewUserMessage(TextContent(name))}, nil
	}

	_ = r.Register(p, handler)

	// With all args
	msgs, err := r.Get(ctx, "format", map[string]string{
		"first": "John",
		"last":  "Doe",
		"title": "Dr.",
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if msgs[0].Content[0].Text != "Dr. John Doe" {
		t.Errorf("message text = %q, want %q", msgs[0].Content[0].Text, "Dr. John Doe")
	}
}

func TestRegistry_Get_MissingRequiredArg(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p := Prompt{
		Name: "greeting",
		Arguments: []Argument{
			{Name: "name", Required: true},
		},
	}
	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return nil, nil
	}

	_ = r.Register(p, handler)

	_, err := r.Get(ctx, "greeting", map[string]string{}) // Missing required arg
	if !errors.Is(err, ErrMissingArgument) {
		t.Errorf("Get() error = %v, want ErrMissingArgument", err)
	}
}

func TestRegistry_Get_HandlerError(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	testErr := errors.New("handler failed")
	p := Prompt{Name: "failing"}
	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return nil, testErr
	}

	_ = r.Register(p, handler)

	_, err := r.Get(ctx, "failing", nil)
	if !errors.Is(err, testErr) {
		t.Errorf("Get() error = %v, want testErr", err)
	}
}

func TestRegistry_ConcurrentSafety(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	handler := func(ctx context.Context, args map[string]string) ([]Message, error) {
		return []Message{NewUserMessage(TextContent("ok"))}, nil
	}

	var wg sync.WaitGroup

	// Concurrent registrations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p := Prompt{Name: "prompt-" + string(rune('a'+i))}
			_ = r.Register(p, handler)
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, _ = r.List(ctx)
		}()
		go func(i int) {
			defer wg.Done()
			name := "prompt-" + string(rune('a'+(i%10)))
			_, _ = r.Get(ctx, name, nil)
		}(i)
	}
	wg.Wait()
}
