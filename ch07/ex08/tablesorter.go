package ex08

import "sort"

type TableSorter struct {
	prevSorter, currSorter sort.Interface
}

func (s *TableSorter) SetSorter(sorter sort.Interface) {
	s.prevSorter = s.currSorter
	s.currSorter = sorter
}

func (s *TableSorter) Len() int {
	return s.currSorter.Len()
}

func (s *TableSorter) Less(i, j int) bool {
	if !equals(s.currSorter, i, j) {
		return s.currSorter.Less(i, j)
	}
	if s.prevSorter != nil && !equals(s.prevSorter, i, j) {
		return s.prevSorter.Less(i, j)
	}
	return false
}

func (s *TableSorter) Swap(i, j int) {
	s.currSorter.Swap(i, j)
}

func equals(sorter sort.Interface, i, j int) bool {
	return !sorter.Less(i, j) && !sorter.Less(j, i)
}
