package resource

import (
	"context"
	"sync"
	"testing"
)

func TestStaticProvider_List(t *testing.T) {
	p := NewStaticProvider()
	ctx := context.Background()

	p.Add(
		Resource{URI: "test://a", Name: "A"},
		Contents{URI: "test://a", Text: "content a"},
	)
	p.Add(
		Resource{URI: "test://b", Name: "B"},
		Contents{URI: "test://b", Text: "content b"},
	)

	resources, err := p.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(resources) != 2 {
		t.Errorf("List() length = %d, want 2", len(resources))
	}
}

func TestStaticProvider_List_Empty(t *testing.T) {
	p := NewStaticProvider()
	ctx := context.Background()

	resources, err := p.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(resources) != 0 {
		t.Errorf("List() length = %d, want 0", len(resources))
	}
}

func TestStaticProvider_Read(t *testing.T) {
	p := NewStaticProvider()
	ctx := context.Background()

	p.Add(
		Resource{URI: "test://doc", Name: "Doc"},
		Contents{URI: "test://doc", MIMEType: "text/plain", Text: "Hello"},
	)

	contents, err := p.Read(ctx, "test://doc")
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if contents.Text != "Hello" {
		t.Errorf("Text = %q, want %q", contents.Text, "Hello")
	}
	if contents.MIMEType != "text/plain" {
		t.Errorf("MIMEType = %q, want %q", contents.MIMEType, "text/plain")
	}
}

func TestStaticProvider_Read_NotFound(t *testing.T) {
	p := NewStaticProvider()
	ctx := context.Background()

	_, err := p.Read(ctx, "test://nonexistent")
	if err != ErrResourceNotFound {
		t.Errorf("Read() error = %v, want ErrResourceNotFound", err)
	}
}

func TestStaticProvider_Templates(t *testing.T) {
	p := NewStaticProvider()
	ctx := context.Background()

	p.AddTemplate(Template{
		URITemplate: "test://users/{id}",
		Name:        "User",
	})

	templates, err := p.Templates(ctx)
	if err != nil {
		t.Fatalf("Templates() error = %v", err)
	}
	if len(templates) != 1 {
		t.Errorf("Templates() length = %d, want 1", len(templates))
	}
	if templates[0].URITemplate != "test://users/{id}" {
		t.Errorf("URITemplate = %q", templates[0].URITemplate)
	}
}

func TestStaticProvider_Add(t *testing.T) {
	p := NewStaticProvider()
	ctx := context.Background()

	p.Add(
		Resource{URI: "test://new", Name: "New"},
		Contents{URI: "test://new", Text: "new content"},
	)

	resources, _ := p.List(ctx)
	if len(resources) != 1 {
		t.Errorf("after Add, List() length = %d, want 1", len(resources))
	}
}

func TestStaticProvider_Remove(t *testing.T) {
	p := NewStaticProvider()
	ctx := context.Background()

	p.Add(
		Resource{URI: "test://doc", Name: "Doc"},
		Contents{URI: "test://doc", Text: "content"},
	)

	p.Remove("test://doc")

	_, err := p.Read(ctx, "test://doc")
	if err != ErrResourceNotFound {
		t.Errorf("after Remove, Read() error = %v, want ErrResourceNotFound", err)
	}
}

func TestStaticProvider_ConcurrentSafety(t *testing.T) {
	p := NewStaticProvider()
	ctx := context.Background()

	var wg sync.WaitGroup

	// Concurrent adds
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			uri := "test://resource-" + string(rune('a'+i%26))
			p.Add(
				Resource{URI: uri, Name: "Test"},
				Contents{URI: uri, Text: "content"},
			)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = p.List(ctx)
		}()
	}

	wg.Wait()
}
