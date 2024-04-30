package file

import (
	"fmt"
	"os"
	"time"
)

const (
	FileName      string = "log"
	FilePermition        = 0644
)

type File struct {
	rawFile *os.File
	path    string
	size    int
}

func NewFile(basePath string, id uint64) (*File, error) {
	path := makeFilePath(basePath, id)
	rawFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, FilePermition)
	if err != nil {
		return nil, err
	}

	info, err := rawFile.Stat()
	if err != nil {
		return nil, err
	}

	f := &File{
		rawFile: rawFile,
		path:    path,
		size:    int(info.Size()),
	}
	return f, nil
}

func (f *File) Write(buf []byte) error {
	n, err := f.rawFile.Write(buf)
	f.size += n
	return err
}

func (f *File) WriteWithSync(buf []byte) error {
	if err := f.Write(buf); err != nil {
		return err
	}
	return f.Sync()
}

func (f *File) ReadAll() ([]byte, error) {
	return os.ReadFile(f.path)
}

func (f *File) Read() ([]byte, error) {
	buf := []byte{}
	_, err := f.rawFile.Read(buf)
	return buf, err
}

func (f *File) Sync() error {
	return f.rawFile.Sync()
}

func (f *File) Close() error {
	if err := f.Sync(); err != nil {
		return err
	}
	return f.rawFile.Close()
}

func makeFilePath(basePath string, id uint64) string {
	return fmt.Sprintf("%s/%s_%d_%s", basePath, FileName, id, time.Now().String())
}
