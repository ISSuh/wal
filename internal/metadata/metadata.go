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
	"github.com/ISSuh/wal/internal/segment"
)

const (
	metadataHeaderByteSize = 12
)

type Data struct {
	Size     int
	Index    uint64
	Metadata []entry.Metadata
}

func newMetadata(size int, index uint64) Data {
	return Data{
		Size:     size,
		Index:    index,
		Metadata: make([]entry.Metadata, size),
	}
}

func EncodeMetadata(m Data) []byte {
	buf := make([]byte, metadataHeaderByteSize)
	binary.BigEndian.PutUint32(buf, uint32(m.Size))
	binary.BigEndian.PutUint64(buf[8:], uint64(m.Index))

	for _, v := range m.Metadata {
		buf = append(buf, entry.EncodeLogMetadata(v)...)
	}
	return buf
}

func DecodeMetadata(data []byte) (Data, error) {
	if len(data) < metadataHeaderByteSize {
		return Data{}, nil
	}

	size := int(binary.BigEndian.Uint32(data[:4]))
	index := binary.BigEndian.Uint64(data[4:12])

	m := Data{
		Size:     size,
		Index:    index,
		Metadata: make([]entry.Metadata, 0),
	}

	segmentMetadataLen := (size - metadataHeaderByteSize) / segment.MetadataByteLen
	for i := 0; i < segmentMetadataLen; i++ {
		beginOffset := metadataHeaderByteSize + (i * segment.MetadataByteLen)
		endOffset := beginOffset + segment.MetadataByteLen

		encodedLogMetadata := data[beginOffset:endOffset]
		logMetadata, err := entry.DecodeLogMetadata(encodedLogMetadata)
		if err != nil {
			return Data{}, fmt.Errorf("failed to decode metadata. %w", err)
		}

		m.Metadata = append(m.Metadata, logMetadata)
	}

	return m, nil
}
