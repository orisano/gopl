package main

import (
	"reflect"
	"testing"
)

func TestParseClock(t *testing.T) {
	tests := []struct {
		server   string // Name=Address
		expected *Clock
	}{
		{
			server: "NewYork=localhost:8010",
			expected: &Clock{
				Location: "NewYork",
				Address:  "localhost:8010",
			},
		},
		{
			server: "Tokyo=localhost:8020",
			expected: &Clock{
				Location: "Tokyo",
				Address:  "localhost:8020",
			},
		},
		{
			server: "London=localhost:8030",
			expected: &Clock{
				Location: "London",
				Address:  "localhost:8030",
			},
		},
	}

	for _, test := range tests {
		got, err := ParseClock(test.server)
		if err != nil {
			t.Errorf("failed to parse: %v", err)
			continue
		}
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("unexpected clock. expected: %+#v, but got: %+#v", test.expected, got)
		}
	}
}
