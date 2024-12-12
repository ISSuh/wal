package entry

import (
	"bytes"
	"testing"
)

func TestEncodeLogMetadata(t *testing.T) {
	metadata := LogMetadata{
		SegmentID: 1,
		Size:      100,
		Sequence:  10,
		CRC:       12345678,
		Offset:    1234567890,
	}

	expected := []byte{
		0, 0, 0, 1, // SegmentID
		0, 0, 0, 100, // Size
		0, 0, 0, 10, // Sequence
		0, 188, 97, 78, // CRC
		0, 0, 0, 0, 73, 150, 2, 210, // Offset
	}

	result := EncodeLogMetadata(metadata)
	if !bytes.Equal(result, expected) {
		t.Errorf("EncodeLogMetadata() = %v, want %v", result, expected)
	}
}

func TestDecodeLogMetadata(t *testing.T) {
	buf := []byte{
		0, 0, 0, 1, // SegmentID
		0, 0, 0, 100, // Size
		0, 0, 0, 10, // Sequence
		0, 188, 97, 78, // CRC
		0, 0, 0, 0, 73, 150, 2, 210, // Offset
	}

	expected := LogMetadata{
		SegmentID: 1,
		Size:      100,
		Sequence:  10,
		CRC:       12345678,
		Offset:    1234567890,
	}

	result, err := DecodeLogMetadata(buf)
	if err != nil {
		t.Fatalf("DecodeLogMetadata() error = %v", err)
	}
	if result != expected {
		t.Errorf("DecodeLogMetadata() = %v, want %v", result, expected)
	}
}
