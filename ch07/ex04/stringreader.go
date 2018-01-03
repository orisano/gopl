package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html"
)

type StringReader struct {
	s      string
	offset int
}

func (r *StringReader) Read(p []byte) (int, error) {
	n := copy(p, r.s[r.offset:])
	r.offset += n
	if n != len(p) {
		return n, io.EOF
	}
	return n, nil
}

func NewReader(s string) io.Reader {
	return &StringReader{
		s: s,
	}
}

func main() {
	for _, text := range os.Args[1:] {
		doc, err := html.Parse(NewReader(text))
		if err != nil {
			fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
			os.Exit(1)
		}
		for _, link := range visit(nil, doc) {
			fmt.Println(link)
		}
	}
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
