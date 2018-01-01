package treesort

import "testing"

func TestTree_String(t *testing.T) {
	ts := []struct {
		t        *tree
		expected string
	}{
		{
			t:        nil,
			expected: "",
		},
		{
			t:        &tree{value: 10},
			expected: "10",
		},
		{
			t: &tree{
				value: 10,
				left:  &tree{value: 9},
			},
			expected: "9 10",
		},
		{
			t: &tree{
				value: 10,
				right: &tree{value: 11},
			},
			expected: "10 11",
		},
		{
			t: &tree{
				value: 2,
				left:  &tree{value: 1},
				right: &tree{value: 3},
			},
			expected: "1 2 3",
		},
	}

	for _, tc := range ts {
		if got := tc.t.String(); got != tc.expected {
			t.Errorf("unexpected string. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
