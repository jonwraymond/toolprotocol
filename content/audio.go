package content

import (
	"encoding/base64"
	"time"
)

// AudioContent represents audio data.
type AudioContent struct {
	// Data is the raw audio bytes.
	Data []byte

	// Mime is the MIME type (e.g., audio/mpeg, audio/wav).
	Mime string

	// Duration is the audio duration.
	Duration time.Duration
}

// NewAudio creates a new audio content.
func NewAudio(data []byte, mimeType string) *AudioContent {
	return &AudioContent{
		Data: data,
		Mime: mimeType,
	}
}

// Type returns TypeAudio.
func (c *AudioContent) Type() Type {
	return TypeAudio
}

// MIMEType returns the audio MIME type.
func (c *AudioContent) MIMEType() string {
	if c.Mime == "" {
		return "audio/mpeg"
	}
	return c.Mime
}

// Bytes returns the audio data.
func (c *AudioContent) Bytes() ([]byte, error) {
	return c.Data, nil
}

// String returns a base64 representation of the audio.
func (c *AudioContent) String() string {
	return base64.StdEncoding.EncodeToString(c.Data)
}
