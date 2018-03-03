package goftp2

type Context struct {
	*ControlConn

	Command string
	Arg     string

	closed bool
}

func (c *Context) Close() {
	c.closed = true
}
