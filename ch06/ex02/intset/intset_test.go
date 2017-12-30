package intset

import (
	"testing"
)

func newIntSet(a ...int) *IntSet {
	intSet := &IntSet{}
	for _, x := range a {
		intSet.Add(x)
	}
	return intSet
}

func TestIntSet_Has(t *testing.T) {
	ts := []struct {
		intSet   *IntSet
		expected map[int]bool
	}{
		{
			intSet: newIntSet(1, 10, 100, 1000, 10000),
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

func TestIntSet_Add(t *testing.T) {
	ts := []struct {
		intSet   *IntSet
		adds     []int
		expected map[int]bool
	}{
		{
			intSet: newIntSet(1, 10, 100, 1000),
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

func TestIntSet_AddAll(t *testing.T) {
	ts := []struct {
		intSet   *IntSet
		adds     [][]int
		expected map[int]bool
	}{
		{
			intSet: newIntSet(1, 10, 100, 1000),
			adds:   [][]int{{0, 1, 2, 3}, {50, 51}, {500, 10000}},
			expected: map[int]bool{
				0:     true,
				1:     true,
				2:     true,
				3:     true,
				4:     false,
				9:     false,
				10:    true,
				11:    false,
				49:    false,
				50:    true,
				51:    true,
				52:    false,
				99:    false,
				100:   true,
				101:   false,
				499:   false,
				500:   true,
				501:   false,
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
			tc.intSet.AddAll(x...)
		}
		for x, expected := range tc.expected {
			if got := tc.intSet.Has(x); got != expected {
				t.Errorf("unexpected Has(%d). expected: %v, but got: %v", x, expected, got)
			}
		}
	}
}

func TestIntSet_UnionWith(t *testing.T) {
	ts := []struct {
		from, to *IntSet
		expected map[int]bool
	}{
		{
			from: newIntSet(1, 10, 100, 1000),
			to:   newIntSet(10000, 50000, 100000),
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
			from: newIntSet(1, 10, 100, 1000),
			to:   newIntSet(2, 10, 1000, 10000),
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
			from: newIntSet(1, 10, 100, 1000),
			to:   newIntSet(0, 100),
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

func TestIntSet_String(t *testing.T) {
	ts := []struct {
		intSet   *IntSet
		expected string
	}{
		{
			intSet:   newIntSet(),
			expected: "{}",
		},
		{
			intSet:   newIntSet(10),
			expected: "{10}",
		},
		{
			intSet:   newIntSet(1, 2, 4, 8, 16, 32),
			expected: "{1 2 4 8 16 32}",
		},
	}

	for _, tc := range ts {
		if got := tc.intSet.String(); got != tc.expected {
			t.Errorf("unexpected string. expected: %q, but got: %q", tc.expected, got)
		}
	}
}

func TestIntSet_Len(t *testing.T) {
	ts := []struct {
		intSet   *IntSet
		expected int
	}{
		{
			intSet:   newIntSet(),
			expected: 0,
		},
		{
			intSet:   newIntSet(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
			expected: 10,
		},
		{
			intSet:   newIntSet(1, 100000),
			expected: 2,
		},
	}

	for _, tc := range ts {
		if got := tc.intSet.Len(); got != tc.expected {
			t.Errorf("unexpected len. expected: %v, but got: %v", tc.expected, got)
		}
	}
}

func TestIntSet_Remove(t *testing.T) {
	ts := []struct {
		intSet   *IntSet
		removes  []int
		expected map[int]bool
	}{
		{
			intSet:  newIntSet(1, 10, 100, 1000),
			removes: []int{1, 1000},
			expected: map[int]bool{
				0:    false,
				1:    false,
				2:    false,
				9:    false,
				10:   true,
				11:   false,
				99:   false,
				100:  true,
				101:  false,
				999:  false,
				1000: false,
				1001: false,
			},
		},
	}

	for _, tc := range ts {
		for _, x := range tc.removes {
			tc.intSet.Remove(x)
		}
		for x, expected := range tc.expected {
			if got := tc.intSet.Has(x); got != expected {
				t.Errorf("unexpected Has(%d). expected: %v, but got: %v", x, expected, got)
			}
		}
	}
}

func TestIntSet_Clear(t *testing.T) {
	ts := []struct {
		intSet   *IntSet
		expected map[int]bool
	}{
		{
			intSet: newIntSet(1, 10, 100, 1000),
			expected: map[int]bool{
				0:    false,
				1:    false,
				2:    false,
				9:    false,
				10:   false,
				11:   false,
				99:   false,
				100:  false,
				101:  false,
				999:  false,
				1000: false,
				1001: false,
			},
		},
	}

	for _, tc := range ts {
		tc.intSet.Clear()
		for x, expected := range tc.expected {
			if got := tc.intSet.Has(x); got != expected {
				t.Errorf("unexpected Has(%d). expected: %v, but got: %v", x, expected, got)
			}
		}
	}
}

func TestIntSet_Copy(t *testing.T) {
	ts := []struct {
		intSet   *IntSet
		expected map[int]bool
	}{
		{
			intSet: newIntSet(1, 10, 100, 1000),
			expected: map[int]bool{
				0:    false,
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
		intSet := tc.intSet.Copy()
		for x, expected := range tc.expected {
			if got := intSet.Has(x); got != expected {
				t.Errorf("unexpected Has(%d). expected: %v, but got: %v", x, expected, got)
			}
		}
	}
}
