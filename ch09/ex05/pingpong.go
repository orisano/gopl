package main

import (
	"sync/atomic"
	"time"
	"fmt"
)

func main() {
	ping := make(chan struct{})
	pong := make(chan struct{})

	var count int64
	go func() {
		for {
			<-pong
			atomic.AddInt64(&count, 1)
			ping <- struct{}{}
		}
	} ()

	go func() {
		for {
			<-ping
			atomic.AddInt64(&count, 1)
			pong <- struct{}{}
		}
	} ()

	go func() {
		ping <- struct{}{}
	}()

	time.Sleep(1 * time.Second)
	c := atomic.LoadInt64(&count)
	fmt.Println(c)
}
