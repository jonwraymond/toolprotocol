package content_test

import (
	"fmt"

	"github.com/jonwraymond/toolprotocol/content"
)

func ExampleNewText() {
	text := content.NewText("Hello, world!")

	fmt.Println("Type:", text.Type())
	fmt.Println("MIME:", text.MIMEType())
	fmt.Println("Content:", text.String())
	// Output:
	// Type: text
	// MIME: text/plain
	// Content: Hello, world!
}

func ExampleTextContent_Bytes() {
	text := content.NewText("Hello")

	data, err := text.Bytes()
	fmt.Println("Error:", err)
	fmt.Println("Bytes:", string(data))
	// Output:
	// Error: <nil>
	// Bytes: Hello
}

func ExampleTextContent_customMIME() {
	text := &content.TextContent{
		Text: "<html>Hello</html>",
		Mime: "text/html",
	}

	fmt.Println("MIME:", text.MIMEType())
	// Output:
	// MIME: text/html
}

func ExampleNewImage() {
	// Create a simple 1x1 PNG (minimal valid PNG)
	pngData := []byte{0x89, 0x50, 0x4E, 0x47}
	img := content.NewImage(pngData, "image/png")

	fmt.Println("Type:", img.Type())
	fmt.Println("MIME:", img.MIMEType())
	// Output:
	// Type: image
	// MIME: image/png
}

func ExampleImageContent_Bytes() {
	data := []byte{0x01, 0x02, 0x03}
	img := content.NewImage(data, "image/png")

	bytes, err := img.Bytes()
	fmt.Println("Error:", err)
	fmt.Println("Length:", len(bytes))
	// Output:
	// Error: <nil>
	// Length: 3
}

func ExampleImageContent_String() {
	// String returns base64 encoding
	data := []byte("test")
	img := content.NewImage(data, "image/png")

	fmt.Println("Base64:", img.String())
	// Output:
	// Base64: dGVzdA==
}

func ExampleImageContent_withAltText() {
	img := &content.ImageContent{
		Data:    []byte{0x01},
		Mime:    "image/jpeg",
		AltText: "A beautiful sunset",
	}

	fmt.Println("Alt text:", img.AltText)
	// Output:
	// Alt text: A beautiful sunset
}

func ExampleNewResource() {
	resource := content.NewResource("file:///path/to/file.txt")

	fmt.Println("Type:", resource.Type())
	fmt.Println("URI:", resource.String())
	// Output:
	// Type: resource
	// URI: file:///path/to/file.txt
}

func ExampleResourceContent_Bytes() {
	resource := &content.ResourceContent{
		URI:  "file:///data.txt",
		Text: "file contents",
	}

	bytes, _ := resource.Bytes()
	fmt.Println("Content:", string(bytes))
	// Output:
	// Content: file contents
}

func ExampleResourceContent_Bytes_blob() {
	resource := &content.ResourceContent{
		URI:  "file:///data.bin",
		Blob: []byte{0x01, 0x02, 0x03},
	}

	bytes, _ := resource.Bytes()
	fmt.Println("Length:", len(bytes))
	// Output:
	// Length: 3
}

func ExampleResourceContent_String() {
	// String prefers Text over URI
	resource := &content.ResourceContent{
		URI:  "file:///path",
		Text: "content text",
	}

	fmt.Println("String:", resource.String())
	// Output:
	// String: content text
}

func ExampleNewAudio() {
	audioData := []byte{0xFF, 0xFB, 0x90, 0x00} // MP3 header
	audio := content.NewAudio(audioData, "audio/mpeg")

	fmt.Println("Type:", audio.Type())
	fmt.Println("MIME:", audio.MIMEType())
	// Output:
	// Type: audio
	// MIME: audio/mpeg
}

func ExampleAudioContent_Bytes() {
	data := []byte{0x01, 0x02}
	audio := content.NewAudio(data, "audio/wav")

	bytes, err := audio.Bytes()
	fmt.Println("Error:", err)
	fmt.Println("Length:", len(bytes))
	// Output:
	// Error: <nil>
	// Length: 2
}

func ExampleAudioContent_String() {
	// String returns base64 encoding
	data := []byte("audio")
	audio := content.NewAudio(data, "audio/mpeg")

	fmt.Println("Base64:", audio.String())
	// Output:
	// Base64: YXVkaW8=
}

func ExampleNewFile() {
	fileData := []byte("file contents")
	file := content.NewFile(fileData, "application/pdf")

	fmt.Println("Type:", file.Type())
	fmt.Println("MIME:", file.MIMEType())
	fmt.Println("Size:", file.Size)
	// Output:
	// Type: file
	// MIME: application/pdf
	// Size: 13
}

func ExampleFileContent_Bytes() {
	data := []byte("hello")
	file := content.NewFile(data, "text/plain")

	bytes, err := file.Bytes()
	fmt.Println("Error:", err)
	fmt.Println("Content:", string(bytes))
	// Output:
	// Error: <nil>
	// Content: hello
}

func ExampleFileContent_String() {
	file := &content.FileContent{
		Data: []byte("content"),
		Mime: "text/plain",
		Path: "/path/to/file.txt",
	}

	fmt.Println("String:", file.String())
	// Output:
	// String: /path/to/file.txt
}

func ExampleFileContent_String_noPath() {
	file := content.NewFile([]byte("data"), "text/plain")

	fmt.Println("String:", file.String())
	// Output:
	// String: [file data]
}

func ExampleNewBuilder() {
	builder := content.NewBuilder()

	fmt.Printf("Type: %T\n", builder)
	// Output:
	// Type: *content.Builder
}

func ExampleBuilder_Text() {
	builder := content.NewBuilder()
	text := builder.Text("Hello, world!")

	fmt.Println("Type:", text.Type())
	fmt.Println("Content:", text.String())
	// Output:
	// Type: text
	// Content: Hello, world!
}

func ExampleBuilder_TextWithMIME() {
	builder := content.NewBuilder()
	text := builder.TextWithMIME("<h1>Hello</h1>", "text/html")

	fmt.Println("MIME:", text.MIMEType())
	// Output:
	// MIME: text/html
}

func ExampleBuilder_Image() {
	builder := content.NewBuilder()
	img := builder.Image([]byte{0x01}, "image/png")

	fmt.Println("Type:", img.Type())
	fmt.Println("MIME:", img.MIMEType())
	// Output:
	// Type: image
	// MIME: image/png
}

func ExampleBuilder_ImageWithAlt() {
	builder := content.NewBuilder()
	img := builder.ImageWithAlt([]byte{0x01}, "image/jpeg", "Sunset photo")

	fmt.Println("Alt:", img.AltText)
	// Output:
	// Alt: Sunset photo
}

func ExampleBuilder_Resource() {
	builder := content.NewBuilder()
	resource := builder.Resource("file:///data.txt")

	fmt.Println("Type:", resource.Type())
	fmt.Println("URI:", resource.String())
	// Output:
	// Type: resource
	// URI: file:///data.txt
}

func ExampleBuilder_ResourceWithText() {
	builder := content.NewBuilder()
	resource := builder.ResourceWithText("file:///doc.txt", "Document contents")

	fmt.Println("String:", resource.String())
	// Output:
	// String: Document contents
}

func ExampleBuilder_Audio() {
	builder := content.NewBuilder()
	audio := builder.Audio([]byte{0xFF, 0xFB}, "audio/mpeg")

	fmt.Println("Type:", audio.Type())
	fmt.Println("MIME:", audio.MIMEType())
	// Output:
	// Type: audio
	// MIME: audio/mpeg
}

func ExampleBuilder_File() {
	builder := content.NewBuilder()
	file := builder.File([]byte("data"), "application/pdf")

	fmt.Println("Type:", file.Type())
	fmt.Println("Size:", file.Size)
	// Output:
	// Type: file
	// Size: 4
}

func ExampleBuilder_FileWithPath() {
	builder := content.NewBuilder()
	file := builder.FileWithPath([]byte("content"), "text/plain", "/docs/readme.txt")

	fmt.Println("Path:", file.Path)
	fmt.Println("Size:", file.Size)
	// Output:
	// Path: /docs/readme.txt
	// Size: 7
}

func Example_contentInterface() {
	// All content types implement Content interface
	var contents []content.Content

	contents = append(contents, content.NewText("text"))
	contents = append(contents, content.NewImage([]byte{0x01}, "image/png"))
	contents = append(contents, content.NewResource("file:///path"))
	contents = append(contents, content.NewAudio([]byte{0x01}, "audio/mpeg"))
	contents = append(contents, content.NewFile([]byte{0x01}, "text/plain"))

	for _, c := range contents {
		fmt.Println(c.Type())
	}
	// Output:
	// text
	// image
	// resource
	// audio
	// file
}

func Example_builderWorkflow() {
	builder := content.NewBuilder()

	// Build various content types
	text := builder.Text("Hello")
	img := builder.ImageWithAlt([]byte{0x89, 0x50}, "image/png", "Logo")
	file := builder.FileWithPath([]byte("data"), "text/csv", "/report.csv")

	fmt.Println("Text type:", text.Type())
	fmt.Println("Image alt:", img.AltText)
	fmt.Println("File path:", file.Path)
	// Output:
	// Text type: text
	// Image alt: Logo
	// File path: /report.csv
}
