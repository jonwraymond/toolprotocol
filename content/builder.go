package content

// Builder creates content instances.
type Builder struct{}

// NewBuilder creates a new content builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// Text creates a text content.
func (b *Builder) Text(text string) *TextContent {
	return NewText(text)
}

// TextWithMIME creates a text content with a custom MIME type.
func (b *Builder) TextWithMIME(text, mimeType string) *TextContent {
	return &TextContent{
		Text: text,
		Mime: mimeType,
	}
}

// Image creates an image content.
func (b *Builder) Image(data []byte, mimeType string) *ImageContent {
	return NewImage(data, mimeType)
}

// ImageWithAlt creates an image content with alt text.
func (b *Builder) ImageWithAlt(data []byte, mimeType, altText string) *ImageContent {
	return &ImageContent{
		Data:    data,
		Mime:    mimeType,
		AltText: altText,
	}
}

// Resource creates a resource content.
func (b *Builder) Resource(uri string) *ResourceContent {
	return NewResource(uri)
}

// ResourceWithText creates a resource content with text representation.
func (b *Builder) ResourceWithText(uri, text string) *ResourceContent {
	return &ResourceContent{
		URI:  uri,
		Text: text,
	}
}

// Audio creates an audio content.
func (b *Builder) Audio(data []byte, mimeType string) *AudioContent {
	return NewAudio(data, mimeType)
}

// File creates a file content.
func (b *Builder) File(data []byte, mimeType string) *FileContent {
	return NewFile(data, mimeType)
}

// FileWithPath creates a file content with path information.
func (b *Builder) FileWithPath(data []byte, mimeType, path string) *FileContent {
	return &FileContent{
		Data: data,
		Mime: mimeType,
		Path: path,
		Size: int64(len(data)),
	}
}
