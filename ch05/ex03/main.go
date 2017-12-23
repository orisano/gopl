package main

import (
	"io"
	"log"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	renderTextNode(os.Stdout, doc)
}

func renderTextNode(w io.Writer, n *html.Node) {
	if n == nil {
		return
	}
	if n.Type == html.TextNode {
		io.WriteString(w, n.Data)
	}
	if !(n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style")) {
		renderTextNode(w, n.FirstChild)
	}
	renderTextNode(w, n.NextSibling)
}
