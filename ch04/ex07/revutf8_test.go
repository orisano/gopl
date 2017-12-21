package main

import (
	"testing"
)

func TestReverseUTF8(t *testing.T) {
	ts := []struct {
		s        string
		expected string
	}{
		{"foo", "oof"},
		{"テスト", "トステ"},
		{"テsuト", "トusテ"},
	}

	for _, tc := range ts {
		b := []byte(tc.s)
		ReverseUTF8(b)
		if got := string(b); got != tc.expected {
			t.Errorf("reverse failed. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
