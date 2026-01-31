package content

import "encoding/base64"

// ImageContent represents image data.
type ImageContent struct {
	// Data is the raw image bytes.
	Data []byte

	// Mime is the MIME type (e.g., image/png, image/jpeg).
	Mime string

	// URI is an optional source URI.
	URI string

	// AltText is accessibility text.
	AltText string
}

// NewImage creates a new image content.
func NewImage(data []byte, mimeType string) *ImageContent {
	return &ImageContent{
		Data: data,
		Mime: mimeType,
	}
}

// Type returns TypeImage.
func (c *ImageContent) Type() Type {
	return TypeImage
}

// MIMEType returns the image MIME type.
func (c *ImageContent) MIMEType() string {
	if c.Mime == "" {
		return "application/octet-stream"
	}
	return c.Mime
}

// Bytes returns the image data.
func (c *ImageContent) Bytes() ([]byte, error) {
	return c.Data, nil
}

// String returns a base64 representation of the image.
func (c *ImageContent) String() string {
	return base64.StdEncoding.EncodeToString(c.Data)
}
