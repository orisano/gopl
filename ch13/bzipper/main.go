package main

import (
	"io"
	"log"
	"os"

	"github.com/orisano/gopl/ch13/bzip"
)

func main() {
	w := bzip.NewWriter(os.Stdout)
	if _, err := io.Copy(w, os.Stdin); err != nil {
		log.Fatalf("bzipper: %v", err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("bzipper: %v", err)
	}
}
