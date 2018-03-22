package ex05

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		s, sep string
		want   int
	}{
		{"a:b:c", ":", 3},
		{"1 2 3 4 5", " ", 5},
		{"", "\n", 1},
		{"one\ttwo", "\t", 2},
	}
	for _, test := range tests {
		words := strings.Split(test.s, test.sep)
		if got, want := len(words), test.want; got != want {
			t.Errorf("Split(%q, %q) returned %d words, want %d", test.s, test.sep, got, want)
		}
	}
}
