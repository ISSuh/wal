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
	MetadataByteLen = 24
)

type LogMetadata struct {
	SegmentID int
	Size      int
	Sequence  int
	CRC       uint32
	Offset    int64
}

func EncodeLogMetadata(m LogMetadata) []byte {
	buf := make([]byte, MetadataByteLen)
	binary.BigEndian.PutUint32(buf[:4], uint32(m.SegmentID))
	binary.BigEndian.PutUint32(buf[4:8], uint32(m.Size))
	binary.BigEndian.PutUint32(buf[8:12], uint32(m.Sequence))
	binary.BigEndian.PutUint32(buf[12:16], m.CRC)
	binary.BigEndian.PutUint64(buf[16:], uint64(m.Offset))
	return buf
}

func DecodeLogMetadata(buf []byte) (LogMetadata, error) {
	if len(buf) != MetadataByteLen {
		return LogMetadata{}, fmt.Errorf("invalid segment metadata size. %d", len(buf))
	}

	segmentID := int(binary.BigEndian.Uint32(buf[:4]))
	size := int(binary.BigEndian.Uint32(buf[4:8]))
	sequence := int(binary.BigEndian.Uint32(buf[8:12]))
	crc := binary.BigEndian.Uint32(buf[12:16])
	offset := int64(binary.BigEndian.Uint64(buf[16:]))
	return LogMetadata{
		SegmentID: segmentID,
		Size:      size,
		Sequence:  sequence,
		Offset:    offset,
		CRC:       crc,
	}, nil
}
