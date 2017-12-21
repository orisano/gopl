package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	categories := map[string]func(rune) bool{
		"control": unicode.IsControl,
		"digit":   unicode.IsDigit,
		"space":   unicode.IsSpace,
		"symbol":  unicode.IsSymbol,
	}
	countsCategories := make(map[string]int)

	counts := make(map[rune]int)
	var utflen [utf8.UTFMax + 1]int
	invalid := 0

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprint(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
		for title, fn := range categories {
			if fn(r) {
				countsCategories[title]++
			}
		}
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}

	fmt.Printf("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}

	fmt.Printf("\ncategory\tcount\n")
	for title, cnt := range countsCategories {
		if cnt == 0 {
			continue
		}
		fmt.Printf("%16s\t%d\n", title, cnt)
	}

	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}
