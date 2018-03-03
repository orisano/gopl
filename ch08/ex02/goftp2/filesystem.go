package goftp2

import "github.com/pkg/errors"

var ErrNotFound = errors.New("no such file or directory")

type FileSystem interface {
	Stat(path string) (*FileStat, error)
}

type FileStat struct {
	IsDir bool
}
