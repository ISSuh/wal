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

package file

import (
	"fmt"
	"os"
)

type File interface {
	Open(filePath string) error
	Close() error
	Write([]byte) error
	ReadAt(int64, int) ([]byte, error)
	Sync() error
	Size() (int64, error)
	Truncate(int64) error
	Path() string
}

type file struct {
	filePath string
	f        *os.File
}

func NewFile() File {
	return &file{}
}

func (f *file) Open(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	f.filePath = filePath
	f.f = file
	return nil
}

func (f *file) Close() error {
	err := f.f.Close()
	if err != nil {
		return err
	}
	return nil
}

func (f *file) Write(data []byte) error {
	n, err := f.f.Write(data)
	if err != nil {
		return err
	}

	if n != len(data) {
		return fmt.Errorf("failed to write all data. %d != %d", n, len(data))
	}

	return nil
}

func (f *file) ReadAt(offset int64, size int) ([]byte, error) {
	buf := make([]byte, size)
	_, err := f.f.ReadAt(buf, offset)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (f *file) Sync() error {
	err := f.f.Sync()
	if err != nil {
		return err
	}
	return nil
}

func (f *file) Size() (int64, error) {
	stat, err := f.f.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func (f *file) Truncate(size int64) error {
	err := f.f.Truncate(size)
	if err != nil {
		return err
	}
	return nil
}

func (f *file) Path() string {
	return f.filePath
}
