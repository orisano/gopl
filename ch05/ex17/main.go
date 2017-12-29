package main

import "golang.org/x/net/html"

func ElementsByTagName(doc *html.Node, name ...string) []*html.Node {
	var nodes []*html.Node

	nameSet := toSet(name)
	forEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode {
			if _, ok := nameSet[n.Data]; ok {
				nodes = append(nodes, n)
			}
		}
	}, nil)

	return nodes
}

func toSet(ss []string) map[string]struct{} {
	x := make(map[string]struct{}, len(ss))
	for _, s := range ss {
		x[s] = struct{}{}
	}
	return x
}

func forEachNode(n *html.Node, pre, post func(*html.Node)) {
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
