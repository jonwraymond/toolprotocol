package content

// Type identifies the type of content.
type Type string

const (
	// TypeText is plain text content.
	TypeText Type = "text"

	// TypeImage is image data.
	TypeImage Type = "image"

	// TypeResource is an external resource reference.
	TypeResource Type = "resource"

	// TypeAudio is audio data.
	TypeAudio Type = "audio"

	// TypeFile is file data.
	TypeFile Type = "file"
)

// Content is the interface for all content types.
//
// Contract:
//   - Concurrency: Implementations are safe for concurrent reads.
//   - Immutability: Content should not be modified after creation.
//   - Errors: Bytes returns nil error for in-memory content types.
//   - Encoding: String returns base64 for binary types (image, audio).
type Content interface {
	// Type returns the content type.
	Type() Type

	// MIMEType returns the MIME type of the content.
	MIMEType() string

	// Bytes returns the raw content bytes.
	Bytes() ([]byte, error)

	// String returns a string representation of the content.
	String() string
}
