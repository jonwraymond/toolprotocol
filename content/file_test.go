package content

import (
	"testing"
)

func TestFileContent_Type(t *testing.T) {
	c := NewFile([]byte{0x00}, "application/pdf")
	if c.Type() != TypeFile {
		t.Errorf("Type() = %q, want %q", c.Type(), TypeFile)
	}
}

func TestFileContent_MIMEType(t *testing.T) {
	c := NewFile([]byte{0x00}, "application/pdf")
	if c.MIMEType() != "application/pdf" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "application/pdf")
	}
}

func TestFileContent_MIMEType_Default(t *testing.T) {
	c := &FileContent{Data: []byte{0x00}}
	if c.MIMEType() != "application/octet-stream" {
		t.Errorf("MIMEType() default = %q, want %q", c.MIMEType(), "application/octet-stream")
	}
}

func TestFileContent_Bytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	c := NewFile(data, "application/octet-stream")
	result, err := c.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	if len(result) != len(data) {
		t.Errorf("len(Bytes()) = %d, want %d", len(result), len(data))
	}
}

func TestFileContent_Path(t *testing.T) {
	c := &FileContent{
		Data: []byte{0x00},
		Mime: "text/plain",
		Path: "/path/to/file.txt",
	}
	if c.Path != "/path/to/file.txt" {
		t.Errorf("Path = %q, want %q", c.Path, "/path/to/file.txt")
	}
	if c.String() != "/path/to/file.txt" {
		t.Errorf("String() = %q, want %q", c.String(), "/path/to/file.txt")
	}
}

func TestFileContent_Size(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	c := NewFile(data, "application/octet-stream")
	if c.Size != 5 {
		t.Errorf("Size = %d, want %d", c.Size, 5)
	}
}

func TestFileContent_Empty(t *testing.T) {
	c := NewFile(nil, "application/octet-stream")
	data, _ := c.Bytes()
	if len(data) != 0 {
		t.Errorf("len(Bytes()) = %d, want 0", len(data))
	}
	if c.String() != "[file data]" {
		t.Errorf("String() = %q, want %q", c.String(), "[file data]")
	}
}
