package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var DefaultSource io.Reader = os.Stdin

func main() {
	counts, err := Dup(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	for line, files := range counts {
		fmt.Printf("%d\t%s\n", len(files), line)
		fmt.Println(strings.Join(Unique(files), ","))
	}
}

func CountLines(name string, r io.Reader, counts map[string][]string) {
	input := bufio.NewScanner(r)
	for input.Scan() {
		text := input.Text()
		counts[text] = append(counts[text], name)
	}
}

func Unique(strs []string) []string {
	exists := make(map[string]bool)
	u := make([]string, 0, len(strs))
	for _, str := range strs {
		if !exists[str] {
			exists[str] = true
			u = append(u, str)
		}
	}
	return u
}

func Dup(files []string) (map[string][]string, error) {
	counts := make(map[string][]string)

	if len(files) == 0 {
		CountLines("<stdin>", DefaultSource, counts)
	} else {
		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				return nil, fmt.Errorf("dup2: %v", err)
			}
			CountLines(file, f, counts)
			f.Close()
		}
	}

	for line, fs := range counts {
		if len(fs) == 1 {
			delete(counts, line)
		}
	}
	return counts, nil
}
