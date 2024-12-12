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

package metadata

import (
	"fmt"

	"github.com/ISSuh/wal/internal/file"
)

const (
	metadataFileName = "metadata"
)

type File struct {
	file.File
	basePath string

	offset       int64
	lastMetadata Data
}

func NewFile(basePath string) *File {
	return &File{
		File:     file.NewFile(),
		basePath: basePath,
	}
}

func (f *File) Open() error {
	filePath := fmt.Sprintf("%s/%s", f.basePath, metadataFileName)
	err := f.File.Open(filePath)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Close() error {
	if f.File != nil {
		return f.File.Close()
	}
	return nil
}

func (f *File) Write(metadata Data) (int64, error) {
	itemOffset := f.offset
	buf := EncodeMetadata(metadata)
	if err := f.File.Write(buf); err != nil {
		return 0, fmt.Errorf("failed to write metadata. %w", err)
	}

	if err := f.File.Sync(); err != nil {
		return 0, fmt.Errorf("failed to sync metadata. %w", err)
	}

	f.offset += int64(len(buf))
	f.lastMetadata = metadata
	return itemOffset, nil
}

func (f *File) Read(offset int64, len int) (Data, error) {
	data, err := f.File.ReadAt(offset, len)
	if err != nil {
		return Data{}, fmt.Errorf("failed to read metadata. %w", err)
	}

	metadata, err := DecodeMetadata(data)
	if err != nil {
		return Data{}, fmt.Errorf("failed to decode metadata. %w", err)
	}

	return metadata, nil
}

func (f *File) LastOffset() int64 {
	return f.offset
}

func (f *File) Rollback() error {
	size, err := f.File.Size()
	if err != nil {
		return fmt.Errorf("failed to get file size. %w", err)
	}

	if size < int64(f.lastMetadata.Size) {
		return nil
	}

	targetSize := size - int64(f.lastMetadata.Size)
	if err := f.File.Truncate(targetSize); err != nil {
		return fmt.Errorf("failed to truncate metadata file. %w", err)
	}

	return nil
}
