package index

import (
	"os"
	"testing"
)

func setup() (*File, func()) {
	basePath := "./testdata"
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		panic(err)
	}

	f := NewFile(basePath)
	return f, func() {
		os.RemoveAll(basePath)
	}
}

func TestNewFile(t *testing.T) {
	basePath := "./testdata"
	f := NewFile(basePath)
	if f == nil {
		t.Errorf("NewFile() returned nil")
	}

	if f.basePath != basePath {
		t.Errorf("NewFile() basePath = %v, want %v", f.basePath, basePath)
	}
}

func TestFile_Open(t *testing.T) {
	f, teardown := setup()
	defer teardown()

	if err := f.Open(); err != nil {
		t.Errorf("File.Open() error = %v", err)
	}
}

func TestFile_Close(t *testing.T) {
	f, teardown := setup()
	defer teardown()

	if err := f.Open(); err != nil {
		t.Errorf("File.Open() error = %v", err)
	}

	if err := f.Close(); err != nil {
		t.Errorf("File.Close() error = %v", err)
	}
}

func TestFile_Write(t *testing.T) {
	f, teardown := setup()
	defer teardown()

	if err := f.Open(); err != nil {
		t.Errorf("File.Open() error = %v", err)
	}
	defer f.Close()

	index := Index{}
	if err := f.Write(index); err != nil {
		t.Errorf("File.Write() error = %v", err)
	}
	if f.LastIndex() != 1 {
		t.Errorf("File.LastIndex() = %v, want %v", f.LastIndex(), 1)
	}
}

func TestFile_Read(t *testing.T) {
	f, teardown := setup()
	defer teardown()

	if err := f.Open(); err != nil {
		t.Errorf("File.Open() error = %v", err)
	}
	defer f.Close()

	index := Index{}
	if err := f.Write(index); err != nil {
		t.Errorf("File.Write() error = %v", err)
	}

	readIndex, err := f.Read(0)
	if err != nil {
		t.Errorf("File.Read() error = %v", err)
	}
	if readIndex != index {
		t.Errorf("File.Read() = %v, want %v", readIndex, index)
	}
}

func TestFile_LastIndex(t *testing.T) {
	f, teardown := setup()
	defer teardown()

	if err := f.Open(); err != nil {
		t.Errorf("File.Open() error = %v", err)
	}
	defer f.Close()

	if f.LastIndex() != 0 {
		t.Errorf("File.LastIndex() = %v, want %v", f.LastIndex(), 0)
	}

	index := Index{}
	if err := f.Write(index); err != nil {
		t.Errorf("File.Write() error = %v", err)
	}

	if f.LastIndex() != 1 {
		t.Errorf("File.LastIndex() = %v, want %v", f.LastIndex(), 1)
	}
}
