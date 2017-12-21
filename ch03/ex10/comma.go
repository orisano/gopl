package main

import (
	"bytes"
	"log"
	"os"
)

func main() {
	log.Println(comma(os.Args[1]))
}

func comma(s string) string {
	b := new(bytes.Buffer)
	for len(s) > 3 {
		size := len(s) % 3
		if size == 0 {
			size = 3
		}
		b.WriteString(s[:size])
		b.WriteString(",")
		s = s[size:]
	}
	b.WriteString(s)
	return b.String()
}
