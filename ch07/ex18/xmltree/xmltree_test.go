package xmltree

import (
	"encoding/xml"
	"reflect"
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	ts := []struct {
		xml      string
		expected []Node
	}{
		{
			xml:      "",
			expected: nil,
		},
		{
			xml: `<xml></xml>`,
			expected: []Node{
				&Element{
					Type: xml.Name{Local: "xml"},
					Attr: []xml.Attr{},
				},
			},
		},
		{
			xml: "<p>foo</p>",
			expected: []Node{
				&Element{
					Type: xml.Name{Local: "p"},
					Attr: []xml.Attr{},
					Children: []Node{
						CharData("foo"),
					},
				},
			},
		},
		{
			xml: `<h1 id="heading">foo</h1>`,
			expected: []Node{
				&Element{
					Type: xml.Name{Local: "h1"},
					Attr: []xml.Attr{
						{Name: xml.Name{Local: "id"}, Value: "heading"},
					},
					Children: []Node{
						CharData("foo"),
					},
				},
			},
		},
		{
			xml: `
<html lang="ja">
<head></head>
<body>
	<h1>It works!</h1>
</body>
</html>`,
			expected: []Node{
				CharData("\n"),
				&Element{
					Type: xml.Name{Local: "html"},
					Attr: []xml.Attr{
						{Name: xml.Name{Local: "lang"}, Value: "ja"},
					},
					Children: []Node{
						CharData("\n"),
						&Element{
							Type: xml.Name{Local: "head"},
							Attr: []xml.Attr{},
						},
						CharData("\n"),
						&Element{
							Type: xml.Name{Local: "body"},
							Attr: []xml.Attr{},
							Children: []Node{
								CharData("\n\t"),
								&Element{
									Type: xml.Name{Local: "h1"},
									Attr: []xml.Attr{},
									Children: []Node{
										CharData("It works!"),
									},
								},
								CharData("\n"),
							},
						},
						CharData("\n"),
					},
				},
			},
		},
		{
			xml: `<x:foo></x:foo>`,
			expected: []Node{
				&Element{
					Type: xml.Name{Space: "x", Local: "foo"},
					Attr: []xml.Attr{},
				},
			},
		},
	}

	for _, tc := range ts {
		r := strings.NewReader(tc.xml)
		got, err := Build(r)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("unexpected tree. expected: %s, but got: %s", tc.expected, got)
		}
	}
}
