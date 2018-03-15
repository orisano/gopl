package goftp2

import (
	"io"

	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("no such file or directory")

type FileSystem interface {
	Stat(path string) (*FileStat, error)
	Open(path string) (io.ReadCloser, error)
}

type FileStat struct {
	IsDir bool
}
