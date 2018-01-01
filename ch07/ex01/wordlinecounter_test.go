package main

import (
	"io"
	"testing"
)

func TestWordLineCounter_Lines(t *testing.T) {
	ts := []struct {
		text     string
		expected int
	}{
		{
			text:     "",
			expected: 0,
		},
		{
			text:     "foo",
			expected: 0,
		},
		{
			text:     "bar\n",
			expected: 1,
		},
		{
			text:     "これはUTF-8においても\n行数が数えられるかの\nサンプルです\n\n",
			expected: 4,
		},
	}

	for _, tc := range ts {
		wc := &WordLineCounter{}
		io.WriteString(wc, tc.text)

		if got := wc.Lines(); got != tc.expected {
			t.Errorf("unexpected lines. expected: %v, but got: %v", tc.expected, got)
		}
	}
}

func TestWordLineCounter_Words(t *testing.T) {
	ts := []struct {
		text     string
		expected int
	}{
		{
			text:     "",
			expected: 0,
		},
		{
			text:     "foo",
			expected: 1,
		},
		{
			text:     "bar\n",
			expected: 1,
		},
		{
			text:     "日\t本 語\nサ  ン  プ  ル",
			expected: 7,
		},
	}

	for _, tc := range ts {
		wc := &WordLineCounter{}
		io.WriteString(wc, tc.text)
		if got := wc.Words(); got != tc.expected {
			t.Errorf("unexpected words. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
