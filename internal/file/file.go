package file

import (
	"io/fs"
	"os"
)

type File struct {
	rawFile *os.File
	path    string
	size    int
}

func NewFile(path string, mode fs.FileMode) (*File, error) {
	rawFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, mode)
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
