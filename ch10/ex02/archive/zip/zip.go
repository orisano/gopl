package zip

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/orisano/gopl/ch10/ex02/archive"
)

func init() {
	archive.RegisterFormat("zip", []byte{'P', 'K', 0x05, 0x06}, Unarchive)
}

type Archive struct {
	zr   *zip.Reader
	curr int
}

func (a *Archive) Next() (*archive.File, error) {
	if a.curr >= len(a.zr.File) {
		return nil, archive.EOA
	}
	f := a.zr.File[a.curr]
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	return &archive.File{
		Name: f.Name,
		Body: rc,
	}, nil
}

type sizeReaderAt interface {
	io.ReaderAt
	Size() int64
}

type File struct {
	*os.File
	size int64
}

func (f *File) Size() int64 { return f.size }

func asSizeReaderAt(r io.Reader) (sizeReaderAt, error) {
	if sra, ok := r.(sizeReaderAt); ok {
		return sra, nil
	}
	if f, ok := r.(*os.File); ok {
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}
		return &File{File: f, size: stat.Size()}, nil
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func Unarchive(r io.Reader) (archive.Archive, error) {
	sra, err := asSizeReaderAt(r)
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(sra, sra.Size())
	if err != nil {
		return nil, err
	}
	return &Archive{zr: zr}, nil
}
