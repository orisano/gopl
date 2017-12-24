package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		if err := outline(url); err != nil {
			log.Print(err)
		}
	}
}

func outline(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get request (%v): %v", url, err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse html (%v): %v", url, err)
	}
	forEachNode(doc, startElement, endElement)
	return nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

var depth int

func Indent(s, indent string) string {
	lines := strings.Split(s, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		result = append(result, indent+line)
	}
	return strings.Join(result, "\n")
}

func startElement(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		fmt.Printf("%*s<%s", depth*2, "", n.Data)
		for _, a := range n.Attr {
			fmt.Printf(" %s=%q", a.Key, a.Val)
		}
		if n.FirstChild == nil {
			fmt.Println("/>")
		} else {
			fmt.Println(">")
			depth++
		}
	case html.TextNode:
		s := strings.TrimSpace(n.Data)
		if len(s) > 0 {
			fmt.Println(Indent(s, strings.Repeat(" ", depth*2)))
		}
	case html.CommentNode:
		fmt.Printf("%*s<!-- %s -->\n", depth*2, "", n.Data)
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode && n.FirstChild != nil {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}
