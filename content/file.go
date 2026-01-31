package content

// FileContent represents file data.
type FileContent struct {
	// Data is the raw file bytes.
	Data []byte

	// Mime is the MIME type.
	Mime string

	// Path is the file path.
	Path string

	// Size is the file size in bytes (may be set independently of Data).
	Size int64
}

// NewFile creates a new file content.
func NewFile(data []byte, mimeType string) *FileContent {
	return &FileContent{
		Data: data,
		Mime: mimeType,
		Size: int64(len(data)),
	}
}

// Type returns TypeFile.
func (c *FileContent) Type() Type {
	return TypeFile
}

// MIMEType returns the file MIME type.
func (c *FileContent) MIMEType() string {
	if c.Mime == "" {
		return "application/octet-stream"
	}
	return c.Mime
}

// Bytes returns the file data.
func (c *FileContent) Bytes() ([]byte, error) {
	return c.Data, nil
}

// String returns the file path or a placeholder.
func (c *FileContent) String() string {
	if c.Path != "" {
		return c.Path
	}
	return "[file data]"
}
