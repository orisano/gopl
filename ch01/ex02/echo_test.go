package main

import (
	"bytes"
	"testing"
)

func TestEcho(t *testing.T) {
	var b bytes.Buffer
	Echo(&b, []string{"./echo", "a", "b", "c", "de", "fgh"})

	got := b.String()
	expected := "0 a\n1 b\n2 c\n3 de\n4 fgh\n"
	if got != expected {
		t.Error("unexpected output. expected: %v, got: %v", expected, got)
	}
}
