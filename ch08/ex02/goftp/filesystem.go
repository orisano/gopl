package goftp

import (
	"io"
	"os"
	"path/filepath"
)

type FileSystem interface {
	Get(path string) (io.ReadCloser, error)
}

type RawFileSystem struct {
	Root string
}

func (f *RawFileSystem) resolve(path string) string {
	return filepath.Join(f.Root, path)
}

func (f *RawFileSystem) Get(path string) (io.ReadCloser, error) {
	p := f.resolve(path)
	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	return file, nil
}
