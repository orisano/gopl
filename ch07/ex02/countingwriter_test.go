package main

import (
	"bytes"
	"io"
	"testing"
)

func TestCountingWriter(t *testing.T) {
	ts := []struct {
		text string
	}{
		{
			text: "",
		},
		{
			text: "foo",
		},
		{
			text: "foo bar",
		},
		{
			text: "日本語",
		},
	}

	for _, tc := range ts {
		var b bytes.Buffer

		cw, nbyte := CountingWriter(&b)
		if *nbyte != 0 {
			t.Errorf("unexpected count. expected: %v, but got: %v", 0, *nbyte)
		}
		io.WriteString(cw, tc.text)

		if *nbyte != int64(len(tc.text)) {
			t.Errorf("unexpected count. expected: %v, but got: %v", len(tc.text), *nbyte)
		}

		if got := b.String(); got != tc.text {
			t.Errorf("unexpected string. expected: %q, but got: %q", tc.text, got)
		}
	}
}
