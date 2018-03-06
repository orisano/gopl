package goftp2

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"path"
	"strings"
)

const (
	EOL = "\r\n"
)

type ControlConn struct {
	logger  *log.Logger
	conn    net.Conn
	fs      FileSystem
	connSrc ConnSource

	workingDir string

	handler CommandHandler
}

func (c *ControlConn) Send(code int, msg string) error {
	_, err := fmt.Fprint(c.conn, code, " ", msg, EOL)
	return err
}

func (c *ControlConn) DataConn() (net.Conn, error) {
	return c.connSrc.Conn()
}

func (c *ControlConn) ChangeSource(connSrc ConnSource) {
	c.closeSource()
	c.connSrc = connSrc
}

func (c *ControlConn) GetWD() string {
	return c.workingDir
}

func (c *ControlConn) ChangeWD(dir string) {
	c.workingDir = c.ResolvePath(dir)
}

func (c *ControlConn) ResolvePath(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = path.Join(c.workingDir, p)
	}
	return path.Clean(p)
}

func (c *ControlConn) FS() FileSystem {
	return c.fs
}

func (c *ControlConn) Logf(format string, v ...interface{}) {
	c.logger.Printf(format, v...)
}

func (c *ControlConn) Addr() *net.TCPAddr {
	return c.conn.LocalAddr().(*net.TCPAddr)
}

func (c *ControlConn) handle() {
	defer c.conn.Close()
	defer c.closeSource()

	c.Send(220, "GoFTP service")

	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		ctx := &Context{}

		line := scanner.Text()
		tokens := strings.SplitN(line, " ", 2)
		ctx.Command = tokens[0]
		if len(tokens) == 2 {
			ctx.Arg = strings.TrimSpace(tokens[1])
		}
		c.handler.Handle(ctx)
		if ctx.closed {
			break
		}
	}
}

func (c *ControlConn) closeSource() {
	if c.connSrc == nil {
		return
	}
	type closer interface {
		Close() error
	}
	if cl, ok := c.connSrc.(closer); ok {
		cl.Close()
	}
}
