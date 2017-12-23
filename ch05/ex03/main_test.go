package main

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func buildHTML(text string) string {
	return `<html><head></head><body>` + text + `</body></html>`
}

func Test_renderTextNode(t *testing.T) {
	ts := []struct {
		HTML     string
		expected string
	}{
		{
			buildHTML(""),
			"",
		},
		{
			buildHTML("<h1>It works!"),
			"It works!",
		},
		{
			buildHTML("<h1>It works!</h1> <p>Index</p>"),
			"It works! Index",
		},
		{
			buildHTML(`<style>h1 { display: none; }</style><h1>It works!</h1>`),
			"It works!",
		},
		{
			buildHTML(`<h1>It works!</h1><script>alert('JavaScript');</script>`),
			"It works!",
		},
		{
			buildHTML(`<style>h1 { display: none; }</style><h1>It works!</h1><script>alert('JavaScript');</script>`),
			"It works!",
		},
	}

	for _, tc := range ts {
		doc, _ := html.Parse(strings.NewReader(tc.HTML))
		var b bytes.Buffer
		renderTextNode(&b, doc)
		if got := b.String(); got != tc.expected {
			t.Errorf("unexpected text node. expected: %q, but got: %q", tc.expected, got)
		}
	}
}
