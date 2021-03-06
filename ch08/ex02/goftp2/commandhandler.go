package goftp2

import (
	"fmt"
	"io"
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
		CommandNotImplemented(ctx)
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

	// 125: Data connection already open; transfer starting.
	// 150: File status okay; about to open data connection.
	// 200: The requested action has been successfully completed.
	// 215:
	// 	NAME system type.
	// 	Where NAME is an official system name from the registry kept by IANA.
	// 220: Service ready for new user.
	// 221: Service closing control connection.
	// 226: Closing data connection.
	//   Requested file action successful (for example, file transfer or file abort).
	// 230: User logged in, proceed. Logged out if appropriate.
	// 250: Requested file action okay, completed.
	// 257: "PATHNAME" created.
	// 331: User name okay, need password.
	// 332: Need account for login.
	// 421:
	// 	Service not available, closing control connection.
	// 	This may be a reply to any command if the service knows it must shut down.
	// 425: Can't open data connection.
	// 426: Connection closed; transfer aborted.
	// 450:
	//  Requested file action not taken.
	//  File unavailable (e.g., file busy).
	// 451: Requested action aborted: local error in processing.
	// 500:
	// 	Syntax error, command unrecognized and the requested action did not take place.
	// 	This may include errors such as command line too long.
	// 501: Syntax error in parameters or arguments.
	// 502: Command not implemented.
	// 530: Not logged in.
	// 550: Requested action not taken. File unavailable (e.g., file not found, no access).

	// returns: 230, 530, 500, 501, 421, 331, 332
	mux.OnFunc("USER", func(ctx *Context) {
		ctx.Send(230, "Welcome")
	})

	// returns: 221, 500
	mux.OnFunc("QUIT", func(ctx *Context) {
		ctx.Send(221, "Bye")
		ctx.Close()
	})

	// returns: 215, 500, 501, 502, 421
	mux.OnFunc("SYST", func(ctx *Context) {
		ctx.Send(215, "OSX system type")
	})

	// returns: 257, 500, 501, 502, 421, 550
	mux.OnFunc("PWD", func(ctx *Context) {
		PathCreated(ctx, ctx.GetWD())
	})

	mux.OnFunc("FEAT", func(ctx *Context) {
		ctx.Send(211, "No Feature")
	})

	// returns: 200, 500, 501, 421, 530
	mux.OnFunc("PORT", func(ctx *Context) {
		tokens := strings.Split(ctx.Arg, ",")
		if len(tokens) != 6 {
			InvalidParametersOrArguments(ctx)
			return
		}

		p := &octetParser{}
		h1 := p.Parse(tokens[0])
		h2 := p.Parse(tokens[1])
		h3 := p.Parse(tokens[2])
		h4 := p.Parse(tokens[3])
		p1 := p.Parse(tokens[4])
		p2 := p.Parse(tokens[5])
		if err := p.Err(); err != nil {
			InvalidParametersOrArguments(ctx)
			return
		}

		src := dataAddr(ctx.Addr())
		dst := &net.TCPAddr{
			IP:   net.IPv4(h1, h2, h3, h4),
			Port: int(p1)*256 + int(p2),
		}
		ctx.ChangeSource(NewActiveMode(src, dst))
		CommandOK(ctx)
	})

	mux.OnFunc("NOOP", func(ctx *Context) {
		CommandOK(ctx)
	})

	// returns: 125, 150, 226, 250, 421, 425, 426, 450, 451, 500, 501, 530, 550
	mux.OnFunc("RETR", func(ctx *Context) {
		if len(ctx.Arg) == 0 {
			InvalidParametersOrArguments(ctx)
			return
		}
		p := ctx.ResolvePath(ctx.Arg)

		stat, err := ctx.FS().Stat(p)
		if err == ErrNotFound {
			ctx.Send(550, "Requested action not taken. File unavailable")
			return
		}
		if err != nil {
			ctx.Send(451, "Requested action aborted: local error in processing.")
			return
		}
		if stat.IsDir {
			ctx.Send(550, fmt.Sprintf("Requested action not taken. %q is directory", p))
			return
		}

		f, err := ctx.FS().Open(p)
		if err != nil {
			ctx.Send(451, "Requested action aborted: local error in processing.")
			return
		}
		defer f.Close()

		conn, err := ctx.DataConn()
		if err != nil {
			ctx.Send(425, "Can't open data connection")
			return
		}
		defer conn.Close()

		ctx.Send(125, "Data connection already open; transfer starting")
		if _, err := io.Copy(conn, f); err != nil {
			ctx.Send(426, "Connection closed; transfer aborted.")
			return
		}
		ctx.Send(226, "Closing data connection; Requested file action successful")
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
