package tar

import (
	"archive/tar"
	"io"
	"io/ioutil"

	"github.com/orisano/gopl/ch10/ex02/archive"
)

func init() {
	archive.RegisterFormat("tar", []byte("ustar.00"), Unarchive)
	archive.RegisterFormat("tar", []byte("ustar . "), Unarchive)
}

type Archive struct {
	tr *tar.Reader
}

func (a *Archive) Next() (*archive.File, error) {
	h, err := a.tr.Next()
	if err == io.EOF {
		return nil, archive.EOA
	}
	if err != nil {
		return nil, err
	}
	return &archive.File{
		Name: h.Name,
		Body: ioutil.NopCloser(a.tr),
	}, nil
}

func Unarchive(r io.Reader) (archive.Archive, error) {
	tr := tar.NewReader(r)
	return &Archive{tr: tr}, nil
}
