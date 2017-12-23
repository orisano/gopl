package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

func getAttr(n *html.Node, attr string) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == attr {
			return a.Val, true
		}
	}
	return "", false
}

var linkMap = map[string]string{
	"a":      "href",
	"script": "src",
	"link":   "href",
	"img":    "src",
}

func visit(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}
	if n.Type == html.ElementNode {
		if attr, ok := linkMap[n.Data]; ok {
			if link, ok := getAttr(n, attr); ok {
				links = append(links, link)
			}
		}
	}
	links = visit(links, n.FirstChild)
	return visit(links, n.NextSibling)
}
