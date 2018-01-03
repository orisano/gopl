package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	ts := []string{
		"",
		"foo",
		"foo bar",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"日本語",
		strings.Repeat("A", 10000),
	}

	for _, tc := range ts {
		r := NewReader(tc)
		var b bytes.Buffer
		n, err := io.Copy(&b, r)
		if err != nil {
			t.Error(err)
			continue
		}
		if got := int(n); got != len(tc) {
			t.Errorf("unexpected write bytes. expected: %v, but got: %v", len(tc), got)
		}
		if got := b.String(); got != tc {
			t.Errorf("unexpected result. expected: %q, but got: %q", tc, got)
		}
	}
}
