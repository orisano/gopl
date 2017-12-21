package main

import (
	"bytes"
	"log"
	"os"
	"strings"
)

func main() {
	log.Println(comma(os.Args[1]))
}

func comma(s string) string {
	b := new(bytes.Buffer)

	if strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") {
		b.WriteString(s[:1])
		s = s[1:]
	}

	tail := ""
	if idx := strings.IndexByte(s, '.'); idx >= 0 {
		s, tail = s[:idx], s[idx:]
	}

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
	if len(tail) > 0 {
		b.WriteString(tail)
	}
	return b.String()
}
