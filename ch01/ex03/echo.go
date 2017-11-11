package main

import (
	"fmt"
	"io"
	"strings"
)

func NaiveEcho(w io.Writer, args []string) {
	s, sep := "", ""
	for _, arg := range args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Fprintln(w, s)
}

func FastEcho(w io.Writer, args []string) {
	fmt.Fprintln(w, strings.Join(args[1:], " "))
}
