package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	count := map[string]int{}
	countTag(count, doc)

	fmt.Printf("tag     \tcount\n")
	for tag, cnt := range count {
		fmt.Printf("%-8s\t%d\n", tag, cnt)
	}
}

func countTag(count map[string]int, n *html.Node) {
	if n == nil {
		return
	}
	if n.Type == html.ElementNode {
		count[n.Data]++
	}
	countTag(count, n.FirstChild)
	countTag(count, n.NextSibling)
}
