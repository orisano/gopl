package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	Echo(os.Stdout, os.Args)
}

func Echo(w io.Writer, args []string) {
	fmt.Fprintln(w, strings.Join(args, " "))
}
