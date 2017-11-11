package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	Echo(os.Stdout, os.Args)
}

func Echo(w io.Writer, args []string) {
	for i, arg := range args[1:] {
		fmt.Fprintln(w, i, arg)
	}
}
