package main

import (
	"fmt"
	"log"
)

var prereqs = map[string]map[string]bool{
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
}

func main() {
	order, ok := topoSort(prereqs)
	if !ok {
		log.Fatal("loop detected")
	}
	for i, course := range order {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string]map[string]bool) ([]string, bool) {
	const (
		Yet = iota
		Temporary
		Visited
	)

	var order []string
	var loopDetected bool
	seen := make(map[string]int)
	var visitAll func(items map[string]bool)

	visitAll = func(items map[string]bool) {
		for item := range items {
			switch seen[item] {
			case Yet:
				seen[item] = Temporary
				visitAll(m[item])
				order = append(order, item)
				seen[item] = Visited
			case Temporary:
				loopDetected = true
			}
		}
	}
	keys := make(map[string]bool)
	for key := range m {
		keys[key] = true
	}
	visitAll(keys)
	return order, !loopDetected
}
