package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/orisano/gopl/ch08/ex10/links"
)

var tokens = make(chan struct{}, 20)

func crawl(ctx context.Context, url string) []string {
	fmt.Println(url)
	tokens <- struct{}{}
	list, err := links.Extract(ctx, url)
	<-tokens
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string)
	unseenLinks := make(chan string)

	go func() { worklist <- os.Args[1:] }()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		os.Stdin.Read(make([]byte, 1))
		cancel()
	}()

	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(ctx, link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}
