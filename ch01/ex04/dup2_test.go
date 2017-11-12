package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestUnique(t *testing.T) {
	ts := []struct {
		strs     []string
		expected []string
	}{
		{
			strs:     []string{"", "a", "b", "c", "d"},
			expected: []string{"", "a", "b", "c", "d"},
		},
		{
			strs:     []string{"a", "a", "a"},
			expected: []string{"a"},
		},
		{
			strs:     []string{"a", "b", "b", "a", "c", "a", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			strs:     []string{"b", "b", "a", "c", "a", "c"},
			expected: []string{"b", "a", "c"},
		},
	}
	for _, tc := range ts {
		got := Unique(tc.strs)
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("unexpected result. expected: %v, but got: %v", tc.expected, got)
		}
	}
}

func TestCountLines(t *testing.T) {
	ts := []struct {
		text     string
		name     string
		expected map[string][]string
	}{
		{
			text: "1\n2\n3\n4\n\n",
			name: "test",
			expected: map[string][]string{
				"1": {"test"},
				"2": {"test"},
				"3": {"test"},
				"4": {"test"},
				"":  {"test"},
			},
		},
		{
			text: "1\n2\n1\n1\n",
			name: "file.txt",
			expected: map[string][]string{
				"1": {"file.txt", "file.txt", "file.txt"},
				"2": {"file.txt"},
			},
		},
		{
			text: "1\n2\n1\n1",
			name: "file.txt",
			expected: map[string][]string{
				"1": {"file.txt", "file.txt", "file.txt"},
				"2": {"file.txt"},
			},
		},
	}
	for _, tc := range ts {
		c := make(map[string][]string)
		CountLines(tc.name, bytes.NewBufferString(tc.text), c)
		if got := c; !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("unexpected result. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
