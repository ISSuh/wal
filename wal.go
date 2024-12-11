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

package wal

import (
	"github.com/ISSuh/wal/internal/entry"
	"github.com/ISSuh/wal/internal/index"
	"github.com/ISSuh/wal/internal/metadata"
	"github.com/ISSuh/wal/internal/segment"
)

type Storage interface {
	Write(data []byte) (uint64, error)
	Read(index uint64) ([]byte, error)
	Close() error
	Commit() error
}

type storage struct {
	options Options

	segment      *segment.Segment
	indexFile    *index.File
	metadataFile *metadata.File

	segmentIDCounter uint64
}

func NewStorage(option Options) (Storage, error) {
	return &storage{
		options:          option,
		segmentIDCounter: 0,
	}, nil
}

func (s *storage) Write(data []byte) (uint64, error) {
	newIndexSeq := s.indexFile.LastIndex() + 1

	log := entry.NewLog(len(data), 0, data)
	logMetadata, err := s.segment.Append(log)
	if err != nil {
		return 0, err
	}

	metadata := metadata.Data{
		Size:     4 + (1 * entry.MetadataByteLen),
		Index:    newIndexSeq,
		Metadata: []entry.Metadata{logMetadata},
	}

	if err := s.metadataFile.Write(metadata); err != nil {
		return 0, err
	}

	index := index.Index{
		Index:          newIndexSeq,
		MetadataOffset: 0,
	}

	if err := s.indexFile.Write(index); err != nil {
		return 0, err
	}

	return newIndexSeq, nil
}

func (s *storage) Read(index uint64) ([]byte, error) {
	return nil, nil
}

func (s *storage) Close() error {
	return nil
}

func (s *storage) Commit() error {
	return nil
}

func (s *storage) newSegment() error {
	s.segmentIDCounter++

	s.segment = segment.NewSegment(s.segmentIDCounter)
	return nil
}
