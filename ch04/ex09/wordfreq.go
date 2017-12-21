package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)

	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanWords)
	for in.Scan() {
		text := in.Text()
		counts[text]++
	}
	fmt.Printf("text\tcount\n")
	for text, n := range counts {
		fmt.Printf("%16s\t%d\n", text, n)
	}
}
