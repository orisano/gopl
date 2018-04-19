package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/orisano/gopl/ch13/bzip"
)

func main() {
	sha := sha1.New()
	w := bzip.NewWriter(sha)
	if _, err := io.Copy(w, os.Stdin); err != nil {
		log.Fatalf("bzipper: %v", err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("bzipper: %v", err)
	}
	fmt.Print(hex.EncodeToString(sha.Sum(nil)))
}
