package goftp2

import (
	"fmt"
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

		CommandOK(ctx)
	})

	return mux
}
