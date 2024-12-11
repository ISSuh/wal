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

package segment

import "github.com/ISSuh/wal/internal/entry"

const (
	ByteSizeOfUint64 = 8
)

type Segment struct {
	file       *File
	id         uint64
	size       int
	firstIndex uint64
	offset     int64
	buffer     []entry.Log
}

func NewSegment(id uint64) *Segment {
	return &Segment{
		id:         id,
		size:       0,
		firstIndex: 0,
		offset:     0,
		buffer:     make([]entry.Log, 0),
	}
}

func (s *Segment) Append(e entry.Log) (entry.Metadata, error) {
	return entry.Metadata{}, nil
}

func (s *Segment) Size() int {
	return s.size
}

func (s *Segment) FirstIndex() uint64 {
	return s.firstIndex
}

func (s *Segment) Offset() int64 {
	return s.offset
}
