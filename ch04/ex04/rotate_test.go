package main

import "testing"

func equals(a, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRotate(t *testing.T) {
	ts := []struct {
		s        []int
		r        int
		expected []int
	}{
		{
			[]int{0, 1, 2, 3, 4, 5},
			1,
			[]int{5, 0, 1, 2, 3, 4},
		},
		{
			[]int{0, 1, 2, 3, 4, 5},
			1,
			[]int{5, 0, 1, 2, 3, 4},
		},
		{
			[]int{0, 1, 2, 3, 4, 5},
			2,
			[]int{4, 5, 0, 1, 2, 3},
		},
		{
			[]int{0, 1, 2, 3, 4, 5},
			3,
			[]int{3, 4, 5, 0, 1, 2},
		},
		{
			[]int{1, 9, 2, 3, 5, 1, 4, 6},
			2,
			[]int{4, 6, 1, 9, 2, 3, 5, 1},
		},
	}

	for _, tc := range ts {
		Rotate(tc.s, tc.r)
		if !equals(tc.s, tc.expected) {
			t.Errorf("unexpected slice. expected: %v, but got: %v", tc.expected, tc.s)
		}
	}
}
