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

import (
	"fmt"

	"github.com/ISSuh/wal/internal/crc"
	"github.com/ISSuh/wal/internal/entry"
	"github.com/ISSuh/wal/internal/file"
)

const (
	segmentFilePrefix = "segment"
)

type Segment struct {
	id        int
	size      int
	offset    int64
	lastIndex int64

	file     file.File
	basePath string
}

func NewSegment(id int, basePath string) (*Segment, error) {
	s := &Segment{
		id:        id,
		size:      0,
		offset:    0,
		lastIndex: 0,
		file:      file.NewFile(),
		basePath:  basePath,
	}

	if err := s.open(id); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Segment) Append(e entry.Log) (entry.LogMetadata, error) {
	if err := s.write(e); err != nil {
		return entry.LogMetadata{}, err
	}

	crc := crc.Encode(e.PayLoad)
	m := entry.LogMetadata{
		SegmentID: s.id,
		Size:      len(e.PayLoad),
		Sequence:  e.Sequence,
		Offset:    s.offset,
		CRC:       crc,
	}

	s.offset += int64(m.Size)
	s.size += m.Size
	return m, nil
}

func (s *Segment) Read(offset int64, len int) (entry.Log, error) {
	data, err := s.file.ReadAt(offset, len)
	if err != nil {
		return entry.Log{}, fmt.Errorf("failed to read segment file. %w", err)
	}

	log, err := entry.DecodeLog(data)
	if err != nil {
		return entry.Log{}, fmt.Errorf("failed to decode segment file. %w", err)
	}

	return log, nil
}

func (s *Segment) ID() int {
	return s.id
}

func (s *Segment) Size() int {
	return s.size
}

func (s *Segment) Offset() int64 {
	return s.offset
}

func (s *Segment) Close() error {
	if s.file != nil {
		if err := s.file.Close(); err != nil {
			return fmt.Errorf("failed to close segment file. %w", err)
		}
	}
	return nil
}

func (s *Segment) Sync() error {
	if err := s.file.Sync(); err != nil {
		return fmt.Errorf("failed to sync segment file. %w", err)
	}
	return nil
}

func (s *Segment) open(id int) error {
	filewithPath := fmt.Sprintf("%s/%s_%d", s.basePath, segmentFilePrefix, id)
	if err := s.file.Open(filewithPath); err != nil {
		return fmt.Errorf("failed to open segment file. %w", err)
	}
	return nil
}

func (s *Segment) write(log entry.Log) error {
	data := entry.EncodeLog(log)
	if err := s.file.Write(data); err != nil {
		return fmt.Errorf("failed to write segment file. %w", err)
	}

	if err := s.file.Sync(); err != nil {
		return fmt.Errorf("failed to write segment file. %w", err)
	}

	return nil
}
