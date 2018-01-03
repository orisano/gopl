package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestLimitReader(t *testing.T) {
	ts := []struct {
		text     string
		limit    int64
		expected string
	}{
		{
			text:     "",
			limit:    10,
			expected: "",
		},
		{
			text:     "日本語",
			limit:    int64(len([]byte("日本"))),
			expected: "日本",
		},
		{
			text:     strings.Repeat("Ab", 100000),
			limit:    6,
			expected: "AbAbAb",
		},
	}

	for _, tc := range ts {
		sr := strings.NewReader(tc.text)
		r := LimitReader(sr, tc.limit)

		var b bytes.Buffer
		if _, err := io.Copy(&b, r); err != nil {
			t.Error(err)
			continue
		}
		if got := b.String(); got != tc.expected {
			t.Errorf("unexpected result. expected: %q, but got: %q", tc.expected, got)
		}
	}
}
