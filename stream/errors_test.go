package stream

import (
	"errors"
	"testing"
)

func TestErrStreamClosed(t *testing.T) {
	if ErrStreamClosed.Error() != "stream closed" {
		t.Errorf("Error() = %q", ErrStreamClosed.Error())
	}
}

func TestErrBufferFull(t *testing.T) {
	if ErrBufferFull.Error() != "stream buffer full" {
		t.Errorf("Error() = %q", ErrBufferFull.Error())
	}
}

func TestStreamError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *StreamError
		want string
	}{
		{
			name: "with underlying error",
			err: &StreamError{
				StreamID: "stream-123",
				Op:       "send",
				Err:      ErrStreamClosed,
			},
			want: "stream stream-123: send: stream closed",
		},
		{
			name: "without underlying error",
			err: &StreamError{
				StreamID: "stream-123",
				Op:       "close",
				Err:      nil,
			},
			want: "stream stream-123: close",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStreamError_Unwrap(t *testing.T) {
	underlying := ErrStreamClosed
	err := &StreamError{
		StreamID: "stream-123",
		Op:       "send",
		Err:      underlying,
	}

	if !errors.Is(err, underlying) {
		t.Error("errors.Is should return true for underlying error")
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != underlying {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlying)
	}
}
