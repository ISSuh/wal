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

package entry

import (
	"encoding/binary"
	"fmt"
)

const (
	MetadataByteLen = 20
)

type Metadata struct {
	Size   int
	Index  int
	Offset int64
	CRC    uint32
}

func EncodeLogMetadata(m Metadata) []byte {
	buf := make([]byte, MetadataByteLen)
	binary.BigEndian.PutUint64(buf, uint64(m.Size))
	binary.BigEndian.PutUint32(buf[8:], uint32(m.Index))
	binary.BigEndian.PutUint64(buf[16:], uint64(m.Offset))
	binary.BigEndian.PutUint32(buf[24:], m.CRC)
	return buf
}

func DecodeLogMetadata(buf []byte) (Metadata, error) {
	if len(buf) != MetadataByteLen {
		return Metadata{}, fmt.Errorf("invalid segment metadata size. %d", len(buf))
	}

	size := int(binary.BigEndian.Uint64(buf[:8]))
	index := int(binary.BigEndian.Uint64(buf[8:16]))
	offset := int64(binary.BigEndian.Uint64(buf[16:24]))
	crc := binary.BigEndian.Uint32(buf[24:])
	return Metadata{
		Size:   size,
		Index:  index,
		Offset: offset,
		CRC:    crc,
	}, nil
}
