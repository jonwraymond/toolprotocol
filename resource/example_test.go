package resource_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonwraymond/toolprotocol/resource"
)

func ExampleNewRegistry() {
	registry := resource.NewRegistry()

	fmt.Printf("Type: %T\n", registry)
	// Output:
	// Type: *resource.Registry
}

func ExampleRegistry_Register() {
	registry := resource.NewRegistry()
	provider := resource.NewStaticProvider()

	err := registry.Register("file", provider)
	fmt.Println("Register error:", err)
	// Output:
	// Register error: <nil>
}

func ExampleRegistry_Register_duplicate() {
	registry := resource.NewRegistry()
	provider := resource.NewStaticProvider()

	_ = registry.Register("file", provider)
	err := registry.Register("file", provider)
	fmt.Println("Is duplicate:", errors.Is(err, resource.ErrDuplicateProvider))
	// Output:
	// Is duplicate: true
}

func ExampleRegistry_Unregister() {
	registry := resource.NewRegistry()
	provider := resource.NewStaticProvider()
	_ = registry.Register("file", provider)

	err := registry.Unregister("file")
	fmt.Println("Unregister error:", err)
	// Output:
	// Unregister error: <nil>
}

func ExampleRegistry_Read() {
	registry := resource.NewRegistry()
	provider := resource.NewStaticProvider()

	provider.Add(
		resource.Resource{URI: "file:///doc.txt", Name: "Document"},
		resource.Contents{URI: "file:///doc.txt", Text: "Hello, World!"},
	)
	_ = registry.Register("file", provider)

	ctx := context.Background()
	contents, err := registry.Read(ctx, "file:///doc.txt")
	fmt.Println("Error:", err)
	fmt.Println("Text:", contents.Text)
	// Output:
	// Error: <nil>
	// Text: Hello, World!
}

func ExampleRegistry_Read_invalidURI() {
	registry := resource.NewRegistry()
	ctx := context.Background()

	_, err := registry.Read(ctx, "invalid-uri")
	fmt.Println("Is invalid URI:", errors.Is(err, resource.ErrInvalidURI))
	// Output:
	// Is invalid URI: true
}

func ExampleRegistry_Read_providerNotFound() {
	registry := resource.NewRegistry()
	ctx := context.Background()

	_, err := registry.Read(ctx, "unknown://resource")
	fmt.Println("Is provider not found:", errors.Is(err, resource.ErrProviderNotFound))
	// Output:
	// Is provider not found: true
}

func ExampleRegistry_List() {
	registry := resource.NewRegistry()
	provider := resource.NewStaticProvider()

	provider.Add(
		resource.Resource{URI: "file:///a.txt", Name: "File A"},
		resource.Contents{URI: "file:///a.txt", Text: "content a"},
	)
	provider.Add(
		resource.Resource{URI: "file:///b.txt", Name: "File B"},
		resource.Contents{URI: "file:///b.txt", Text: "content b"},
	)
	_ = registry.Register("file", provider)

	ctx := context.Background()
	resources, _ := registry.List(ctx)
	fmt.Println("Count:", len(resources))
	// Output:
	// Count: 2
}

func ExampleRegistry_Templates() {
	registry := resource.NewRegistry()
	provider := resource.NewStaticProvider()

	provider.AddTemplate(resource.Template{
		URITemplate: "file:///users/{id}/profile",
		Name:        "User Profile",
	})
	_ = registry.Register("file", provider)

	ctx := context.Background()
	templates, _ := registry.Templates(ctx)
	fmt.Println("Count:", len(templates))
	fmt.Println("Template:", templates[0].URITemplate)
	// Output:
	// Count: 1
	// Template: file:///users/{id}/profile
}

func ExampleNewStaticProvider() {
	provider := resource.NewStaticProvider()

	fmt.Printf("Type: %T\n", provider)
	// Output:
	// Type: *resource.StaticProvider
}

func ExampleStaticProvider_Add() {
	provider := resource.NewStaticProvider()

	provider.Add(
		resource.Resource{
			URI:         "file:///readme.md",
			Name:        "README",
			Description: "Project readme",
			MIMEType:    "text/markdown",
		},
		resource.Contents{
			URI:      "file:///readme.md",
			MIMEType: "text/markdown",
			Text:     "# Project\n\nWelcome!",
		},
	)

	ctx := context.Background()
	resources, _ := provider.List(ctx)
	fmt.Println("Count:", len(resources))
	// Output:
	// Count: 1
}

func ExampleStaticProvider_Remove() {
	provider := resource.NewStaticProvider()
	provider.Add(
		resource.Resource{URI: "file:///temp.txt"},
		resource.Contents{URI: "file:///temp.txt", Text: "temp"},
	)

	provider.Remove("file:///temp.txt")

	ctx := context.Background()
	resources, _ := provider.List(ctx)
	fmt.Println("Count:", len(resources))
	// Output:
	// Count: 0
}

func ExampleStaticProvider_Read() {
	provider := resource.NewStaticProvider()
	provider.Add(
		resource.Resource{URI: "file:///doc.txt"},
		resource.Contents{URI: "file:///doc.txt", Text: "Hello!"},
	)

	ctx := context.Background()
	contents, err := provider.Read(ctx, "file:///doc.txt")
	fmt.Println("Error:", err)
	fmt.Println("Text:", contents.Text)
	// Output:
	// Error: <nil>
	// Text: Hello!
}

func ExampleStaticProvider_Read_notFound() {
	provider := resource.NewStaticProvider()
	ctx := context.Background()

	_, err := provider.Read(ctx, "file:///missing.txt")
	fmt.Println("Is not found:", errors.Is(err, resource.ErrResourceNotFound))
	// Output:
	// Is not found: true
}

func ExampleResource_Clone() {
	res := &resource.Resource{
		URI:         "file:///doc.txt",
		Name:        "Document",
		Description: "A document",
		MIMEType:    "text/plain",
		Annotations: map[string]any{"key": "value"},
	}

	clone := res.Clone()
	fmt.Println("URI matches:", clone.URI == res.URI)
	fmt.Println("Name matches:", clone.Name == res.Name)
	// Output:
	// URI matches: true
	// Name matches: true
}

func ExampleContents_IsText() {
	textContent := &resource.Contents{Text: "Hello"}
	binaryContent := &resource.Contents{Blob: []byte{0x01, 0x02}}

	fmt.Println("Text is text:", textContent.IsText())
	fmt.Println("Binary is text:", binaryContent.IsText())
	// Output:
	// Text is text: true
	// Binary is text: false
}

func ExampleContents_IsBinary() {
	textContent := &resource.Contents{Text: "Hello"}
	binaryContent := &resource.Contents{Blob: []byte{0x01, 0x02}}

	fmt.Println("Text is binary:", textContent.IsBinary())
	fmt.Println("Binary is binary:", binaryContent.IsBinary())
	// Output:
	// Text is binary: false
	// Binary is binary: true
}

func ExampleTemplate_Expand() {
	tmpl := &resource.Template{
		URITemplate: "file:///users/{userId}/docs/{docId}",
		Name:        "User Document",
	}

	uri := tmpl.Expand(map[string]string{
		"userId": "123",
		"docId":  "456",
	})
	fmt.Println("Expanded:", uri)
	// Output:
	// Expanded: file:///users/123/docs/456
}

func ExampleNewSubscriptionManager() {
	mgr := resource.NewSubscriptionManager()

	fmt.Printf("Type: %T\n", mgr)
	// Output:
	// Type: *resource.SubscriptionManager
}

func ExampleNewSubscriptionManager_withBufferSize() {
	mgr := resource.NewSubscriptionManager(resource.WithBufferSize(100))

	fmt.Printf("Type: %T\n", mgr)
	// Output:
	// Type: *resource.SubscriptionManager
}

func ExampleSubscriptionManager_Subscribe() {
	mgr := resource.NewSubscriptionManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := mgr.Subscribe(ctx, "file:///config.json")
	fmt.Println("Error:", err)
	fmt.Printf("Channel type: %T\n", ch)
	// Output:
	// Error: <nil>
	// Channel type: <-chan *resource.Contents
}

func ExampleSubscriptionManager_Notify() {
	mgr := resource.NewSubscriptionManager()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, _ := mgr.Subscribe(ctx, "file:///config.json")

	// Notify subscribers
	mgr.Notify("file:///config.json", &resource.Contents{
		URI:  "file:///config.json",
		Text: `{"updated": true}`,
	})

	// Receive update
	contents := <-ch
	fmt.Println("Received:", contents.Text)
	// Output:
	// Received: {"updated": true}
}

func ExampleSubscriptionManager_Unsubscribe() {
	mgr := resource.NewSubscriptionManager()
	ctx := context.Background()

	_, _ = mgr.Subscribe(ctx, "file:///config.json")

	err := mgr.Unsubscribe(ctx, "file:///config.json")
	fmt.Println("Unsubscribe error:", err)
	// Output:
	// Unsubscribe error: <nil>
}

func ExampleResourceError() {
	err := &resource.ResourceError{
		URI: "file:///missing.txt",
		Op:  "read",
		Err: resource.ErrResourceNotFound,
	}

	fmt.Println(err.Error())
	fmt.Println("Unwraps:", errors.Is(err, resource.ErrResourceNotFound))
	// Output:
	// resource file:///missing.txt: read: resource: not found
	// Unwraps: true
}

func Example_resourceWorkflow() {
	// Create registry and provider
	registry := resource.NewRegistry()
	provider := resource.NewStaticProvider()

	// Add resources
	provider.Add(
		resource.Resource{
			URI:      "file:///config.json",
			Name:     "Configuration",
			MIMEType: "application/json",
		},
		resource.Contents{
			URI:      "file:///config.json",
			MIMEType: "application/json",
			Text:     `{"debug": true}`,
		},
	)

	// Register provider
	_ = registry.Register("file", provider)

	// Read resource
	ctx := context.Background()
	contents, _ := registry.Read(ctx, "file:///config.json")
	fmt.Println("Config:", contents.Text)
	// Output:
	// Config: {"debug": true}
}
