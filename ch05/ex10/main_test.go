package main

import (
	"testing"
)

func Test_topoSort(t *testing.T) {
	t.Run("IsTopoOrdered", func(t *testing.T) {
		ord := make(map[string]int)
		for i, course := range topoSort(prereqs) {
			ord[course] = i
		}

		for course, requires := range prereqs {
			for require := range requires {
				if ord[course] < ord[require] {
					t.Errorf("failed to topoSort. must be ord[%q] < ord[%q], but got (ord[%q]=%d, ord[%q]=%d)",
						require, course, require, ord[require], course, ord[course])
				}
			}
		}
	})
}
