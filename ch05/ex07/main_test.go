package main

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestOutline(t *testing.T) {
	ts := []struct {
		HTML string
	}{
		{`
<html>
<head>
</head>
<body>
	<h1>It works!</h1>
	<p><a href="https://gopl.example/">example</a></p>
	<div id="foo">
		<span id="bar">
		</span>
	</div>
</body>
</html>
`},
	}

	for _, tc := range ts {
		doc, err := html.Parse(strings.NewReader(tc.HTML))
		if err != nil {
			t.Error(err)
			continue
		}

		b := bytes.NewBuffer(nil)
		defaultWriter = b
		forEachNode(doc, startElement, endElement)

		if _, err := html.Parse(b); err != nil {
			t.Errorf("unexpected error. %v", err)
		}

	}
}
