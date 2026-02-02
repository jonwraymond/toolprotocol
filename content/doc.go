// Package content provides unified content/part abstraction for
// text, images, resources, audio, and files across all protocols.
//
// This package enables consistent handling of response content
// regardless of the underlying protocol (MCP, A2A, ACP).
//
// # Ecosystem Position
//
// content provides unified content abstraction for protocol responses:
//
//	┌─────────────────────────────────────────────────────────────────┐
//	│                      Content Flow                               │
//	├─────────────────────────────────────────────────────────────────┤
//	│                                                                 │
//	│   Protocol              content                  Output         │
//	│   ┌─────────┐        ┌───────────┐         ┌─────────┐         │
//	│   │  Tool   │────────│  Builder  │─────────│  wire   │         │
//	│   │Response │ create │           │ encode  │(encode) │         │
//	│   └─────────┘        │ ┌───────┐ │         └─────────┘         │
//	│        │             │ │Content│ │              │               │
//	│        │             │ │       │ │              │               │
//	│        ▼             │ └───────┘ │              │               │
//	│   ┌─────────┐        │     │     │              │               │
//	│   │  Text   │────────│─────┼─────│──────────────▼               │
//	│   │  Image  │        │     │     │         ┌─────────┐         │
//	│   │ Resource│        │     ▼     │         │Transport│         │
//	│   │  Audio  │        │ ┌───────┐ │         │(stream) │         │
//	│   │  File   │        │ │ Types │ │         └─────────┘         │
//	│   └─────────┘        │ └───────┘ │              ▲               │
//	│                      └───────────┘──────────────┘               │
//	│                                                                 │
//	└─────────────────────────────────────────────────────────────────┘
//
// # Content Types
//
// The package provides five content types:
//
//   - [TextContent]: Plain text (text/plain, text/html, etc.)
//   - [ImageContent]: Image data with optional alt text (PNG, JPEG, etc.)
//   - [ResourceContent]: External resource references with optional content
//   - [AudioContent]: Audio data with duration (MP3, WAV, etc.)
//   - [FileContent]: File data with path and size metadata
//
// # Core Components
//
//   - [Content]: Interface implemented by all content types
//   - [Type]: Content type constants (text, image, resource, audio, file)
//   - [Builder]: Factory for creating content instances
//   - [NewText], [NewImage], [NewResource], [NewAudio], [NewFile]: Constructors
//
// # Quick Start
//
//	// Direct construction
//	text := content.NewText("Hello, world!")
//	img := content.NewImage(pngData, "image/png")
//
//	// Using the builder
//	builder := content.NewBuilder()
//	text := builder.Text("Hello")
//	img := builder.ImageWithAlt(data, "image/png", "Photo description")
//	file := builder.FileWithPath(data, "text/csv", "/reports/data.csv")
//
// # Content Interface
//
// All content types implement the Content interface:
//
//	type Content interface {
//	    Type() Type           // Returns content type constant
//	    MIMEType() string     // Returns MIME type (e.g., "text/plain")
//	    Bytes() ([]byte, error) // Returns raw content bytes
//	    String() string       // Returns string representation
//	}
//
// # Builder Pattern
//
// The Builder provides a fluent API for creating content:
//
//	builder := content.NewBuilder()
//
//	// Basic creation
//	text := builder.Text("message")
//	img := builder.Image(data, "image/png")
//
//	// With additional metadata
//	html := builder.TextWithMIME("<h1>Hello</h1>", "text/html")
//	imgAlt := builder.ImageWithAlt(data, "image/jpeg", "Sunset photo")
//	fileWithPath := builder.FileWithPath(data, "text/csv", "/path/to/file.csv")
//
// # String Encoding
//
// Different content types encode differently in String():
//
//   - TextContent: Returns the text content directly
//   - ImageContent: Returns base64-encoded data
//   - ResourceContent: Returns Text if set, otherwise URI
//   - AudioContent: Returns base64-encoded data
//   - FileContent: Returns Path if set, otherwise "[file data]"
//
// # Thread Safety
//
// Content types are designed for concurrent read access:
//
//   - All content types are safe for concurrent reads after creation
//   - Content should be treated as immutable after construction
//   - Builder is stateless and safe for concurrent use
//   - Bytes() returns the underlying slice; do not modify
//
// # Integration with ApertureStack
//
// content integrates with other ApertureStack packages:
//
//   - wire: Content encoded to protocol-specific formats
//   - stream: Content delivered via streaming events
//   - transport: Content transmitted over HTTP/SSE/stdio
package content
