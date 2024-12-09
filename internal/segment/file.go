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

	"github.com/ISSuh/wal/internal/entry"
	"github.com/ISSuh/wal/internal/file"
)

const (
	SegmentFilePrefix = "segment"
)

type File struct {
	file.File
	basePath string
}

func NewFile(basePath string) *File {
	return &File{
		File:     file.NewFile(),
		basePath: basePath,
	}
}

func (f *File) Open(id int) error {
	filewithPath := fmt.Sprintf("%s/%s_%d", f.basePath, SegmentFilePrefix, id)
	if err := f.File.Open(filewithPath); err != nil {
		return fmt.Errorf("failed to open segment file. %w", err)
	}
	return nil
}

func (f *File) Close() error {
	if f.File != nil {
		if err := f.File.Close(); err != nil {
			return fmt.Errorf("failed to close segment file. %w", err)
		}
	}
	return nil
}

func (f *File) Write(log entry.Log) error {
	data := entry.EncodeLog(log)
	if err := f.File.Write(data); err != nil {
		return fmt.Errorf("failed to write segment file. %w", err)
	}

	if err := f.File.Sync(); err != nil {
		return fmt.Errorf("failed to write segment file. %w", err)
	}

	return nil
}

func (f *File) Read(offset int64, len int) (entry.Log, error) {
	data, err := f.File.ReadAt(offset, len)
	if err != nil {
		return entry.Log{}, fmt.Errorf("failed to read segment file. %w", err)
	}

	log, err := entry.DecodeLog(data)
	if err != nil {
		return entry.Log{}, fmt.Errorf("failed to decode segment file. %w", err)
	}

	return log, nil
}
