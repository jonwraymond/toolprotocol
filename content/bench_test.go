package content

import (
	"testing"
)

// BenchmarkNewText measures text content creation performance.
func BenchmarkNewText(b *testing.B) {
	text := "Hello, world! This is a test message."

	b.ResetTimer()
	for b.Loop() {
		_ = NewText(text)
	}
}

// BenchmarkTextContent_Bytes measures text bytes retrieval.
func BenchmarkTextContent_Bytes(b *testing.B) {
	text := NewText("Hello, world!")

	b.ResetTimer()
	for b.Loop() {
		_, _ = text.Bytes()
	}
}

// BenchmarkTextContent_String measures text string conversion.
func BenchmarkTextContent_String(b *testing.B) {
	text := NewText("Hello, world!")

	b.ResetTimer()
	for b.Loop() {
		_ = text.String()
	}
}

// BenchmarkNewImage measures image content creation performance.
func BenchmarkNewImage(b *testing.B) {
	data := make([]byte, 1024)

	b.ResetTimer()
	for b.Loop() {
		_ = NewImage(data, "image/png")
	}
}

// BenchmarkImageContent_Bytes measures image bytes retrieval.
func BenchmarkImageContent_Bytes(b *testing.B) {
	img := NewImage(make([]byte, 1024), "image/png")

	b.ResetTimer()
	for b.Loop() {
		_, _ = img.Bytes()
	}
}

// BenchmarkImageContent_String measures image base64 encoding.
func BenchmarkImageContent_String(b *testing.B) {
	img := NewImage(make([]byte, 1024), "image/png")

	b.ResetTimer()
	for b.Loop() {
		_ = img.String()
	}
}

// BenchmarkNewResource measures resource content creation.
func BenchmarkNewResource(b *testing.B) {
	uri := "file:///path/to/resource.txt"

	b.ResetTimer()
	for b.Loop() {
		_ = NewResource(uri)
	}
}

// BenchmarkResourceContent_Bytes measures resource bytes retrieval.
func BenchmarkResourceContent_Bytes(b *testing.B) {
	resource := &ResourceContent{
		URI:  "file:///path",
		Text: "resource content text",
	}

	b.ResetTimer()
	for b.Loop() {
		_, _ = resource.Bytes()
	}
}

// BenchmarkNewAudio measures audio content creation.
func BenchmarkNewAudio(b *testing.B) {
	data := make([]byte, 1024)

	b.ResetTimer()
	for b.Loop() {
		_ = NewAudio(data, "audio/mpeg")
	}
}

// BenchmarkAudioContent_String measures audio base64 encoding.
func BenchmarkAudioContent_String(b *testing.B) {
	audio := NewAudio(make([]byte, 1024), "audio/mpeg")

	b.ResetTimer()
	for b.Loop() {
		_ = audio.String()
	}
}

// BenchmarkNewFile measures file content creation.
func BenchmarkNewFile(b *testing.B) {
	data := make([]byte, 1024)

	b.ResetTimer()
	for b.Loop() {
		_ = NewFile(data, "application/pdf")
	}
}

// BenchmarkFileContent_Bytes measures file bytes retrieval.
func BenchmarkFileContent_Bytes(b *testing.B) {
	file := NewFile(make([]byte, 1024), "application/pdf")

	b.ResetTimer()
	for b.Loop() {
		_, _ = file.Bytes()
	}
}

// BenchmarkNewBuilder measures builder creation.
func BenchmarkNewBuilder(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewBuilder()
	}
}

// BenchmarkBuilder_Text measures builder text creation.
func BenchmarkBuilder_Text(b *testing.B) {
	builder := NewBuilder()

	b.ResetTimer()
	for b.Loop() {
		_ = builder.Text("Hello, world!")
	}
}

// BenchmarkBuilder_Image measures builder image creation.
func BenchmarkBuilder_Image(b *testing.B) {
	builder := NewBuilder()
	data := make([]byte, 1024)

	b.ResetTimer()
	for b.Loop() {
		_ = builder.Image(data, "image/png")
	}
}

// BenchmarkBuilder_Chain measures chained builder calls.
func BenchmarkBuilder_Chain(b *testing.B) {
	builder := NewBuilder()
	data := make([]byte, 256)

	b.ResetTimer()
	for b.Loop() {
		_ = builder.Text("text")
		_ = builder.Image(data, "image/png")
		_ = builder.Resource("file:///path")
		_ = builder.Audio(data, "audio/mpeg")
		_ = builder.File(data, "text/plain")
	}
}

// BenchmarkContent_Type measures type retrieval across content types.
func BenchmarkContent_Type(b *testing.B) {
	contents := []Content{
		NewText("text"),
		NewImage([]byte{0x01}, "image/png"),
		NewResource("file:///path"),
		NewAudio([]byte{0x01}, "audio/mpeg"),
		NewFile([]byte{0x01}, "text/plain"),
	}

	b.ResetTimer()
	for b.Loop() {
		for _, c := range contents {
			_ = c.Type()
		}
	}
}

// BenchmarkContent_MIMEType measures MIME type retrieval.
func BenchmarkContent_MIMEType(b *testing.B) {
	contents := []Content{
		NewText("text"),
		NewImage([]byte{0x01}, "image/png"),
		NewResource("file:///path"),
		NewAudio([]byte{0x01}, "audio/mpeg"),
		NewFile([]byte{0x01}, "text/plain"),
	}

	b.ResetTimer()
	for b.Loop() {
		for _, c := range contents {
			_ = c.MIMEType()
		}
	}
}
