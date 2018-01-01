package main

import (
	"bufio"
)

type WordLineCounter struct {
	words, lines int
	terminated   bool
}

func (c *WordLineCounter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if b == '\n' {
			c.lines++
		}
	}
	for b := p; len(b) > 0; {
		advance, _, _ := bufio.ScanWords(b, true)
		b = b[advance:]
		c.words++
	}
	return len(p), nil
}

func (c *WordLineCounter) Words() int {
	return c.words
}

func (c *WordLineCounter) Lines() int {
	return c.lines
}
