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

package index

import (
	"encoding/binary"
	"testing"
)

func TestEncodeIndex(t *testing.T) {
	tests := []struct {
		name string
		i    Index
		want []byte
	}{
		{
			name: "valid index",
			i:    Index{Index: 1, MetadataOffset: 2, MetadataSize: 3},
			want: func() []byte {
				buf := make([]byte, indexSize)
				binary.LittleEndian.PutUint64(buf, 1)
				binary.LittleEndian.PutUint64(buf[8:], 2)
				binary.LittleEndian.PutUint32(buf[16:], 3)
				return buf
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeIndex(tt.i); !equal(got, tt.want) {
				t.Errorf("EncodeIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeIndex(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Index
		wantErr bool
	}{
		{
			name: "valid index",
			data: func() []byte {
				buf := make([]byte, indexSize)
				binary.LittleEndian.PutUint64(buf, 1)
				binary.LittleEndian.PutUint64(buf[8:], 2)
				binary.LittleEndian.PutUint32(buf[16:], 3)
				return buf
			}(),
			want:    Index{Index: 1, MetadataOffset: 2, MetadataSize: 3},
			wantErr: false,
		},
		{
			name:    "invalid index size",
			data:    []byte{1, 2, 3},
			want:    Index{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeIndex(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
