package ex08

import (
	"sort"
	"testing"
)

type columns struct {
	A string
	B string
	C int
}

type columnsSorter struct {
	table []*columns
	less  func(x, y *columns) bool
}

func (s *columnsSorter) Len() int {
	return len(s.table)
}

func (s *columnsSorter) Less(i, j int) bool {
	return s.less(s.table[i], s.table[j])
}

func (s *columnsSorter) Swap(i, j int) {
	s.table[i], s.table[j] = s.table[j], s.table[i]
}

func byA(x, y *columns) bool { return x.A < y.A }
func byB(x, y *columns) bool { return x.B < y.B }
func byC(x, y *columns) bool { return x.C < y.C }

func equalsTable(a, b []*columns) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if *a[i] != *b[i] {
			return false
		}
	}
	return true
}

func TestTableSorter(t *testing.T) {
	ts := []struct {
		table         []*columns
		lessFunctions []func(x, y *columns) bool
		expected      [][]*columns
	}{
		{
			table: []*columns{
				{A: "B"},
				{A: "A"},
				{A: ""},
				{A: "_"},
			},
			lessFunctions: []func(x, y *columns) bool{byA},
			expected: [][]*columns{
				{
					{A: ""},
					{A: "A"},
					{A: "B"},
					{A: "_"},
				},
			},
		},
		{
			table: []*columns{
				{A: "B", B: "bar"},
				{A: "A", B: "bar"},
				{A: "", B: "foo"},
				{A: "_", B: "bar"},
			},
			lessFunctions: []func(x, y *columns) bool{byA, byB},
			expected: [][]*columns{
				{
					{A: "", B: "foo"},
					{A: "A", B: "bar"},
					{A: "B", B: "bar"},
					{A: "_", B: "bar"},
				},
				{
					{A: "A", B: "bar"},
					{A: "B", B: "bar"},
					{A: "_", B: "bar"},
					{A: "", B: "foo"},
				},
			},
		},
		{
			table: []*columns{
				{A: "B", B: "bar", C: -1},
				{A: "A", B: "bar", C: 3},
				{A: "", B: "foo", C: -1},
				{A: "_", B: "bar", C: 20},
			},
			lessFunctions: []func(x, y *columns) bool{byA, byB, byC},
			expected: [][]*columns{
				{
					{A: "", B: "foo", C: -1},
					{A: "A", B: "bar", C: 3},
					{A: "B", B: "bar", C: -1},
					{A: "_", B: "bar", C: 20},
				},
				{
					{A: "A", B: "bar", C: 3},
					{A: "B", B: "bar", C: -1},
					{A: "_", B: "bar", C: 20},
					{A: "", B: "foo", C: -1},
				},
				{
					{A: "B", B: "bar", C: -1},
					{A: "", B: "foo", C: -1},
					{A: "A", B: "bar", C: 3},
					{A: "_", B: "bar", C: 20},
				},
			},
		},
	}

	for i, tc := range ts {
		if len(tc.lessFunctions) != len(tc.expected) {
			t.Errorf("test case %v is broken", i)
			continue
		}

		tsorter := &TableSorter{}
		for j, less := range tc.lessFunctions {
			tsorter.SetSorter(&columnsSorter{tc.table, less})
			sort.Sort(tsorter)

			if !equalsTable(tc.table, tc.expected[j]) {
				t.Errorf("unexpected order %v. expected: %#v, but got: %#v", j, tc.expected[j], tc.table)
			}
		}
	}
}

func ExampleTableSorter() {
	table := []*columns{
		{"A", "B", 2},
		{"B", "A", 5},
		{"1", "foo", 1},
		{"2", "bar", 4},
	}
	tsorter := &TableSorter{}

	tsorter.SetSorter(&columnsSorter{table, byA})
	sort.Sort(tsorter)

	tsorter.SetSorter(&columnsSorter{table, byB})
	sort.Sort(tsorter)

	tsorter.SetSorter(&columnsSorter{table, byC})
	sort.Sort(tsorter)

	// Output:
}

func ExampleSortStable() {
	table := []*columns{
		{"A", "B", 2},
		{"B", "A", 5},
		{"1", "foo", 1},
		{"2", "bar", 4},
	}

	sort.Stable(&columnsSorter{table, byA})
	sort.Stable(&columnsSorter{table, byB})
	sort.Stable(&columnsSorter{table, byC})

	// Output:
}
