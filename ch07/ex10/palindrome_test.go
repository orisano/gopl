package ex10

import (
	"sort"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	ts := []struct {
		sorter   sort.Interface
		expected bool
	}{
		{
			sorter:   sort.IntSlice([]int{}),
			expected: true,
		},
		{
			sorter:   sort.IntSlice([]int{1}),
			expected: true,
		},
		{
			sorter:   sort.IntSlice([]int{1, 2}),
			expected: false,
		},
		{
			sorter:   sort.IntSlice([]int{2, 2}),
			expected: true,
		},
		{
			sorter:   sort.IntSlice([]int{1, 2, 3}),
			expected: false,
		},
		{
			sorter:   sort.IntSlice([]int{1, 2, 1}),
			expected: true,
		},
	}
	for i, tc := range ts {
		if got := IsPalindrome(tc.sorter); got != tc.expected {
			t.Errorf("unexpected result #%v. expected: %v, but got: %v", i, tc.expected, got)
		}
	}
}
