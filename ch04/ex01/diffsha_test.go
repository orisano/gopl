package main

import "testing"

func TestShaDiffBit(t *testing.T) {
	ts := []struct {
		a, b     []byte
		expected int
	}{
		{[]byte{0x10}, []byte{0x10}, 0},
		{[]byte{0x01}, []byte{0x10}, 2},
		{[]byte{0x10, 0x20}, []byte{0x10, 0x30}, 1},
	}
	for _, tc := range ts {
		if got := ShaDiffBit(tc.a, tc.b); got != tc.expected {
			t.Errorf("unexpected diff bits. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
