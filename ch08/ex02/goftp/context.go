package goftp

import (
	"net"
)

type Context struct {
	Command string
	Args    []string

	controlConn *ControlConn
}

func (c *Context) Send(code int, msg string) error {
	return c.controlConn.Send(code, msg)
}

func (c *Context) DataConn() (net.Conn, error) {
	return c.controlConn.DataConn()
}
