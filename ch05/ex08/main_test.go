package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func buildHTML(text string) string {
	return "<html><head></head><body>" + text + "</body></html>"
}

func getID(n *html.Node) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == "id" {
			return a.Val, true
		}
	}
	return "", false
}

func TestElementByID(t *testing.T) {
	ts := []struct {
		HTML  string
		ID    string
		IsNil bool
		Tag   string
	}{
		{
			HTML: buildHTML(`<p id="test"></p>`),
			ID:   "test",
			Tag:  "p",
		},
		{
			HTML:  buildHTML(``),
			ID:    "test",
			IsNil: true,
		},
		{
			HTML: buildHTML(`<p id="foo"><a id="bar"></a></p>`),
			ID:   "bar",
			Tag:  "a",
		},
		{
			HTML: buildHTML(`<h1 id="conflict">head</h1><p id="conflict">text</p>`),
			ID:   "conflict",
			Tag:  "h1",
		},
	}

	for _, tc := range ts {
		t.Run("Case", func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.HTML))
			if err != nil {
				t.Fatal(err)
			}
			n := ElementByID(doc, tc.ID)

			if tc.IsNil != (n == nil) {
				t.Fatalf("unexpected node. expected: IsNil(%v), but got: %v", tc.IsNil, n)
			}
			if n != nil {
				id, ok := getID(n)
				if !ok {
					t.Fatal("id attribute is not found")
				}
				if id != tc.ID {
					t.Errorf("unexpected id. expected: %v, but got: %v", tc.ID, id)
				}
				if n.Data != tc.Tag {
					t.Errorf("unexpected tag. expected: %v, but got: %v", tc.Tag, n.Data)
				}
			}
		})
	}
}
