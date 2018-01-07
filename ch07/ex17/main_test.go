package main

import (
	"encoding/xml"
	"testing"
)

func TestAttributeSelector_Match(t *testing.T) {
	ts := []struct {
		selector *AttributeSelector
		element  xml.StartElement
		expected bool
	}{
		{
			selector: &AttributeSelector{Name: "id", Value: "foo"},
			element: xml.StartElement{
				Name: xml.Name{Local: "p"},
				Attr: []xml.Attr{
					{xml.Name{Local: "id"}, "foo"},
				},
			},
			expected: true,
		},
		{
			selector: &AttributeSelector{Name: "id", Value: "foo"},
			element: xml.StartElement{
				Name: xml.Name{Local: "p"},
				Attr: []xml.Attr{
					{xml.Name{Local: "id"}, "bar"},
				},
			},
			expected: false,
		},
		{
			selector: &AttributeSelector{Name: "class", Value: "foo-class"},
			element: xml.StartElement{
				Name: xml.Name{Local: "h1"},
				Attr: []xml.Attr{
					{xml.Name{Local: "class"}, "foo-class"},
				},
			},
			expected: true,
		},
		{
			selector: &AttributeSelector{Name: "id", Value: "foo-class"},
			element: xml.StartElement{
				Name: xml.Name{Local: "h1"},
				Attr: []xml.Attr{
					{xml.Name{Local: "class"}, "bar-class"},
				},
			},
			expected: false,
		},
		{
			selector: &AttributeSelector{Name: "id", Value: "foo"},
			element: xml.StartElement{
				Name: xml.Name{Local: "p"},
			},
			expected: false,
		},
		{
			selector: &AttributeSelector{Name: "id", Value: "foo"},
			element: xml.StartElement{
				Name: xml.Name{Local: "img"},
				Attr: []xml.Attr{
					{xml.Name{Local: "src"}, "http://gopl.example/image.png"},
					{xml.Name{Local: "id"}, "foo"},
				},
			},
			expected: true,
		},
	}

	for _, tc := range ts {
		if got := tc.selector.Match(tc.element); got != tc.expected {
			t.Errorf("unexpected result. expected: %v, but got: %v", tc.expected, got)
		}
	}
}

func TestQuerySelector(t *testing.T) {
	ts := []struct {
		elements  []xml.StartElement
		selectors []Selector
		expected  bool
	}{
		{
			elements: []xml.StartElement{
				{Attr: []xml.Attr{{xml.Name{Local: "class"}, "foo"}}},
			},
			selectors: []Selector{
				&AttributeSelector{"class", "foo"},
			},
			expected: true,
		},
		{
			elements: []xml.StartElement{
				{},
				{},
				{Attr: []xml.Attr{{xml.Name{Local: "id"}, "bar"}}},
				{},
				{},
			},
			selectors: []Selector{
				&AttributeSelector{"id", "bar"},
			},
			expected: true,
		},
	}

	for _, tc := range ts {
		if got := querySelector(tc.elements, tc.selectors); got != tc.expected {
			t.Errorf("unexpected result. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
