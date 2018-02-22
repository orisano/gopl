package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Clock struct {
	Location string
	Address  string
}

type Time struct {
	idx  int
	text string
}

func run() error {
	var clocks []*Clock
	for _, server := range os.Args[1:] {
		clock, err := ParseClock(server)
		if err != nil {
			return err
		}
		clocks = append(clocks, clock)
	}

	ch := make(chan Time, len(clocks))
	for i, clock := range clocks {
		conn, err := net.Dial("tcp", clock.Address)
		if err != nil {
			return err
		}
		defer conn.Close()
		go handleConn(conn, i, ch)
	}

	for _, clock := range clocks {
		fmt.Printf("%16s", clock.Location)
	}
	fmt.Println()

	times := make([]string, len(clocks))
	for t := range ch {
		times[t.idx] = t.text
		fmt.Print("\r")
		for _, time := range times {
			fmt.Printf("%16s", time)
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func ParseClock(server string) (*Clock, error) {
	tokens := strings.Split(server, "=")
	if len(tokens) != 2 {
		return nil, fmt.Errorf("invalid format: %v", server)
	}
	return &Clock{
		Location: tokens[0],
		Address:  tokens[1],
	}, nil
}

func handleConn(conn net.Conn, i int, ch chan<- Time) {
	b := bufio.NewScanner(conn)
	for b.Scan() {
		ch <- Time{i, b.Text()}
	}
}
