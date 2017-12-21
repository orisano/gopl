package main

import (
	"crypto/sha256"
	"fmt"
)

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func ShaDiffBit(a, b []byte) int {
	diff := 0
	for i := range a {
		diff += int(pc[a[i]^b[i]])
	}
	return diff
}

func main() {
	foo := sha256.New().Sum([]byte("foo"))
	bar := sha256.New().Sum([]byte("bar"))
	fmt.Printf("foo: %v, bar: %v => %v\n", foo, bar, ShaDiffBit(foo, bar))
}
