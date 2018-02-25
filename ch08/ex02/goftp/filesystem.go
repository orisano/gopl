package goftp

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FileSystem interface {
	Get(path string) (io.ReadCloser, error)
	Create(path string) (io.WriteCloser, error)
	List(path string) ([]string, error)
	LS(path string) ([]string, error)
	IsDir(path string) bool
}

type RawFileSystem struct {
	Root string
}

var _ FileSystem = &RawFileSystem{}

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

func (f *RawFileSystem) Create(path string) (io.WriteCloser, error) {
	p := f.resolve(path)
	file, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *RawFileSystem) List(path string) ([]string, error) {
	p := f.resolve(path)
	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var stats []os.FileInfo
	if stat.IsDir() {
		ss, err := file.Readdir(0)
		if err != nil {
			return nil, err
		}
		stats = ss
	} else {
		stats = append(stats, stat)
	}

	fileNames := make([]string, 0, len(stats))
	for _, s := range stats {
		fileNames = append(fileNames, s.Name())
	}
	return fileNames, nil
}

func (f *RawFileSystem) LS(path string) ([]string, error) {
	p := f.resolve(path)
	cmd := exec.Command("/bin/ls", "-lah", p)
	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf
	cmd.Stderr = buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return strings.Split(buf.String(), "\n"), nil
}

func (f *RawFileSystem) IsDir(path string) bool {
	p := f.resolve(path)
	stat, err := os.Stat(p)
	if err != nil {
		return false
	}
	return stat.IsDir()
}