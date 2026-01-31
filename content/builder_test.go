package content

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	b := NewBuilder()
	if b == nil {
		t.Fatal("NewBuilder() returned nil")
	}
}

func TestBuilder_Text(t *testing.T) {
	b := NewBuilder()
	c := b.Text("hello")
	if c == nil {
		t.Fatal("Text() returned nil")
	}
	if c.Text != "hello" {
		t.Errorf("Text = %q, want %q", c.Text, "hello")
	}
	if c.Type() != TypeText {
		t.Errorf("Type() = %q, want %q", c.Type(), TypeText)
	}
}

func TestBuilder_TextWithMIME(t *testing.T) {
	b := NewBuilder()
	c := b.TextWithMIME("data", "application/json")
	if c == nil {
		t.Fatal("TextWithMIME() returned nil")
	}
	if c.Text != "data" {
		t.Errorf("Text = %q, want %q", c.Text, "data")
	}
	if c.MIMEType() != "application/json" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "application/json")
	}
}

func TestBuilder_Image(t *testing.T) {
	b := NewBuilder()
	data := []byte{0x89, 0x50, 0x4E, 0x47}
	c := b.Image(data, "image/png")
	if c == nil {
		t.Fatal("Image() returned nil")
	}
	if len(c.Data) != len(data) {
		t.Errorf("len(Data) = %d, want %d", len(c.Data), len(data))
	}
	if c.MIMEType() != "image/png" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "image/png")
	}
}

func TestBuilder_ImageWithAlt(t *testing.T) {
	b := NewBuilder()
	data := []byte{0x89, 0x50, 0x4E, 0x47}
	c := b.ImageWithAlt(data, "image/png", "A test image")
	if c == nil {
		t.Fatal("ImageWithAlt() returned nil")
	}
	if c.AltText != "A test image" {
		t.Errorf("AltText = %q, want %q", c.AltText, "A test image")
	}
	if c.MIMEType() != "image/png" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "image/png")
	}
}

func TestBuilder_Resource(t *testing.T) {
	b := NewBuilder()
	c := b.Resource("file:///path/to/file")
	if c == nil {
		t.Fatal("Resource() returned nil")
	}
	if c.URI != "file:///path/to/file" {
		t.Errorf("URI = %q, want %q", c.URI, "file:///path/to/file")
	}
}

func TestBuilder_ResourceWithText(t *testing.T) {
	b := NewBuilder()
	c := b.ResourceWithText("file:///test.txt", "file contents")
	if c == nil {
		t.Fatal("ResourceWithText() returned nil")
	}
	if c.URI != "file:///test.txt" {
		t.Errorf("URI = %q, want %q", c.URI, "file:///test.txt")
	}
	if c.Text != "file contents" {
		t.Errorf("Text = %q, want %q", c.Text, "file contents")
	}
}

func TestBuilder_Audio(t *testing.T) {
	b := NewBuilder()
	data := []byte{0xFF, 0xFB, 0x90}
	c := b.Audio(data, "audio/mpeg")
	if c == nil {
		t.Fatal("Audio() returned nil")
	}
	if len(c.Data) != len(data) {
		t.Errorf("len(Data) = %d, want %d", len(c.Data), len(data))
	}
	if c.MIMEType() != "audio/mpeg" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "audio/mpeg")
	}
}

func TestBuilder_File(t *testing.T) {
	b := NewBuilder()
	data := []byte{0x25, 0x50, 0x44, 0x46}
	c := b.File(data, "application/pdf")
	if c == nil {
		t.Fatal("File() returned nil")
	}
	if len(c.Data) != len(data) {
		t.Errorf("len(Data) = %d, want %d", len(c.Data), len(data))
	}
	if c.MIMEType() != "application/pdf" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "application/pdf")
	}
}

func TestBuilder_FileWithPath(t *testing.T) {
	b := NewBuilder()
	data := []byte{0x25, 0x50, 0x44, 0x46}
	c := b.FileWithPath(data, "application/pdf", "/docs/report.pdf")
	if c == nil {
		t.Fatal("FileWithPath() returned nil")
	}
	if c.Path != "/docs/report.pdf" {
		t.Errorf("Path = %q, want %q", c.Path, "/docs/report.pdf")
	}
	if c.Size != int64(len(data)) {
		t.Errorf("Size = %d, want %d", c.Size, len(data))
	}
}

func TestBuilder_Chaining(t *testing.T) {
	b := NewBuilder()

	// Verify builder can create multiple different content types
	text := b.Text("text content")
	image := b.Image([]byte{0x00}, "image/png")
	resource := b.Resource("file:///test")

	if text.Type() != TypeText {
		t.Errorf("text.Type() = %q, want %q", text.Type(), TypeText)
	}
	if image.Type() != TypeImage {
		t.Errorf("image.Type() = %q, want %q", image.Type(), TypeImage)
	}
	if resource.Type() != TypeResource {
		t.Errorf("resource.Type() = %q, want %q", resource.Type(), TypeResource)
	}
}
