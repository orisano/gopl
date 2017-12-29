package main

import (
	"flag"
	"fmt"
)

var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":      {"discrete math"},
	"databases":            {"data structures"},
	"discrete math":        {"intro to programming"},
	"formal languages":     {"discrete math"},
	"networks":             {"operating system"},
	"operating system":     {"data structures", "computer organization"},
	"programming language": {"data structures", "computer organization"},
}

func main() {
	s := flag.String("s", "compilers", "start point")
	flag.Parse()

	dist := map[string]int{}
	dist[*s] = 0

	breadthFirst(func(item string) []string {
		for _, course := range prereqs[item] {
			if _, ok := dist[course]; ok {
				continue
			}
			dist[course] = dist[item] + 1
		}
		return prereqs[item]
	}, []string{*s})

	for course, d := range dist {
		fmt.Println(course, ":", d)
	}
}

func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}
