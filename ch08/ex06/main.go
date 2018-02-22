package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/orisano/gopl/ch05/links"
)

var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{}
	list, err := links.Extract(url)
	<-tokens
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	depth := flag.Int("depth", 0, "depth")
	flag.Parse()

	type Links struct {
		List  []string
		Depth int
	}

	type Link struct {
		URL   string
		Depth int
	}

	worklist := make(chan *Links)
	unseenLinks := make(chan *Link)

	go func() { worklist <- &Links{flag.Args(), 0} }()

	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link.URL)
				go func(depth int) {
					worklist <- &Links{foundLinks, depth + 1}
				}(link.Depth)
			}
		}()
	}

	seen := make(map[string]bool)
	for list := range worklist {
		if list.Depth > *depth {
			continue
		}
		for _, link := range list.List {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- &Link{link, list.Depth}
			}
		}
	}
}
