/*
MIT License

Copyright (c) 2024 ISSuh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package metadata

import (
	"testing"

	"github.com/ISSuh/wal/internal/entry"
)

func TestNewMetadata(t *testing.T) {
	size := 12 + (24 * 2)
	index := int64(1)
	logMetadata := []entry.LogMetadata{
		{
			SegmentID: 1,
			Size:      1,
			Sequence:  1,
			CRC:       1,
			Offset:    0,
		},
		{
			SegmentID: 1,
			Size:      2,
			Sequence:  2,
			CRC:       2,
			Offset:    233,
		},
	}

	m := NewMetadata(1, logMetadata)

	if m.Size != size {
		t.Errorf("expected size %d, got %d", size, m.Size)
	}
	if m.Index != index {
		t.Errorf("expected index %d, got %d", index, m.Index)
	}
	if len(m.LogMetadata) != len(logMetadata) {
		t.Errorf("expected log metadata size %d, got %d", len(logMetadata), len(m.LogMetadata))
	}
}

func TestEncodeMetadata(t *testing.T) {
	logMetadata := []entry.LogMetadata{
		{
			SegmentID: 1,
			Size:      1,
			Sequence:  1,
			CRC:       1,
			Offset:    0,
		},
		{
			SegmentID: 1,
			Size:      2,
			Sequence:  2,
			CRC:       2,
			Offset:    233,
		},
	}
	m := NewMetadata(1, logMetadata)

	encoded := EncodeMetadata(m)
	expectedSize := metadataHeaderByteSize + len(entry.EncodeLogMetadata(logMetadata[0]))*2

	if len(encoded) != expectedSize {
		t.Errorf("expected encoded size %d, got %d", expectedSize, len(encoded))
	}
}

func TestDecodeMetadata(t *testing.T) {
	logMetadata := []entry.LogMetadata{
		{
			SegmentID: 1,
			Size:      1,
			Sequence:  1,
			CRC:       1,
			Offset:    0,
		},
		{
			SegmentID: 1,
			Size:      2,
			Sequence:  2,
			CRC:       2,
			Offset:    233,
		},
	}
	m := NewMetadata(1, logMetadata)
	encoded := EncodeMetadata(m)

	decoded, err := DecodeMetadata(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if decoded.Size != m.Size {
		t.Errorf("expected size %d, got %d", m.Size, decoded.Size)
	}
	if decoded.Index != m.Index {
		t.Errorf("expected index %d, got %d", m.Index, decoded.Index)
	}
	if len(decoded.LogMetadata) != len(m.LogMetadata) {
		t.Errorf("expected log metadata size %d, got %d", len(m.LogMetadata), len(decoded.LogMetadata))
	}
	for i, v := range decoded.LogMetadata {
		if v != m.LogMetadata[i] {
			t.Errorf("expected log metadata %v, got %v", m.LogMetadata[i], v)
		}
	}
}
