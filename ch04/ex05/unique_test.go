package main

import "testing"

func equals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestUnique(t *testing.T) {
	ts := []struct {
		ss       []string
		expected []string
	}{
		{
			[]string{"a", "a", "a", "a"},
			[]string{"a"},
		},
		{
			[]string{"a", "b", "c", "a"},
			[]string{"a", "b", "c", "a"},
		},
		{
			[]string{"a", "b", "b", "a"},
			[]string{"a", "b", "a"},
		},
	}
	for _, tc := range ts {
		if got := Unique(tc.ss); !equals(got, tc.expected) {
			t.Errorf("unexpected slice. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
