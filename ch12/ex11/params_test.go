package params

import "testing"

func TestPack(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected string
	}{
		{
			in: &struct {
				A string
			}{
				"lower",
			},
			expected: "a=lower",
		},
		{
			in: &struct {
				one int
				two int
			}{
				1, 2,
			},
			expected: "one=1&two=2",
		},
		{
			in: &struct {
				b bool
			}{
				true,
			},
			expected: "b=true",
		},
		{
			in: &struct {
				q []string
			}{
				[]string{"hello", "world"},
			},
			expected: "q=hello&q=world",
		},
	}

	for _, test := range tests {
		if got := Pack(test.in); got != test.expected {
			t.Errorf("unexpected query string. expected: %v, but got: %v", test.expected, got)
		}
	}
}
