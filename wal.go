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
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/ISSuh/wal/internal/entry"
	"github.com/ISSuh/wal/internal/index"
	"github.com/ISSuh/wal/internal/metadata"
	"github.com/ISSuh/wal/internal/segment"
)

type rollbackkey int

const (
	rollbackIndex rollbackkey = iota
	rollbackMetadata
	rollbackSegment
)

type Storage interface {
	Write(data []byte) (int64, error)
	Read(index int64) ([]byte, error)
	Close() error
}

type storage struct {
	options Options

	segment      *segment.Segment
	indexFile    *index.File
	metadataFile *metadata.File

	segmentIDCounter int
	mutex            sync.RWMutex
}

func NewStorage(option Options) (Storage, error) {
	indexFile := index.NewFile(option.Path)
	if err := indexFile.Open(); err != nil {
		return nil, fmt.Errorf("failed to open index file. %w", err)
	}

	metadataFile := metadata.NewFile(option.Path)
	if err := metadataFile.Open(); err != nil {
		return nil, fmt.Errorf("failed to open metadata file. %w", err)
	}

	segment, err := segment.NewSegment(0, option.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to create segment. %w", err)
	}

	return &storage{
		options:          option,
		segment:          segment,
		indexFile:        indexFile,
		metadataFile:     metadataFile,
		segmentIDCounter: 0,
	}, nil
}

func (s *storage) Write(data []byte) (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newIndexSeq := s.indexFile.LastIndex() + 1

	// append data to segment
	logMetadata, err := s.appendLogToSegment(newIndexSeq, data)
	if err != nil {
		return 0, fmt.Errorf("failed to append data to segment. %w", err)
	}

	// append metadata to metadata file
	metadata := metadata.NewMetadata(newIndexSeq, logMetadata)
	metadataOffset, err := s.metadataFile.Write(metadata)
	if err != nil {
		if rollbackErr := s.rollback(rollbackIndex); rollbackErr != nil {
			errs := errors.Join(rollbackErr, fmt.Errorf("failed to write metadata. %w", err))
			return 0, fmt.Errorf("failed to rollback index. %w ", errs)
		}
		return 0, fmt.Errorf("failed to write metadata. %w", err)
	}

	// append index to index file
	index := index.NewIndex(newIndexSeq, metadataOffset, metadata.Size)
	if err := s.indexFile.Write(index); err != nil {
		if rollbackErr := s.rollback(rollbackIndex); rollbackErr != nil {
			errs := errors.Join(rollbackErr, fmt.Errorf("failed to write metadata. %w", err))
			return 0, fmt.Errorf("failed to rollback index. %w. %w", rollbackErr, errs)
		}

		if rollbackErr := s.rollback(rollbackMetadata); rollbackErr != nil {
			errs := errors.Join(rollbackErr, fmt.Errorf("failed to write metadata. %w", err))
			return 0, fmt.Errorf("failed to rollback metadata. %w. %w", rollbackErr, errs)
		}

		return 0, fmt.Errorf("failed to write index. %w", err)
	}

	return newIndexSeq, nil
}

func (s *storage) Read(i int64) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// read index from index file
	index, err := s.indexFile.Read(i)
	if err != nil {
		return nil, fmt.Errorf("failed to read index. %w", err)
	}

	// read metadata from metadata file
	metadata, err := s.readMetadata(index)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata. %w", err)
	}

	// read data from segment
	data, err := s.readLogFromSegment(metadata.LogMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from segment. %w", err)
	}

	return data, nil
}

// closes storage
func (s *storage) Close() error {
	if err := s.segment.Close(); err != nil {
		return fmt.Errorf("failed to close segment. %w", err)
	}

	if err := s.indexFile.Close(); err != nil {
		return fmt.Errorf("failed to close index file. %w", err)
	}

	if err := s.metadataFile.Close(); err != nil {
		return fmt.Errorf("failed to close metadata file. %w", err)
	}

	return nil
}

// calculateOffsetFromData calculates offset of data and need new segment after append
func (s *storage) calculateOffsetFromData(len int) (int, bool) {
	offset := 0
	needNewSegmentAfterAppend := false
	switch {
	case len > s.options.SegmentFileSize:
		offset = s.options.SegmentFileSize
		needNewSegmentAfterAppend = true
	case len+s.segment.Size() > s.options.SegmentFileSize:
		offset = s.options.SegmentFileSize - s.segment.Size()
		needNewSegmentAfterAppend = true
	default:
		offset = len
	}

	return offset, needNewSegmentAfterAppend
}

// appendLogToSegment appends log to segment
func (s *storage) appendLogToSegment(newIndex int64, data []byte) ([]entry.LogMetadata, error) {
	offset := 0
	dataSize := len(data)
	remainedDataSize := dataSize
	sequence := 0
	segmentMetadata := make([]entry.LogMetadata, 0)
	for offset <= dataSize {
		// calculate offset of data
		needNewSegmentAfterAppend := false
		offset, needNewSegmentAfterAppend = s.calculateOffsetFromData(remainedDataSize)

		// create log
		partaialData := data[0:offset]
		log := entry.NewLog(newIndex, sequence, partaialData)

		// append log to segment
		m, err := s.segment.Append(log)
		if err != nil {
			return nil, fmt.Errorf("failed to append log to segment. %w", err)
		}

		if needNewSegmentAfterAppend {
			s.segmentIDCounter++

			// create new segment
			segment, err := segment.NewSegment(s.segmentIDCounter, s.options.Path)
			if err != nil {
				return nil, fmt.Errorf("failed to create new segment. %w", err)
			}

			// switch segment
			s.segment = segment
		}

		segmentMetadata = append(segmentMetadata, m)
		remainedDataSize = len(data[offset:])
		sequence++
	}

	return segmentMetadata, nil
}

// readMetadata reads metadata from metadata file
func (s *storage) readMetadata(index index.Index) (metadata.Data, error) {
	offset := index.MetadataOffset
	size := index.MetadataSize
	data, err := s.metadataFile.Read(offset, size)
	if err != nil {
		return metadata.Data{}, fmt.Errorf("failed to read metadata. %w", err)
	}

	// sort metadata by sequence
	sort.Slice(data.LogMetadata, func(i, j int) bool {
		return data.LogMetadata[i].Sequence < data.LogMetadata[j].Sequence
	})

	return data, nil
}

// readLogFromSegment reads log from segment
func (s *storage) readLogFromSegment(logMetadata []entry.LogMetadata) ([]byte, error) {
	data := make([]byte, 0)
	for _, m := range logMetadata {
		// open segment if segment id is different
		seg := s.segment
		if m.SegmentID != s.segment.ID() {
			var err error
			seg, err = segment.NewSegment(m.SegmentID, s.options.Path)
			if err != nil {
				return nil, fmt.Errorf("failed to open segment. %w", err)
			}
		}

		log, err := seg.Read(m.Offset, m.Size)
		if err != nil {
			return nil, fmt.Errorf("failed to read log. %w", err)
		}

		data = append(data, log.PayLoad...)
	}

	return data, nil
}

func (s *storage) rollback(key rollbackkey) error {
	switch key {
	case rollbackIndex:
		return s.indexFile.Rollback()
	case rollbackMetadata:
		return s.metadataFile.Rollback()
	default:
		return nil
	}
}
