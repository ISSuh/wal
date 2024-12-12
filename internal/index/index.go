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
	"fmt"
)

const (
	indexSize = 20
)

type Index struct {
	Index          int64
	MetadataOffset int64
	MetadataSize   int
}

func NewIndex(index int64, metadataOffset int64, metadataSize int) Index {
	return Index{
		Index:          index,
		MetadataOffset: metadataOffset,
		MetadataSize:   metadataSize,
	}
}

func EncodeIndex(i Index) []byte {
	buf := make([]byte, indexSize)
	binary.LittleEndian.PutUint64(buf, uint64(i.Index))
	binary.LittleEndian.PutUint64(buf[8:], uint64(i.MetadataOffset))
	binary.LittleEndian.PutUint32(buf[16:], uint32(i.MetadataSize))
	return buf
}

func DecodeIndex(data []byte) (Index, error) {
	if len(data) != indexSize {
		return Index{}, fmt.Errorf("invalid index size. %d", len(data))
	}

	i := Index{}
	i.Index = int64(binary.LittleEndian.Uint64(data[:8]))
	i.MetadataOffset = int64(binary.LittleEndian.Uint64(data[8:16]))
	i.MetadataSize = int(binary.LittleEndian.Uint32(data[16:20]))
	return i, nil
}
