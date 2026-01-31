package content

import (
	"testing"
)

func TestTextContent_Type(t *testing.T) {
	c := NewText("hello")
	if c.Type() != TypeText {
		t.Errorf("Type() = %q, want %q", c.Type(), TypeText)
	}
}

func TestTextContent_MIMEType_Default(t *testing.T) {
	c := NewText("hello")
	if c.MIMEType() != "text/plain" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "text/plain")
	}
}

func TestTextContent_MIMEType_Custom(t *testing.T) {
	c := &TextContent{
		Text: "# Markdown",
		Mime: "text/markdown",
	}
	if c.MIMEType() != "text/markdown" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "text/markdown")
	}
}

func TestTextContent_Bytes(t *testing.T) {
	c := NewText("hello")
	data, err := c.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("Bytes() = %q, want %q", string(data), "hello")
	}
}

func TestTextContent_String(t *testing.T) {
	c := NewText("hello, world!")
	if c.String() != "hello, world!" {
		t.Errorf("String() = %q, want %q", c.String(), "hello, world!")
	}
}

func TestTextContent_Empty(t *testing.T) {
	c := NewText("")
	if c.String() != "" {
		t.Errorf("String() = %q, want empty", c.String())
	}
	data, _ := c.Bytes()
	if len(data) != 0 {
		t.Errorf("len(Bytes()) = %d, want 0", len(data))
	}
}

func TestTextContent_Unicode(t *testing.T) {
	c := NewText("Hello, 世界! 🌍")
	data, err := c.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	if string(data) != "Hello, 世界! 🌍" {
		t.Errorf("Unicode roundtrip failed")
	}
}
