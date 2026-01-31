// Package content provides unified content/part abstraction for
// text, images, resources, audio, and files across all protocols.
//
// This package enables consistent handling of response content
// regardless of the underlying protocol (MCP, A2A, ACP).
//
// # Content Interface
//
// All content types implement the Content interface:
//
//	type Content interface {
//	    Type() Type
//	    MIMEType() string
//	    Bytes() ([]byte, error)
//	    String() string
//	}
//
// # Available Content Types
//
//   - TextContent: Plain text content
//   - ImageContent: Image data (PNG, JPEG, etc.)
//   - ResourceContent: External resource references
//   - AudioContent: Audio data
//   - FileContent: File data
//
// # Usage
//
//	// Create text content
//	text := content.NewText("Hello, world!")
//
//	// Create image content
//	img := content.NewImage(pngData, "image/png")
//
//	// Use the builder
//	builder := content.NewBuilder()
//	text := builder.Text("Hello")
//	img := builder.Image(data, "image/png")
package content
