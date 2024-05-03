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
	"encoding/binary"
	"fmt"
	"os"
)

const (
	SegmentMetadataByteLen = 12
	ByteSizeOfUint64       = 8
)

type segmentHeader struct {
	metadataOffset int
	metadataLen    int
}

type segmentMeta struct {
	index  uint64
	offset uint32
}

func (m *segmentMeta) encode() []byte {
	buf := make([]byte, SegmentMetadataByteLen)
	binary.BigEndian.PutUint64(buf, m.index)
	binary.BigEndian.PutUint32(buf, uint32(m.offset))
	return buf
}

func decodeSegmentMeta(buf []byte) segmentMeta {
	index := binary.BigEndian.Uint64(buf[:ByteSizeOfUint64])
	offset := binary.BigEndian.Uint32(buf[ByteSizeOfUint64:])
	return segmentMeta{
		index:  index,
		offset: offset,
	}
}

func loadSegmentMetadatas(buf []byte) ([]segmentMeta, error) {
	if (len(buf) % SegmentMetadataByteLen) != 0 {
		return nil, fmt.Errorf("invalid metadata buffer size")
	}

	metas := []segmentMeta{}
	offset := 0
	size := len(buf)
	for offset < size {
		b := buf[offset:SegmentMetadataByteLen]
		m := decodeSegmentMeta(b)
		metas = append(metas, m)

		offset += SegmentMetadataByteLen
	}
	return metas, nil
}

type Segment struct {
	id          uint64
	logFile     *os.File
	logMetaFile *os.File
	segmentSize int
	offset      uint32
	buffer      []byte
	metadatas   []segmentMeta
}

func NewSegment(basePath string, id uint64, segmentSize int) (*Segment, error) {
	logFile, err := openLogFile(basePath, id)
	if err != nil {
		return nil, err
	}

	logMetaFile, err := openLogMetaFile(basePath, id)
	if err != nil {
		return nil, err
	}

	return &Segment{
		id:          id,
		segmentSize: segmentSize,
		offset:      0,
		logFile:     logFile,
		logMetaFile: logMetaFile,
		buffer:      []byte{},
		metadatas:   []segmentMeta{},
	}, nil
}

func (s *Segment) Write(index uint64, buf []byte) error {
	n, err := s.logFile.Write(buf)
	if err != nil {
		return err
	}

	s.offset += uint32(n)
	meta := segmentMeta{
		index:  index,
		offset: s.offset,
	}

	s.metadatas = append(s.metadatas, meta)
	_, err = s.logMetaFile.Write(meta.encode())
	if err != nil {
		return err
	}
	return nil
}

func (s *Segment) WriteWithSync(index uint64, buf []byte) error {
	if err := s.Write(index, buf); err != nil {
		return nil
	}

	if err := s.logFile.Sync(); err != nil {
		return err
	}
	return nil
}

func (s *Segment) Read(index uint64) ([]byte, error) {
	buf := []byte{}
	_, err := s.logFile.Read(buf)
	return buf, err
}

func (s *Segment) Load() {
}

func (s *Segment) Commit() error {
	return s.logFile.Sync()
}

func (s *Segment) Rollback() {
}

func (s *Segment) Close() error {
	if err := s.logFile.Sync(); err != nil {
		return err
	}

	if err := s.logMetaFile.Sync(); err != nil {
		return err
	}

	if err := s.logFile.Close(); err != nil {
		return nil
	}

	if err := s.logMetaFile.Close(); err != nil {
		return nil
	}

	return nil
}
