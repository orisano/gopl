package goftp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/orisano/gopl/ch08/ex02/goftp/ftpcodes"
)

type ControlConn struct {
	logger           *log.Logger
	conn             net.Conn
	fs               FileSystem
	workingDirectory string

	dataPort *net.TCPAddr

	handler CommandHandler
}

func (c *ControlConn) Handle() {
	defer c.conn.Close()

	c.Send(ftpcodes.ServiceReadyForNewUser, "goftp service")

	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		line := scanner.Text()
		c.logger.Print(">>>> ", line)

		tokens := strings.Split(line, " ")
		cmd, args := tokens[0], tokens[1:]

		next := c.handler.Handle(&Context{
			Command: cmd,
			Args:    args,

			controlConn: c,
		})
		if !next {
			break
		}
	}
}

func (c *ControlConn) Send(code int, msg string) error {
	_, err := fmt.Fprint(c.conn, code, " ", msg, "\r\n")
	return err
}

func (c *ControlConn) DataConn() (net.Conn, error) {
	src := *c.conn.LocalAddr().(*net.TCPAddr)
	src.Port = src.Port - 1
	var dst net.TCPAddr
	if c.dataPort != nil {
		dst = *c.dataPort
	} else {
		dst = *c.conn.RemoteAddr().(*net.TCPAddr)
	}
	return net.DialTCP("tcp", &src, &dst)
}
