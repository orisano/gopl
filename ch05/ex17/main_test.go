package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func getAttr(n *html.Node, name string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == name {
			return attr.Val, true
		}
	}
	return "", false
}

func buildHTML(text string) string {
	return `<html><head></head><body>` + text + `</body></html>`
}

func TestElementsByTagName(t *testing.T) {
	type result struct {
		tag, id string
	}
	ts := []struct {
		HTML     string
		tags     []string
		expected []result
	}{
		{
			HTML: buildHTML(`<img src="https://gopl.example/foo.png" id="1" />`),
			tags: []string{"img"},
			expected: []result{
				{tag: "img", id: "1"},
			},
		},
		{
			HTML: buildHTML(`
<h1 id="1">It works!</h1>
<img src="https://gopl.example/foo.png" id="2" />
<div id="3">
	<h2 id="4">foo</h2>
	<h2 id="5">bar</h2>
	<span id="6">
		<h3 id="7">example</h3>
		<h3 id="8">test</h3>
		<div id="9">
			<h4 id="10">nested</h4>
			<h5 id="11"></h5>
		</div>
	</span>
	<h2 id="12">foo</h2>
</div>
`),
			tags: []string{"h1", "h2", "h3", "h4"},
			expected: []result{
				{tag: "h1", id: "1"},
				{tag: "h2", id: "4"},
				{tag: "h2", id: "5"},
				{tag: "h3", id: "7"},
				{tag: "h3", id: "8"},
				{tag: "h4", id: "10"},
				{tag: "h2", id: "12"},
			},
		},
	}

	for _, tc := range ts {
		doc, err := html.Parse(strings.NewReader(tc.HTML))
		if err != nil {
			t.Error(err)
			continue
		}
		nodes := ElementsByTagName(doc, tc.tags...)

		if len(nodes) != len(tc.expected) {
			t.Errorf("unexpected nodes length. expected: %v, but got: %v", len(tc.expected), len(nodes))
			continue
		}
		for i, node := range nodes {
			if node.Data != tc.expected[i].tag {
				t.Errorf("unexpected tag name. expected: %v, but got: %v", tc.expected[i].tag, node.Data)
			}
			if got, ok := getAttr(node, "id"); ok {
				if got != tc.expected[i].id {
					t.Errorf("unexpected id. expected: %v, but got: %v", tc.expected[i].id, got)
				}
			} else {
				t.Errorf("id attribute not found")
			}
		}
	}
}
