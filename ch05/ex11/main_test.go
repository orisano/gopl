package main

import (
	"testing"
)

func Test_topoSort(t *testing.T) {
	ts := []struct {
		graph        map[string]map[string]bool
		loopContains bool
	}{
		{
			graph: map[string]map[string]bool{
				"algorithms": {"data structures": true},
				"calculus":   {"linear algebra": true},
				"compilers": {
					"data structures":       true,
					"formal languages":      true,
					"computer organization": true,
				},
				"data structures":      {"discrete math": true},
				"databases":            {"data structures": true},
				"discrete math":        {"intro to programming": true},
				"formal languages":     {"discrete math": true},
				"networks":             {"operating system": true},
				"operating system":     {"data structures": true, "computer organization": true},
				"programming language": {"data structures": true, "computer organization": true},
			},
			loopContains: false,
		},
		{
			graph: map[string]map[string]bool{
				"algorithms":     {"data structures": true},
				"calculus":       {"linear algebra": true},
				"linear algebra": {"calculus": true},
				"compilers": {
					"data structures":       true,
					"formal languages":      true,
					"computer organization": true,
				},
				"data structures":      {"discrete math": true},
				"databases":            {"data structures": true},
				"discrete math":        {"intro to programming": true},
				"formal languages":     {"discrete math": true},
				"networks":             {"operating system": true},
				"operating system":     {"data structures": true, "computer organization": true},
				"programming language": {"data structures": true, "computer organization": true},
			},
			loopContains: true,
		},
	}

	for _, tc := range ts {
		t.Run("Case", func(t *testing.T) {
			order, ok := topoSort(tc.graph)
			if tc.loopContains == ok {
				t.Fatalf("unexpected loop. expected: %v, but got: %v", tc.loopContains, ok)
			}
			if ok {
				ord := make(map[string]int)
				for i, course := range order {
					ord[course] = i
				}

				for course, requires := range tc.graph {
					for require := range requires {
						if ord[course] < ord[require] {
							t.Errorf("failed to topoSort. must be ord[%q] < ord[%q], but got (ord[%q]=%d, ord[%q]=%d)",
								require, course, require, ord[require], course, ord[course])
						}
					}
				}
			}
		})
	}
}
