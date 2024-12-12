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
	"fmt"

	"github.com/ISSuh/wal/internal/file"
)

const (
	IndexFileName = "index"
)

type File struct {
	file.File
	basePath string

	lastIndex Index
	size      int
}

func NewFile(basePath string) *File {
	return &File{
		File:     file.NewFile(),
		basePath: basePath,
	}
}

func (f *File) Open() error {
	filePath := fmt.Sprintf("%s/%s", f.basePath, IndexFileName)
	if err := f.File.Open(filePath); err != nil {
		return fmt.Errorf("failed to open index file. %w", err)
	}
	return nil
}

func (f *File) Close() error {
	if f.File != nil {
		if err := f.File.Close(); err != nil {
			return fmt.Errorf("failed to close index file. %w", err)
		}
	}
	return nil
}

func (f *File) Write(i Index) error {
	buf := EncodeIndex(i)
	if err := f.File.Write(buf); err != nil {
		return fmt.Errorf("failed to write index. %w", err)
	}

	if err := f.File.Sync(); err != nil {
		return fmt.Errorf("failed to sync index file. %w", err)
	}

	f.lastIndex = i
	f.size = len(buf)
	return nil
}

func (f *File) Read(i int64) (Index, error) {
	offset := i * indexSize
	buf, err := f.File.ReadAt(offset, indexSize)
	if err != nil {
		return Index{}, fmt.Errorf("failed to read index. %w", err)
	}

	index, err := DecodeIndex(buf)
	if err != nil {
		return index, fmt.Errorf("failed to decode index. %w", err)
	}

	return index, nil
}

func (f *File) LastIndex() int64 {
	return f.lastIndex.Index
}

func (f *File) Rollback() error {
	if f.size < indexSize {
		return nil
	}

	targetSize := int64(f.size - indexSize)
	if err := f.File.Truncate(targetSize); err != nil {
		return fmt.Errorf("failed to truncate index file. %w", err)
	}

	return nil
}
