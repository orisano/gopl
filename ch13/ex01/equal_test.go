package equal

import "testing"

func TestEqual(t *testing.T) {
	tests := []struct {
		x, y     float64
		expected bool
	}{
		{100, 100, true},
		{100, 100 + 1e-10, true},
		{100, 100 + 1e-8, false},
	}
	for _, test := range tests {
		got := Equal(test.x, test.y)
		if got != test.expected {
			t.Errorf("unexpected result. expected: %v, but got: %v", test.expected, got)
		}
	}
}
