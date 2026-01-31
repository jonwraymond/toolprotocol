package content

import (
	"encoding/base64"
	"testing"
)

func TestImageContent_Type(t *testing.T) {
	c := NewImage([]byte{0x89, 0x50}, "image/png")
	if c.Type() != TypeImage {
		t.Errorf("Type() = %q, want %q", c.Type(), TypeImage)
	}
}

func TestImageContent_MIMEType(t *testing.T) {
	c := NewImage([]byte{0x89}, "image/png")
	if c.MIMEType() != "image/png" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "image/png")
	}
}

func TestImageContent_Bytes(t *testing.T) {
	data := []byte{0x89, 0x50, 0x4e, 0x47}
	c := NewImage(data, "image/png")
	result, err := c.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	if len(result) != len(data) {
		t.Errorf("len(Bytes()) = %d, want %d", len(result), len(data))
	}
}

func TestImageContent_String_Base64(t *testing.T) {
	data := []byte{0x89, 0x50, 0x4e, 0x47}
	c := NewImage(data, "image/png")
	expected := base64.StdEncoding.EncodeToString(data)
	if c.String() != expected {
		t.Errorf("String() = %q, want %q", c.String(), expected)
	}
}

func TestImageContent_WithURI(t *testing.T) {
	c := &ImageContent{
		Data: []byte{0x89},
		Mime: "image/png",
		URI:  "https://example.com/image.png",
	}
	if c.URI != "https://example.com/image.png" {
		t.Errorf("URI = %q, want %q", c.URI, "https://example.com/image.png")
	}
}

func TestImageContent_WithAltText(t *testing.T) {
	c := &ImageContent{
		Data:    []byte{0x89},
		Mime:    "image/png",
		AltText: "A test image",
	}
	if c.AltText != "A test image" {
		t.Errorf("AltText = %q, want %q", c.AltText, "A test image")
	}
}

func TestImageContent_PNG(t *testing.T) {
	// PNG magic bytes
	pngData := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	c := NewImage(pngData, "image/png")
	if c.MIMEType() != "image/png" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "image/png")
	}
}

func TestImageContent_JPEG(t *testing.T) {
	// JPEG magic bytes
	jpegData := []byte{0xff, 0xd8, 0xff, 0xe0}
	c := NewImage(jpegData, "image/jpeg")
	if c.MIMEType() != "image/jpeg" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "image/jpeg")
	}
}

func TestImageContent_Empty(t *testing.T) {
	c := NewImage(nil, "image/png")
	data, _ := c.Bytes()
	if len(data) != 0 {
		t.Errorf("len(Bytes()) = %d, want 0", len(data))
	}
	if c.String() != "" {
		t.Errorf("String() = %q, want empty", c.String())
	}
}
