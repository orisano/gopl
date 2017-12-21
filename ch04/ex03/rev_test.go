package main

import "testing"

func equals(a, b *[16]int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Test_reverse(t *testing.T) {
	ts := []struct {
		array    [16]int
		expected [16]int
	}{
		{
			[...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6},
			[...]int{6, 5, 4, 3, 2, 1, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
	}

	for _, tc := range ts {
		reverse(&tc.array)
		if !equals(&tc.array, &tc.expected) {
			t.Errorf("unexpected behavior. expected: %v, but got: %v", tc.expected, tc.array)
		}
	}
}
