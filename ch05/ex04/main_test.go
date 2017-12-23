package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func equals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildHTML(text string) string {
	return `<html><head></head><body>` + text + `</body></html>`
}

func Test_visit(t *testing.T) {
	ts := []struct {
		HTML     string
		expected []string
	}{
		{
			buildHTML(""),
			[]string{},
		},
		{
			buildHTML(`<a href="https://google.com"></a>`),
			[]string{"https://google.com"},
		},
		{
			buildHTML(`<script src="https://gopl.example/bundle.js"></script>`),
			[]string{"https://gopl.example/bundle.js"},
		},
		{
			buildHTML(`<link rel="stylesheet" href="https://gopl.example/style.css" />`),
			[]string{"https://gopl.example/style.css"},
		},
		{
			buildHTML(`<img src="https://gopl.example/findlinks.png" />`),
			[]string{"https://gopl.example/findlinks.png"},
		},
		{
			buildHTML(`<a href="https://google.com"></a><a href="https://yahoo.com"></a>`),
			[]string{"https://google.com", "https://yahoo.com"},
		},
		{
			buildHTML(`
<link rel="stylesheet" href="http://gopl.example/style.css" />
<a href="https://google.com"></a>
<a href="/">Home</a>
<a href="#">#</a>
<div>
	<a href="https://facebook.com"></a>
	<a href="https://twitter.com"></a>
	<a href="https://golang.org"></a>
	<div>
		<a href="https://github.com"></a>
	</div>
	<img src="http://gopl.example/findlinks.png" />
</div>
<a href="https://yahoo.com"></a>
<script src="http://gopl.example/bundle.js"></script>
<script>
alert('http://example.example');
</script>
`),
			[]string{
				"http://gopl.example/style.css",
				"https://google.com",
				"/",
				"#",
				"https://facebook.com",
				"https://twitter.com",
				"https://golang.org",
				"https://github.com",
				"http://gopl.example/findlinks.png",
				"https://yahoo.com",
				"http://gopl.example/bundle.js",
			},
		},
	}
	for _, tc := range ts {
		doc, _ := html.Parse(strings.NewReader(tc.HTML))
		got := visit(nil, doc)
		if !equals(got, tc.expected) {
			t.Errorf("unexpected links. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
