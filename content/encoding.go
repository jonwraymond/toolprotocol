package content

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// contentJSON is the JSON representation of content.
type contentJSON struct {
	Type     Type   `json:"type"`
	MIMEType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"` // base64 encoded
	URI      string `json:"uri,omitempty"`
	AltText  string `json:"altText,omitempty"`
	Path     string `json:"path,omitempty"`
}

// MarshalJSON marshals text content to JSON.
func (c *TextContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(contentJSON{
		Type:     TypeText,
		MIMEType: c.MIMEType(),
		Text:     c.Text,
	})
}

// UnmarshalJSON unmarshals JSON to text content.
func (c *TextContent) UnmarshalJSON(data []byte) error {
	var j contentJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	c.Text = j.Text
	c.Mime = j.MIMEType
	return nil
}

// MarshalJSON marshals image content to JSON.
func (c *ImageContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(contentJSON{
		Type:     TypeImage,
		MIMEType: c.MIMEType(),
		Data:     base64.StdEncoding.EncodeToString(c.Data),
		URI:      c.URI,
		AltText:  c.AltText,
	})
}

// UnmarshalJSON unmarshals JSON to image content.
func (c *ImageContent) UnmarshalJSON(data []byte) error {
	var j contentJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	if j.Data != "" {
		decoded, err := base64.StdEncoding.DecodeString(j.Data)
		if err != nil {
			return fmt.Errorf("decode image data: %w", err)
		}
		c.Data = decoded
	}
	c.Mime = j.MIMEType
	c.URI = j.URI
	c.AltText = j.AltText
	return nil
}

// MarshalJSON marshals resource content to JSON.
func (c *ResourceContent) MarshalJSON() ([]byte, error) {
	j := contentJSON{
		Type:     TypeResource,
		MIMEType: c.MIMEType(),
		URI:      c.URI,
		Text:     c.Text,
	}
	if len(c.Blob) > 0 {
		j.Data = base64.StdEncoding.EncodeToString(c.Blob)
	}
	return json.Marshal(j)
}

// UnmarshalJSON unmarshals JSON to resource content.
func (c *ResourceContent) UnmarshalJSON(data []byte) error {
	var j contentJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	c.URI = j.URI
	c.Mime = j.MIMEType
	c.Text = j.Text
	if j.Data != "" {
		decoded, err := base64.StdEncoding.DecodeString(j.Data)
		if err != nil {
			return fmt.Errorf("decode resource data: %w", err)
		}
		c.Blob = decoded
	}
	return nil
}

// EncodeBase64 encodes bytes to base64 string.
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes base64 string to bytes.
func DecodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
