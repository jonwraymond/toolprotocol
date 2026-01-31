package resource

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()
	p := NewStaticProvider()

	err := r.Register("test", p)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
}

func TestRegistry_Register_Duplicate(t *testing.T) {
	r := NewRegistry()
	p := NewStaticProvider()

	_ = r.Register("test", p)
	err := r.Register("test", p)

	if !errors.Is(err, ErrDuplicateProvider) {
		t.Errorf("Register() error = %v, want ErrDuplicateProvider", err)
	}
}

func TestRegistry_Unregister(t *testing.T) {
	r := NewRegistry()
	p := NewStaticProvider()

	_ = r.Register("test", p)
	err := r.Unregister("test")

	if err != nil {
		t.Fatalf("Unregister() error = %v", err)
	}
}

func TestRegistry_Unregister_NotFound(t *testing.T) {
	r := NewRegistry()

	err := r.Unregister("nonexistent")
	if !errors.Is(err, ErrProviderNotFound) {
		t.Errorf("Unregister() error = %v, want ErrProviderNotFound", err)
	}
}

func TestRegistry_List(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p1 := NewStaticProvider()
	p1.Add(Resource{URI: "test1://a"}, Contents{})

	p2 := NewStaticProvider()
	p2.Add(Resource{URI: "test2://b"}, Contents{})

	_ = r.Register("test1", p1)
	_ = r.Register("test2", p2)

	resources, err := r.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(resources) != 2 {
		t.Errorf("List() length = %d, want 2", len(resources))
	}
}

func TestRegistry_List_MergesProviders(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p1 := NewStaticProvider()
	p1.Add(Resource{URI: "a://1"}, Contents{})
	p1.Add(Resource{URI: "a://2"}, Contents{})

	p2 := NewStaticProvider()
	p2.Add(Resource{URI: "b://1"}, Contents{})

	_ = r.Register("a", p1)
	_ = r.Register("b", p2)

	resources, err := r.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(resources) != 3 {
		t.Errorf("List() should merge all providers, got %d resources", len(resources))
	}
}

func TestRegistry_Read(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p := NewStaticProvider()
	p.Add(
		Resource{URI: "test://doc"},
		Contents{URI: "test://doc", Text: "content"},
	)
	_ = r.Register("test", p)

	contents, err := r.Read(ctx, "test://doc")
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if contents.Text != "content" {
		t.Errorf("Text = %q, want %q", contents.Text, "content")
	}
}

func TestRegistry_Read_RoutesToProvider(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p1 := NewStaticProvider()
	p1.Add(Resource{URI: "scheme1://doc"}, Contents{Text: "from p1"})

	p2 := NewStaticProvider()
	p2.Add(Resource{URI: "scheme2://doc"}, Contents{Text: "from p2"})

	_ = r.Register("scheme1", p1)
	_ = r.Register("scheme2", p2)

	c1, _ := r.Read(ctx, "scheme1://doc")
	c2, _ := r.Read(ctx, "scheme2://doc")

	if c1.Text != "from p1" {
		t.Errorf("scheme1 should route to p1, got %q", c1.Text)
	}
	if c2.Text != "from p2" {
		t.Errorf("scheme2 should route to p2, got %q", c2.Text)
	}
}

func TestRegistry_Read_UnknownScheme(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	_, err := r.Read(ctx, "unknown://doc")
	if !errors.Is(err, ErrProviderNotFound) {
		t.Errorf("Read() error = %v, want ErrProviderNotFound", err)
	}
}

func TestRegistry_Read_InvalidURI(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	_, err := r.Read(ctx, "no-scheme")
	if !errors.Is(err, ErrInvalidURI) {
		t.Errorf("Read() error = %v, want ErrInvalidURI", err)
	}
}

func TestRegistry_Templates(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	p := NewStaticProvider()
	p.AddTemplate(Template{URITemplate: "test://{id}"})
	_ = r.Register("test", p)

	templates, err := r.Templates(ctx)
	if err != nil {
		t.Fatalf("Templates() error = %v", err)
	}
	if len(templates) != 1 {
		t.Errorf("Templates() length = %d, want 1", len(templates))
	}
}

func TestRegistry_ConcurrentSafety(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	var wg sync.WaitGroup

	// Concurrent registrations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p := NewStaticProvider()
			p.Add(Resource{URI: "s://r"}, Contents{Text: "c"})
			scheme := "scheme" + string(rune('a'+i))
			_ = r.Register(scheme, p)
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, _ = r.List(ctx)
		}()
		go func() {
			defer wg.Done()
			_, _ = r.Templates(ctx)
		}()
	}
	wg.Wait()
}
