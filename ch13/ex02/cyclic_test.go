package cyclic

import "testing"

func TestIsCyclic(t *testing.T) {
	type Link struct {
		next *Link
	}
	a, b, c := &Link{}, &Link{}, &Link{}
	a.next, b.next, c.next = b, a, c
	d, e, f := &Link{}, &Link{}, &Link{}
	d.next, e.next = e, f

	tests := []struct {
		in       interface{}
		expected bool
	}{
		{in: a, expected: true},
		{in: b, expected: true},
		{in: c, expected: true},
		{in: d, expected: false},
		{in: e, expected: false},
		{in: f, expected: false},
	}

	for _, test := range tests {
		got := IsCyclic(test.in)
		if got != test.expected {
			t.Errorf("unexpected result. expected: %v, but got: %v", test.expected, got)
		}
	}
}
