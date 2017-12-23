package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func equals(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func Test_countTag(t *testing.T) {
	ts := []struct {
		HTML     string
		expected map[string]int
	}{
		{
			`<html><head></head><body></body></html>`,
			map[string]int{"html": 1, "head": 1, "body": 1},
		},
		{
			`<html><head></head><body><p></p><p></p><span><div><p></p></div><p></p><span></span><pre></pre></span></body></html>`,
			map[string]int{"html": 1, "head": 1, "body": 1, "p": 4, "div": 1, "span": 2, "pre": 1},
		},
	}

	for _, tc := range ts {
		t.Run("Case", func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.HTML))
			if err != nil {
				t.Fatal(err)
			}
			got := map[string]int{}
			countTag(got, doc)
			if !equals(got, tc.expected) {
				t.Errorf("unexpected tag count. expected: %v, but got: %v", tc.expected, got)
			}
		})
	}
}
