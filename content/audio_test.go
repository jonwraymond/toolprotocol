package content

import (
	"testing"
	"time"
)

func TestAudioContent_Type(t *testing.T) {
	c := NewAudio([]byte{0x00}, "audio/mpeg")
	if c.Type() != TypeAudio {
		t.Errorf("Type() = %q, want %q", c.Type(), TypeAudio)
	}
}

func TestAudioContent_MIMEType(t *testing.T) {
	c := NewAudio([]byte{0x00}, "audio/wav")
	if c.MIMEType() != "audio/wav" {
		t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), "audio/wav")
	}
}

func TestAudioContent_MIMEType_Default(t *testing.T) {
	c := &AudioContent{Data: []byte{0x00}}
	if c.MIMEType() != "audio/mpeg" {
		t.Errorf("MIMEType() default = %q, want %q", c.MIMEType(), "audio/mpeg")
	}
}

func TestAudioContent_Bytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	c := NewAudio(data, "audio/mpeg")
	result, err := c.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	if len(result) != len(data) {
		t.Errorf("len(Bytes()) = %d, want %d", len(result), len(data))
	}
}

func TestAudioContent_Duration(t *testing.T) {
	c := &AudioContent{
		Data:     []byte{0x00},
		Mime:     "audio/mpeg",
		Duration: 30 * time.Second,
	}
	if c.Duration != 30*time.Second {
		t.Errorf("Duration = %v, want %v", c.Duration, 30*time.Second)
	}
}

func TestAudioContent_Formats(t *testing.T) {
	formats := []string{"audio/mpeg", "audio/wav", "audio/ogg", "audio/flac"}
	for _, format := range formats {
		c := NewAudio([]byte{0x00}, format)
		if c.MIMEType() != format {
			t.Errorf("MIMEType() = %q, want %q", c.MIMEType(), format)
		}
	}
}
