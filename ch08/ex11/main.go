package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	b := mirroredQuery(os.Args[1:])
	os.Stdout.Write(b)
}

func mirroredQuery(urls []string) []byte {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	responses := make(chan []byte)
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			b, err := request(ctx, u)
			if err != nil {
				log.Print(err)
				return
			}
			responses <- b
		}(url)
	}
	r := <-responses
	cancel()
	wg.Wait()
	return r
}

func request(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}
