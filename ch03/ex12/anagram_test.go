package main

import "testing"

func TestIsAnagram(t *testing.T) {
	ts := []struct {
		a, b     string
		expected bool
	}{
		{"a", "b", false},
		{"a", "a", true},
		{"bab", "abb", true},
		{"babb", "abb", false},
		{"aacc", "accaz", false},
		{"ふー", "ばー", false},
		{"みかん", "かんみ", true},
	}
	for _, tc := range ts {
		if got := IsAnagram(tc.a, tc.b); got != tc.expected {
			t.Errorf("unexpected result. expected: %v, but got: IsAnagram(%v, %v) #=> %v", tc.expected, tc.a, tc.b, got)
		}
	}
}
