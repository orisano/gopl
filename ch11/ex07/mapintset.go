package intset

import (
	"sort"
	"strconv"
)

type MapIntSet struct {
	set map[int]bool
}

func (m *MapIntSet) Has(x int) bool {
	if m.set == nil {
		return false
	}
	return m.set[x]
}

func (m *MapIntSet) Add(x int) {
	if m.set == nil {
		m.set = make(map[int]bool)
	}
	m.set[x] = true
}

func (m *MapIntSet) UnionWith(t *MapIntSet) {
	if t.set != nil && m.set == nil {
		m.set = make(map[int]bool)
	}
	for x := range t.set {
		m.set[x] = true
	}
}

func (m *MapIntSet) String() string {
	var xs []int
	for x := range m.set {
		xs = append(xs, x)
	}
	sort.Ints(xs)

	b := []byte{'{'}
	for i, x := range xs {
		if i > 0 {
			b = append(b, ' ')
		}
		b = strconv.AppendInt(b, int64(x), 10)
	}
	b = append(b, '}')
	return string(b)
}
