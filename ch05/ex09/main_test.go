package main

import "testing"

func Test_expand(t *testing.T) {
	ts := []struct {
		text     string
		fn       func(string) string
		expected string
	}{
		{
			text: "!$foo!",
			fn: func(s string) string {
				return "(" + s + ")"
			},
			expected: "!(foo)!",
		},
		{
			text: "$foo is (foo)",
			fn: func(s string) string {
				return "(" + s + ")"
			},
			expected: "(foo) is (foo)",
		},
		{
			text: "$foo$bar",
			fn: func(s string) string {
				return "(" + s + ")"
			},
			expected: "(foo)(bar)",
		},
	}

	for _, tc := range ts {
		if got := expand(tc.text, tc.fn); got != tc.expected {
			t.Errorf("unexpected string. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
