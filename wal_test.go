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
	"os"
	"testing"
)

func createTempDir(path string) {
	deleteAllFilesOnDir(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
}

func deleteAllFilesOnDir(path string) {
	os.RemoveAll(path)
}

func TestNewStorage(t *testing.T) {
	path := "./tmp"
	createTempDir(path)
	defer deleteAllFilesOnDir(path)

	options := Options{Path: path}
	storage, err := NewStorage(options)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if storage == nil {
		t.Fatalf("expected storage to be non-nil")
	}
}

func TestStorage_Write(t *testing.T) {
	t.Run("Write_1", func(t *testing.T) {
		path := "./tmp"
		createTempDir(path)
		defer deleteAllFilesOnDir(path)

		options := Options{Path: path}
		storage, err := NewStorage(options)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		data := []byte("test data")
		index, err := storage.Write(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if index != 0 {
			t.Errorf("expected index to be 0, got %d", index)
		}
	})

	t.Run("Write_2", func(t *testing.T) {
		path := "./tmp"
		createTempDir(path)
		defer deleteAllFilesOnDir(path)

		options := Options{
			Path:            path,
			SegmentFileSize: 1 * 1024,
		}
		storage, err := NewStorage(options)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		data := make([]byte, 0)
		for i := 0; i < 4000; i++ {
			data = append(data, byte(i))
		}

		index, err := storage.Write(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if index != 0 {
			t.Errorf("expected index to be 0, got %d", index)
		}
	})

	t.Run("Write_3", func(t *testing.T) {
		path := "./tmp"
		createTempDir(path)
		defer deleteAllFilesOnDir(path)

		options := Options{
			Path:            path,
			SegmentFileSize: 1 * 1024,
		}
		storage, err := NewStorage(options)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		data := make([]byte, 0)
		for i := 0; i < 100; i++ {
			data = append(data, byte(i))
		}

		index1, err := storage.Write(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if index1 != 0 {
			t.Errorf("expected index1 to be 1, got %d", index1)
		}

		index2, err := storage.Write(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if index2 != 1 {
			t.Errorf("expected index2 to be 2, got %d", index2)
		}
	})
}

func TestStorage_Read(t *testing.T) {
	t.Run("Read_1", func(t *testing.T) {
		path := "./tmp"
		createTempDir(path)
		defer deleteAllFilesOnDir(path)

		options := Options{Path: path}
		storage, err := NewStorage(options)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		data := []byte("test data")
		i, err := storage.Write(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		readData, err := storage.Read(i)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(readData) != "test data" {
			t.Errorf("expected data to be 'test data', got %s", string(readData))
		}
	})

	t.Run("Read_2", func(t *testing.T) {
		path := "./tmp"
		createTempDir(path)
		defer deleteAllFilesOnDir(path)

		options := Options{
			Path:            path,
			SegmentFileSize: 1 * 1024,
		}
		storage, err := NewStorage(options)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		data1 := []byte("test data1")
		index1, err := storage.Write(data1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if index1 != 0 {
			t.Errorf("expected index1 to be 1, got %d", index1)
		}

		data2 := []byte("test data2")
		index2, err := storage.Write(data2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if index2 != 1 {
			t.Errorf("expected index2 to be 2, got %d", index2)
		}

		readData1, err := storage.Read(index2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(readData1) != "test data2" {
			t.Errorf("expected data to be 'test data2', got %s", string(readData1))
		}

		readData2, err := storage.Read(index1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(readData2) != "test data1" {
			t.Errorf("expected data to be 'test data1', got %s", string(readData2))
		}
	})

	t.Run("Read_3", func(t *testing.T) {
		path := "./tmp"
		createTempDir(path)
		defer deleteAllFilesOnDir(path)

		options := Options{
			Path:            path,
			SegmentFileSize: 10,
		}
		storage, err := NewStorage(options)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		data1 := []byte("aaaaaaaaaabbbb")
		index1, err := storage.Write(data1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if index1 != 0 {
			t.Errorf("expected index1 to be 1, got %d", index1)
		}

		readData1, err := storage.Read(index1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(readData1) != "aaaaaaaaaabbbb" {
			t.Errorf("expected data to be 'aaaaaaaaaabbbb', got %s", string(readData1))
		}
	})
}

func TestStorage_Close(t *testing.T) {
	path := "./tmp"
	createTempDir(path)
	defer deleteAllFilesOnDir(path)

	options := Options{Path: path}
	storage, err := NewStorage(options)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = storage.Close()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func BenchmarkWrite(b *testing.B) {
	path := "./tmp"
	createTempDir(path)
	defer deleteAllFilesOnDir(path)

	options := Options{Path: path}
	storage, err := NewStorage(options)
	if err != nil {
		b.Fatalf("expected no error, got %v", err)
	}

	data := []byte("1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := storage.Write(data)
		if err != nil {
			b.Fatalf("expected no error, got %v", err)
		}
	}
}

func BenchmarkWriteWithSyncAfterWrite(b *testing.B) {
	path := "./tmp"
	createTempDir(path)
	defer deleteAllFilesOnDir(path)

	options := Options{Path: path, SyncAfterWrite: true}
	storage, err := NewStorage(options)
	if err != nil {
		b.Fatalf("expected no error, got %v", err)
	}

	data := []byte("1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := storage.Write(data)
		if err != nil {
			b.Fatalf("expected no error, got %v", err)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	path := "./tmp"
	createTempDir(path)
	defer deleteAllFilesOnDir(path)

	options := Options{Path: path}
	storage, err := NewStorage(options)
	if err != nil {
		b.Fatalf("expected no error, got %v", err)
	}

	data := []byte("1")
	index, err := storage.Write(data)
	if err != nil {
		b.Fatalf("expected no error, got %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := storage.Read(index)
		if err != nil {
			b.Fatalf("expected no error, got %v", err)
		}
	}
}
