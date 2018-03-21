package archive

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

var EOA = errors.New("EOA")

type Archive interface {
	Next() (*File, error)
}

type File struct {
	Name string
	Body io.ReadCloser
}

type format struct {
	name      string
	magic     []byte
	unarchive func(io.Reader) (Archive, error)
}

var formats []format

func RegisterFormat(name string, magic []byte, unarchive func(io.Reader) (Archive, error)) {
	formats = append(formats, format{
		name:      name,
		magic:     magic,
		unarchive: unarchive,
	})
}

type peeker interface {
	Peek(int) ([]byte, error)
}

type atPeeker struct {
	io.ReaderAt
}

func (p *atPeeker) Peek(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := p.ReadAt(b, 0)
	return b, err
}

func asPeeker(r io.Reader) (io.Reader, peeker) {
	switch t := r.(type) {
	case peeker:
		return r, t
	case io.ReaderAt:
		return r, &atPeeker{t}
	default:
		br := bufio.NewReader(r)
		return br, br
	}
}

func sniff(p peeker) format {
	for _, f := range formats {
		b, err := p.Peek(len(f.magic))
		if err == nil && bytes.Equal(b, f.magic) {
			return f
		}
	}
	return format{}
}

func Unarchive(r io.Reader) (Archive, string, error) {
	r, p := asPeeker(r)
	f := sniff(p)
	if f.unarchive == nil {
		return nil, "", errors.New("unknown format")
	}
	a, err := f.unarchive(r)
	if err != nil {
		return nil, f.name, err
	}
	return a, f.name, nil
}
