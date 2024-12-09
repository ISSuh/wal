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
	"os"
	"testing"
)

func TestFile_Open(t *testing.T) {
	f := NewFile()
	err := f.Open("testfile.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer os.Remove("testfile.txt")
	defer f.Close()

	if f.Path() != "testfile.txt" {
		t.Errorf("expected file path to be 'testfile.txt', got %s", f.Path())
	}
}

func TestFile_Write(t *testing.T) {
	f := NewFile()
	err := f.Open("testfile.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer os.Remove("testfile.txt")
	defer f.Close()

	data := []byte("hello world")
	err = f.Write(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFile_ReadAt(t *testing.T) {
	f := NewFile()
	err := f.Open("testfile.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer os.Remove("testfile.txt")
	defer f.Close()

	data := []byte("hello world")
	err = f.Write(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	readData, err := f.ReadAt(0, len(data))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(readData) != string(data) {
		t.Errorf("expected %s, got %s", string(data), string(readData))
	}
}

func TestFile_Sync(t *testing.T) {
	f := NewFile()
	err := f.Open("testfile.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer os.Remove("testfile.txt")
	defer f.Close()

	err = f.Sync()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFile_Size(t *testing.T) {
	f := NewFile()
	err := f.Open("testfile.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer os.Remove("testfile.txt")
	defer f.Close()

	data := []byte("hello world")
	err = f.Write(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	size, err := f.Size()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if size != int64(len(data)) {
		t.Errorf("expected size %d, got %d", len(data), size)
	}
}

func TestFile_Truncate(t *testing.T) {
	f := NewFile()
	err := f.Open("testfile.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer os.Remove("testfile.txt")
	defer f.Close()

	data := []byte("hello world")
	err = f.Write(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = f.Truncate(5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	size, err := f.Size()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if size != 5 {
		t.Errorf("expected size 5, got %d", size)
	}
}
