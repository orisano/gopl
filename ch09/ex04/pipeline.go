package main

import (
	"time"
	"fmt"
)

func main() {
	var to []chan struct{}
	to = append(to, make(chan struct{}))

	id := 0
	for {
		to = append(to, make(chan struct{}))
		go func(recv, send chan struct{}) {
			for {
				<-recv
				send <- struct{}{}
			}
		} (to[id], to[id + 1])
		id++

		t := time.Now()
		to[0] <- struct{}{}
		<-to[id]
		d := time.Now().Sub(t)

		fmt.Printf("pipelines: %d, duration: %v\n", id, d)
	}
}