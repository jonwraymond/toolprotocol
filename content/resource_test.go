package content

import (
	"testing"
)

func TestResourceContent_Type(t *testing.T) {
	c := NewResource("file:///path/to/file")
	if c.Type() != TypeResource {
		t.Errorf("Type() = %q, want %q", c.Type(), TypeResource)
	}
}

func TestResourceContent_MIMEType(t *testing.T) {
	c := &ResourceContent{
		URI:  "file:///path/to/file.txt",
		Mime: "text/plain",
	}
	if c.MIMEType() != "text/plain" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "text/plain")
	}
}

func TestResourceContent_Bytes_Text(t *testing.T) {
	c := &ResourceContent{
		URI:  "file:///path/to/file",
		Text: "file content",
	}
	data, err := c.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	if string(data) != "file content" {
		t.Errorf("Bytes() = %q, want %q", string(data), "file content")
	}
}

func TestResourceContent_Bytes_Blob(t *testing.T) {
	blob := []byte{0x01, 0x02, 0x03}
	c := &ResourceContent{
		URI:  "file:///path/to/binary",
		Blob: blob,
		Text: "ignored", // Blob takes priority
	}
	data, err := c.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	if len(data) != len(blob) {
		t.Errorf("len(Bytes()) = %d, want %d", len(data), len(blob))
	}
}

func TestResourceContent_String(t *testing.T) {
	c := &ResourceContent{
		URI:  "file:///path/to/file",
		Text: "readable content",
	}
	if c.String() != "readable content" {
		t.Errorf("String() = %q, want %q", c.String(), "readable content")
	}
}

func TestResourceContent_String_FallbackToURI(t *testing.T) {
	c := NewResource("file:///path/to/file")
	if c.String() != "file:///path/to/file" {
		t.Errorf("String() = %q, want %q", c.String(), "file:///path/to/file")
	}
}

func TestResourceContent_URI(t *testing.T) {
	c := NewResource("https://example.com/data.json")
	if c.URI != "https://example.com/data.json" {
		t.Errorf("URI = %q, want %q", c.URI, "https://example.com/data.json")
	}
}

func TestResourceContent_Empty(t *testing.T) {
	c := NewResource("")
	if c.String() != "" {
		t.Errorf("String() = %q, want empty", c.String())
	}
}

func TestResourceContent_BothTextAndBlob(t *testing.T) {
	// When both are set, Blob takes priority for Bytes()
	c := &ResourceContent{
		URI:  "file:///path",
		Text: "text representation",
		Blob: []byte{0x01, 0x02},
	}

	data, _ := c.Bytes()
	if len(data) != 2 {
		t.Errorf("Bytes() should return Blob when both set, got len=%d", len(data))
	}

	// But String() should still return Text
	if c.String() != "text representation" {
		t.Errorf("String() = %q, want %q", c.String(), "text representation")
	}
}
