package main

import (
	"reflect"
	"testing"
)

func TestUnique(t *testing.T) {
	ts := []struct {
		strs     []string
		expected []string
	}{
		{
			strs:     []string{"", "a", "b", "c", "d"},
			expected: []string{"", "a", "b", "c", "d"},
		},
		{
			strs:     []string{"a", "a", "a"},
			expected: []string{"a"},
		},
		{
			strs:     []string{"a", "b", "b", "a", "c", "a", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			strs:     []string{"b", "b", "a", "c", "a", "c"},
			expected: []string{"b", "a", "c"},
		},
	}
	for _, tc := range ts {
		got := Unique(tc.strs)
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("unexpected result. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
