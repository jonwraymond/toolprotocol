package content

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSON_Text(t *testing.T) {
	c := &TextContent{Text: "hello world", Mime: "text/plain"}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var j map[string]any
	if err := json.Unmarshal(data, &j); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if j["type"] != string(TypeText) {
		t.Errorf("type = %v, want %v", j["type"], TypeText)
	}
	if j["text"] != "hello world" {
		t.Errorf("text = %v, want %v", j["text"], "hello world")
	}
	if j["mimeType"] != "text/plain" {
		t.Errorf("mimeType = %v, want %v", j["mimeType"], "text/plain")
	}
}

func TestMarshalJSON_Text_NoMIME(t *testing.T) {
	c := NewText("simple text")
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var j map[string]any
	if err := json.Unmarshal(data, &j); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if j["type"] != string(TypeText) {
		t.Errorf("type = %v, want %v", j["type"], TypeText)
	}
	// Default MIME type should be text/plain
	if j["mimeType"] != "text/plain" {
		t.Errorf("mimeType = %v, want %v", j["mimeType"], "text/plain")
	}
}

func TestMarshalJSON_Image(t *testing.T) {
	data := []byte{0x89, 0x50, 0x4E, 0x47}
	c := &ImageContent{Data: data, Mime: "image/png", AltText: "test image"}

	jsonData, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var j map[string]any
	if err := json.Unmarshal(jsonData, &j); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if j["type"] != string(TypeImage) {
		t.Errorf("type = %v, want %v", j["type"], TypeImage)
	}
	if j["mimeType"] != "image/png" {
		t.Errorf("mimeType = %v, want %v", j["mimeType"], "image/png")
	}
	if j["altText"] != "test image" {
		t.Errorf("altText = %v, want %v", j["altText"], "test image")
	}
	// Data should be base64 encoded
	if j["data"] != "iVBORw==" {
		t.Errorf("data = %v, want %v", j["data"], "iVBORw==")
	}
}

func TestMarshalJSON_Image_WithURI(t *testing.T) {
	c := &ImageContent{
		Mime: "image/jpeg",
		URI:  "https://example.com/image.jpg",
	}

	jsonData, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var j map[string]any
	if err := json.Unmarshal(jsonData, &j); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if j["uri"] != "https://example.com/image.jpg" {
		t.Errorf("uri = %v, want %v", j["uri"], "https://example.com/image.jpg")
	}
}

func TestMarshalJSON_Resource(t *testing.T) {
	c := &ResourceContent{
		URI:  "file:///path/to/file",
		Mime: "text/plain",
		Text: "file contents",
	}

	jsonData, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var j map[string]any
	if err := json.Unmarshal(jsonData, &j); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if j["type"] != string(TypeResource) {
		t.Errorf("type = %v, want %v", j["type"], TypeResource)
	}
	if j["uri"] != "file:///path/to/file" {
		t.Errorf("uri = %v, want %v", j["uri"], "file:///path/to/file")
	}
	if j["text"] != "file contents" {
		t.Errorf("text = %v, want %v", j["text"], "file contents")
	}
}

func TestMarshalJSON_Resource_WithBlob(t *testing.T) {
	blob := []byte{0x01, 0x02, 0x03}
	c := &ResourceContent{
		URI:  "file:///binary",
		Mime: "application/octet-stream",
		Blob: blob,
	}

	jsonData, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var j map[string]any
	if err := json.Unmarshal(jsonData, &j); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	// Blob should be base64 encoded
	if j["data"] != "AQID" {
		t.Errorf("data = %v, want %v", j["data"], "AQID")
	}
}

func TestUnmarshalJSON_Text(t *testing.T) {
	jsonData := `{"type":"text","mimeType":"text/plain","text":"hello"}`

	var c TextContent
	if err := json.Unmarshal([]byte(jsonData), &c); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if c.Text != "hello" {
		t.Errorf("Text = %q, want %q", c.Text, "hello")
	}
	if c.Mime != "text/plain" {
		t.Errorf("Mime = %q, want %q", c.Mime, "text/plain")
	}
}

func TestUnmarshalJSON_Image(t *testing.T) {
	jsonData := `{"type":"image","mimeType":"image/png","data":"iVBORw==","altText":"test"}`

	var c ImageContent
	if err := json.Unmarshal([]byte(jsonData), &c); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if c.Mime != "image/png" {
		t.Errorf("Mime = %q, want %q", c.Mime, "image/png")
	}
	if c.AltText != "test" {
		t.Errorf("AltText = %q, want %q", c.AltText, "test")
	}
	// Data should be decoded from base64
	if len(c.Data) != 4 {
		t.Errorf("len(Data) = %d, want 4", len(c.Data))
	}
}

func TestUnmarshalJSON_Image_WithURI(t *testing.T) {
	jsonData := `{"type":"image","mimeType":"image/jpeg","uri":"https://example.com/img.jpg"}`

	var c ImageContent
	if err := json.Unmarshal([]byte(jsonData), &c); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if c.URI != "https://example.com/img.jpg" {
		t.Errorf("URI = %q, want %q", c.URI, "https://example.com/img.jpg")
	}
}

func TestUnmarshalJSON_Image_InvalidBase64(t *testing.T) {
	jsonData := `{"type":"image","mimeType":"image/png","data":"!!!invalid!!!"}`

	var c ImageContent
	err := json.Unmarshal([]byte(jsonData), &c)
	if err == nil {
		t.Fatal("expected error for invalid base64, got nil")
	}
}

func TestUnmarshalJSON_Resource(t *testing.T) {
	jsonData := `{"type":"resource","mimeType":"text/plain","uri":"file:///test","text":"contents"}`

	var c ResourceContent
	if err := json.Unmarshal([]byte(jsonData), &c); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if c.URI != "file:///test" {
		t.Errorf("URI = %q, want %q", c.URI, "file:///test")
	}
	if c.Text != "contents" {
		t.Errorf("Text = %q, want %q", c.Text, "contents")
	}
}

func TestUnmarshalJSON_Resource_WithBlob(t *testing.T) {
	jsonData := `{"type":"resource","mimeType":"application/octet-stream","uri":"file:///bin","data":"AQID"}`

	var c ResourceContent
	if err := json.Unmarshal([]byte(jsonData), &c); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if len(c.Blob) != 3 {
		t.Errorf("len(Blob) = %d, want 3", len(c.Blob))
	}
	if c.Blob[0] != 0x01 || c.Blob[1] != 0x02 || c.Blob[2] != 0x03 {
		t.Errorf("Blob = %v, want [1 2 3]", c.Blob)
	}
}

func TestUnmarshalJSON_Resource_InvalidBase64(t *testing.T) {
	jsonData := `{"type":"resource","uri":"file:///test","data":"!!!invalid!!!"}`

	var c ResourceContent
	err := json.Unmarshal([]byte(jsonData), &c)
	if err == nil {
		t.Fatal("expected error for invalid base64, got nil")
	}
}

func TestBase64_Encode(t *testing.T) {
	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte{}, ""},
		{[]byte{0x00}, "AA=="},
		{[]byte{0x01, 0x02, 0x03}, "AQID"},
		{[]byte("hello"), "aGVsbG8="},
	}

	for _, tt := range tests {
		got := EncodeBase64(tt.input)
		if got != tt.want {
			t.Errorf("EncodeBase64(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBase64_Decode(t *testing.T) {
	tests := []struct {
		input   string
		want    []byte
		wantErr bool
	}{
		{"", []byte{}, false},
		{"AA==", []byte{0x00}, false},
		{"AQID", []byte{0x01, 0x02, 0x03}, false},
		{"aGVsbG8=", []byte("hello"), false},
		{"!!!invalid!!!", nil, true},
	}

	for _, tt := range tests {
		got, err := DecodeBase64(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("DecodeBase64(%q) error = %v, wantErr = %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && string(got) != string(tt.want) {
			t.Errorf("DecodeBase64(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestRoundTrip_Text(t *testing.T) {
	original := &TextContent{Text: "round trip test", Mime: "text/plain"}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var restored TextContent
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if restored.Text != original.Text {
		t.Errorf("Text = %q, want %q", restored.Text, original.Text)
	}
	if restored.Mime != original.Mime {
		t.Errorf("Mime = %q, want %q", restored.Mime, original.Mime)
	}
}

func TestRoundTrip_Image(t *testing.T) {
	original := &ImageContent{
		Data:    []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
		Mime:    "image/png",
		URI:     "https://example.com/img.png",
		AltText: "PNG image",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var restored ImageContent
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if string(restored.Data) != string(original.Data) {
		t.Errorf("Data mismatch")
	}
	if restored.Mime != original.Mime {
		t.Errorf("Mime = %q, want %q", restored.Mime, original.Mime)
	}
	if restored.URI != original.URI {
		t.Errorf("URI = %q, want %q", restored.URI, original.URI)
	}
	if restored.AltText != original.AltText {
		t.Errorf("AltText = %q, want %q", restored.AltText, original.AltText)
	}
}

func TestRoundTrip_Resource(t *testing.T) {
	original := &ResourceContent{
		URI:  "file:///test/resource",
		Mime: "application/json",
		Text: `{"key": "value"}`,
		Blob: []byte{0x01, 0x02},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var restored ResourceContent
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if restored.URI != original.URI {
		t.Errorf("URI = %q, want %q", restored.URI, original.URI)
	}
	if restored.Mime != original.Mime {
		t.Errorf("Mime = %q, want %q", restored.Mime, original.Mime)
	}
	if restored.Text != original.Text {
		t.Errorf("Text = %q, want %q", restored.Text, original.Text)
	}
	if string(restored.Blob) != string(original.Blob) {
		t.Errorf("Blob mismatch")
	}
}

func TestUnmarshalJSON_Text_Invalid(t *testing.T) {
	var c TextContent
	err := json.Unmarshal([]byte("not json"), &c)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestUnmarshalJSON_Image_Invalid(t *testing.T) {
	var c ImageContent
	err := json.Unmarshal([]byte("not json"), &c)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestUnmarshalJSON_Resource_Invalid(t *testing.T) {
	var c ResourceContent
	err := json.Unmarshal([]byte("not json"), &c)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
