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

func Test_visit(t *testing.T) {
	ts := []struct {
		HTML     string
		expected []string
	}{
		{
			`<html><body></body></html>`,
			[]string{},
		},
		{
			`<html><body><a href="https://google.com"></a></body></html>`,
			[]string{"https://google.com"},
		},
		{
			`<html>
<body>
	<a href="https://google.com"></a>
	<a href="https://yahoo.com"></a>
</body>
</html>`,
			[]string{"https://google.com", "https://yahoo.com"},
		},
		{
			`<html>
<body>
	<a href="https://google.com"></a>
	<div>
		<a href="https://facebook.com"></a>
		<a href="https://twitter.com"></a>
		<a href="https://golang.org"></a>
		<div>
			<a href="https://github.com"></a>
		</div>
	</div>
	<a href="https://yahoo.com"></a>
</body>
</html>`,
			[]string{"https://google.com", "https://facebook.com", "https://twitter.com", "https://golang.org", "https://github.com", "https://yahoo.com"},
		},
	}
	for _, tc := range ts {
		t.Run("Case", func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.HTML))
			if err != nil {
				t.Fatal(err)
			}
			got := visit(nil, doc)
			if !equals(got, tc.expected) {
				t.Errorf("unexpected links. expected: %v, but got: %v", tc.expected, got)
			}
		})
	}
}
