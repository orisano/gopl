package intset

import (
	"testing"
)

func newMapIntSet(a ...int) *MapIntSet {
	intSet := &MapIntSet{}
	for _, x := range a {
		intSet.Add(x)
	}
	return intSet
}

func TestMapIntSet_Has(t *testing.T) {
	ts := []struct {
		intSet   *MapIntSet
		expected map[int]bool
	}{
		{
			intSet: newMapIntSet(1, 10, 100, 1000, 10000),
			expected: map[int]bool{
				0:     false,
				1:     true,
				2:     false,
				9:     false,
				10:    true,
				11:    false,
				99:    false,
				100:   true,
				101:   false,
				999:   false,
				1000:  true,
				1001:  false,
				9999:  false,
				10000: true,
				10001: false,
			},
		},
	}
	for _, tc := range ts {
		for x, expected := range tc.expected {
			if got := tc.intSet.Has(x); got != expected {
				t.Errorf("unexpected Has(%d). expected: %v, but got: %v", x, expected, got)
			}
		}
	}
}

func TestMapIntSet_Add(t *testing.T) {
	ts := []struct {
		intSet   *MapIntSet
		adds     []int
		expected map[int]bool
	}{
		{
			intSet: newMapIntSet(1, 10, 100, 1000),
			adds:   []int{0, 50, 10000},
			expected: map[int]bool{
				0:     true,
				1:     true,
				2:     false,
				9:     false,
				10:    true,
				11:    false,
				49:    false,
				50:    true,
				51:    false,
				99:    false,
				100:   true,
				101:   false,
				999:   false,
				1000:  true,
				1001:  false,
				9999:  false,
				10000: true,
				10001: false,
			},
		},
	}

	for _, tc := range ts {
		for _, x := range tc.adds {
			tc.intSet.Add(x)
		}
		for x, expected := range tc.expected {
			if got := tc.intSet.Has(x); got != expected {
				t.Errorf("unexpected Has(%d). expected: %v, but got: %v", x, expected, got)
			}
		}
	}
}

func TestMapIntSet_UnionWith(t *testing.T) {
	ts := []struct {
		from, to *MapIntSet
		expected map[int]bool
	}{
		{
			from: newMapIntSet(1, 10, 100, 1000),
			to:   newMapIntSet(10000, 50000, 100000),
			expected: map[int]bool{
				0:      false,
				1:      true,
				2:      false,
				9:      false,
				10:     true,
				11:     false,
				99:     false,
				100:    true,
				101:    false,
				999:    false,
				1000:   true,
				1001:   false,
				9999:   false,
				10000:  true,
				10001:  false,
				49999:  false,
				50000:  true,
				50001:  false,
				99999:  false,
				100000: true,
				100001: false,
			},
		},
		{
			from: newMapIntSet(1, 10, 100, 1000),
			to:   newMapIntSet(2, 10, 1000, 10000),
			expected: map[int]bool{
				0:     false,
				1:     true,
				2:     true,
				3:     false,
				9:     false,
				10:    true,
				11:    false,
				99:    false,
				100:   true,
				101:   false,
				999:   false,
				1000:  true,
				1001:  false,
				9999:  false,
				10000: true,
				10001: false,
			},
		},
		{
			from: newMapIntSet(1, 10, 100, 1000),
			to:   newMapIntSet(0, 100),
			expected: map[int]bool{
				0:    true,
				1:    true,
				2:    false,
				9:    false,
				10:   true,
				11:   false,
				99:   false,
				100:  true,
				101:  false,
				999:  false,
				1000: true,
				1001: false,
			},
		},
	}

	for _, tc := range ts {
		tc.to.UnionWith(tc.from)
		for x, expected := range tc.expected {
			if got := tc.to.Has(x); got != expected {
				t.Errorf("unexpected Has(%d). expected: %v, but got: %v", x, expected, got)
			}
		}
	}
}

func TestMapIntSet_String(t *testing.T) {
	ts := []struct {
		intSet   *MapIntSet
		expected string
	}{
		{
			intSet:   newMapIntSet(),
			expected: "{}",
		},
		{
			intSet:   newMapIntSet(10),
			expected: "{10}",
		},
		{
			intSet:   newMapIntSet(1, 2, 4, 8, 16, 32),
			expected: "{1 2 4 8 16 32}",
		},
	}

	for _, tc := range ts {
		if got := tc.intSet.String(); got != tc.expected {
			t.Errorf("unexpected string. expected: %q, but got: %q", tc.expected, got)
		}
	}
}
