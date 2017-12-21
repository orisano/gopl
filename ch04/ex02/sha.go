package main

import (
	"crypto/sha256"
	"flag"
	"hash"
	"io"
	"log"
	"os"

	"crypto/sha512"
	"encoding/hex"
)

func main() {
	t := flag.Int("type", 256, "sha size. support [256, 384, 512]")
	flag.Parse()

	var h hash.Hash
	switch *t {
	case 256:
		h = sha256.New()
	case 384:
		h = sha512.New384()
	case 512:
		h = sha512.New()
	default:
		log.Fatalf("unsupported type: %v", *t)
	}
	io.Copy(h, os.Stdin)
	io.WriteString(os.Stdout, hex.EncodeToString(h.Sum(nil)))
}
