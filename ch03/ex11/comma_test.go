package main

import (
	"testing"
)

func Test_comma(t *testing.T) {
	ts := []struct {
		s        string
		expected string
	}{
		{"", ""},
		{"1", "1"},
		{"11", "11"},
		{"111", "111"},
		{"1234", "1,234"},
		{"12345678", "12,345,678"},
		{"+12345678", "+12,345,678"},
		{"-1234567", "-1,234,567"},
		{"-1234.5678", "-1,234.5678"},
	}
	for _, tc := range ts {
		if got := comma(tc.s); got != tc.expected {
			t.Errorf("unexpected result. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
