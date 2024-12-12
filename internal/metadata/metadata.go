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
	"encoding/binary"
	"fmt"

	"github.com/ISSuh/wal/internal/entry"
)

const (
	metadataHeaderByteSize = 12
)

type Data struct {
	Size        int
	Index       int64
	LogMetadata []entry.LogMetadata
}

func NewMetadata(index int64, m []entry.LogMetadata) Data {
	return Data{
		Size:        metadataHeaderByteSize + (len(m) * entry.MetadataByteLen),
		Index:       index,
		LogMetadata: m,
	}
}

func EncodeMetadata(m Data) []byte {
	buf := make([]byte, metadataHeaderByteSize)
	binary.BigEndian.PutUint32(buf[0:4], uint32(m.Size))
	binary.BigEndian.PutUint64(buf[4:12], uint64(m.Index))

	for _, v := range m.LogMetadata {
		buf = append(buf, entry.EncodeLogMetadata(v)...)
	}
	return buf
}

func DecodeMetadata(data []byte) (Data, error) {
	if len(data) < metadataHeaderByteSize {
		return Data{}, nil
	}

	size := int(binary.BigEndian.Uint32(data[:4]))
	index := int64(binary.BigEndian.Uint64(data[4:12]))

	m := Data{
		Size:        size,
		Index:       index,
		LogMetadata: make([]entry.LogMetadata, 0),
	}

	segmentMetadataLen := (size - metadataHeaderByteSize) / entry.MetadataByteLen
	for i := 0; i < segmentMetadataLen; i++ {
		beginOffset := metadataHeaderByteSize + (i * entry.MetadataByteLen)
		endOffset := beginOffset + entry.MetadataByteLen

		encodedLogMetadata := data[beginOffset:endOffset]
		logMetadata, err := entry.DecodeLogMetadata(encodedLogMetadata)
		if err != nil {
			return Data{}, fmt.Errorf("failed to decode metadata. %w", err)
		}

		m.LogMetadata = append(m.LogMetadata, logMetadata)
	}

	return m, nil
}
