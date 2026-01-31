package content

import (
	"testing"
)

func TestType_String(t *testing.T) {
	tests := []struct {
		ct   Type
		want string
	}{
		{TypeText, "text"},
		{TypeImage, "image"},
		{TypeResource, "resource"},
		{TypeAudio, "audio"},
		{TypeFile, "file"},
	}

	for _, tt := range tests {
		if got := string(tt.ct); got != tt.want {
			t.Errorf("Type = %q, want %q", got, tt.want)
		}
	}
}

func TestType_Valid(t *testing.T) {
	valid := []Type{TypeText, TypeImage, TypeResource, TypeAudio, TypeFile}

	for _, ct := range valid {
		if ct == "" {
			t.Errorf("Type %q should not be empty", ct)
		}
	}
}

func TestContentInterface_Contract(t *testing.T) {
	var _ Content = (*TextContent)(nil)
	var _ Content = (*ImageContent)(nil)
	var _ Content = (*ResourceContent)(nil)
}

func TestMIMEType_Defaults(t *testing.T) {
	text := &TextContent{}
	if text.MIMEType() != "text/plain" {
		t.Errorf("TextContent default MIMEType = %q, want %q", text.MIMEType(), "text/plain")
	}

	img := &ImageContent{}
	if img.MIMEType() != "application/octet-stream" {
		t.Errorf("ImageContent default MIMEType = %q, want %q", img.MIMEType(), "application/octet-stream")
	}

	res := &ResourceContent{}
	if res.MIMEType() != "application/octet-stream" {
		t.Errorf("ResourceContent default MIMEType = %q, want %q", res.MIMEType(), "application/octet-stream")
	}
}
