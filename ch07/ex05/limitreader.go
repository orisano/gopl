package main

import "io"

type limitReader struct {
	r io.Reader
	n int64
}

func (r *limitReader) Read(p []byte) (n int, err error) {
	if int64(len(p)) > r.n {
		p = p[:r.n]
	}
	n, err = r.r.Read(p)
	r.n -= int64(n)
	if r.n == 0 && err == nil {
		err = io.EOF
	}
	return n, err
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &limitReader{
		r: r,
		n: n,
	}
}
