package goftp2

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type CommandHandler interface {
	Handle(ctx *Context)
}

type CommandHandlerFunc func(ctx *Context)

func (f CommandHandlerFunc) Handle(ctx *Context) { f(ctx) }

type CommandMux struct {
	handlers map[string]CommandHandler
}

func (m *CommandMux) Handle(ctx *Context) {
	cmd := strings.ToLower(ctx.Command)
	h, ok := m.handlers[cmd]
	if !ok {
		return
	}
	h.Handle(ctx)
}

func (m *CommandMux) On(cmd string, handler CommandHandler) {
	m.handlers[strings.ToLower(cmd)] = handler
}

func (m *CommandMux) OnFunc(cmd string, fn func(ctx *Context)) {
	m.On(cmd, CommandHandlerFunc(fn))
}

func NewCommandMux() *CommandMux {
	return &CommandMux{
		handlers: make(map[string]CommandHandler),
	}
}

func DefaultCommandMux() *CommandMux {
	mux := NewCommandMux()

	mux.OnFunc("USER", func(ctx *Context) {
		ctx.Send(230, "Welcome")
	})

	mux.OnFunc("QUIT", func(ctx *Context) {
		ctx.Send(221, "Bye")
		ctx.Close()
	})

	mux.OnFunc("SYST", func(ctx *Context) {
		ctx.Send(215, "OSX system type")
	})

	mux.OnFunc("PWD", func(ctx *Context) {
		ctx.Send(227, fmt.Sprintf("%q current working directory", ctx.GetWD()))
	})

	mux.OnFunc("FEAT", func(ctx *Context) {
		ctx.Send(211, "No Feature")
	})

	mux.OnFunc("PORT", func(ctx *Context) {
		tokens := strings.Split(ctx.Arg, ",")
		p := &octetParser{}
		h1 := p.Parse(tokens[0])
		h2 := p.Parse(tokens[1])
		h3 := p.Parse(tokens[2])
		h4 := p.Parse(tokens[3])
		p1 := p.Parse(tokens[4])
		p2 := p.Parse(tokens[5])

		src := dataAddr(ctx.Addr())
		dst := &net.TCPAddr{
			IP:   net.IPv4(h1, h2, h3, h4),
			Port: int(p1)*256 + int(p2),
		}
		ctx.ChangeSource(NewActiveMode(src, dst))
		CommandOK(ctx)
	})

	return mux
}

type octetParser struct {
	err error
}

func (p *octetParser) Parse(s string) byte {
	if p.err != nil {
		return 0
	}
	x, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		p.err = err
		return 0
	}
	return byte(x)
}

func (p *octetParser) Err() error {
	return p.err
}
