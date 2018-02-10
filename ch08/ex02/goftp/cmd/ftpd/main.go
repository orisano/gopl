package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/orisano/gopl/ch08/ex02/goftp"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get working directory: ", err)
	}

	d := flag.String("d", wd, "ftp server mount directory")
	port := flag.Int("port", 21, "ftp server control port")
	dataPort := flag.Int("data-port", -1, "ftp server data port")

	flag.Parse()

	if *dataPort < 0 {
		*dataPort = *port - 1
	}

	log.Fatal(goftp.ListenAndServe(fmt.Sprint(":", *port), &goftp.RawFileSystem{Root: *d}))
}
