package goftp

import (
	"fmt"
	"io"
	"net"
	"path"
	"strconv"
	"strings"

	"github.com/orisano/gopl/ch08/ex02/goftp/ftpcodes"
)

type CommandHandler interface {
	Handle(ctx *Context) bool
}

type CommandHandlerFunc func(ctx *Context) bool

func (f CommandHandlerFunc) Handle(ctx *Context) bool {
	return f(ctx)
}

type CommandMux struct {
	handlers map[string]CommandHandler
}

func (m *CommandMux) Handle(ctx *Context) bool {
	cmd := strings.ToLower(ctx.Command)
	h, ok := m.handlers[cmd]
	if !ok {
		CommandNotImplemented(ctx)
		return true
	}
	return h.Handle(ctx)
}

func (m *CommandMux) On(cmd string, handler CommandHandler) {
	m.handlers[strings.ToLower(cmd)] = handler
}

func (m *CommandMux) OnFunc(cmd string, fn func(ctx *Context) bool) {
	m.handlers[strings.ToLower(cmd)] = CommandHandlerFunc(fn)
}

func NewCommandMux() *CommandMux {
	return &CommandMux{
		handlers: make(map[string]CommandHandler),
	}
}

func DefaultCommandMux() *CommandMux {
	mux := NewCommandMux()

	mux.OnFunc("USER", func(ctx *Context) bool {
		ctx.Send(ftpcodes.UserLoggedOn, "Welcome")
		return true
	})

	mux.OnFunc("QUIT", func(ctx *Context) bool {
		ctx.Send(ftpcodes.ServiceClosingTELNETConnection, "Good bye")
		return false
	})

	mux.OnFunc("SYST", func(ctx *Context) bool {
		ctx.Send(ftpcodes.SystemType, "OSX system type")
		return true
	})

	mux.OnFunc("PWD", func(ctx *Context) bool {
		ctx.Send(ftpcodes.Entering, fmt.Sprintf("%q current working directory", ctx.controlConn.GetWD()))
		return true
	})

	mux.OnFunc("FEAT", func(ctx *Context) bool {
		ctx.Send(ftpcodes.SystemStatus, "No Feature")
		return true
	})

	mux.OnFunc("PORT", func(ctx *Context) bool {
		tokens := strings.Split(ctx.Args[0], ",")
		p := &int8Parser{}
		h1 := p.Parse(tokens[0])
		h2 := p.Parse(tokens[1])
		h3 := p.Parse(tokens[2])
		h4 := p.Parse(tokens[3])
		p1 := p.Parse(tokens[4])
		p2 := p.Parse(tokens[5])
		ctx.controlConn.dataPort = &net.TCPAddr{
			IP:   net.IPv4(h1, h2, h3, h4),
			Port: int(p1)*256 + int(p2),
		}
		CommandOK(ctx)
		return true
	})

	mux.OnFunc("TYPE", func(ctx *Context) bool {
		if len(ctx.Args) == 0 {
			CommandSyntaxError(ctx)
			return true
		}
		switch ctx.Args[0] {
		case "A":
			if len(ctx.Args) > 2 {
				CommandSyntaxError(ctx)
				return true
			}
			if len(ctx.Args) == 1 {
				switch ctx.Args[1] {
				case "N":
				case "T", "C":
					CommandNotImplementedForParameter(ctx)
					return true
				}
			}
		case "E", "I", "L":
			CommandNotImplementedForParameter(ctx)
			return true
		}
		CommandOK(ctx)
		return true
	})

	mux.OnFunc("STRU", func(ctx *Context) bool {
		if len(ctx.Args) != 1 {
			CommandSyntaxError(ctx)
			return true
		}
		switch ctx.Args[0] {
		case "F":
		case "R", "P":
			CommandNotImplementedForParameter(ctx)
			return true
		}
		CommandOK(ctx)
		return true
	})

	mux.OnFunc("NOOP", func(ctx *Context) bool {
		CommandOK(ctx)
		return true
	})

	mux.OnFunc("MODE", func(ctx *Context) bool {
		if len(ctx.Args) != 1 {
			CommandSyntaxError(ctx)
			return true
		}
		switch ctx.Args[0] {
		case "S":
		case "B", "C":
			CommandNotImplementedForParameter(ctx)
			return true
		}
		CommandOK(ctx)
		return true
	})

	mux.OnFunc("RETR", func(ctx *Context) bool {
		conn, err := ctx.DataConn()
		if err != nil {
			ctx.Logf("connection error occurred: %v", err)
			ctx.Send(ftpcodes.ConnectionTrouble, "Connection trouble")
			return true
		}
		defer conn.Close()

		f, err := ctx.controlConn.fs.Get(ctx.Args[0])
		if err != nil {
			ctx.Logf("failed to open file: %v", err)
			ctx.Send(ftpcodes.FileUnavailable, "File unavailable")
			return true
		}
		defer f.Close()

		ctx.Send(ftpcodes.FileStatusOkay, "File status OK")
		io.Copy(conn, f)
		ctx.Send(ftpcodes.ClosingDataConnection, "Transfer completed")
		return true
	})

	mux.OnFunc("STOR", func(ctx *Context) bool {
		conn, err := ctx.DataConn()
		if err != nil {
			ctx.Logf("connection error occurred: %v", err)
			ctx.Send(ftpcodes.ConnectionTrouble, "Connection trouble")
			return true
		}
		defer conn.Close()

		f, err := ctx.controlConn.fs.Create(ctx.Args[0])
		if err != nil {
			ctx.controlConn.logger.Print("failed to create file: ", err)
			ctx.Send(ftpcodes.FileUnavailable, "File unavailable")
			return true
		}
		defer f.Close()

		ctx.Send(ftpcodes.FileStatusOkay, "File status OK")
		io.Copy(f, conn)
		ctx.Send(ftpcodes.ClosingDataConnection, "Store completed")
		return true
	})

	mux.OnFunc("LIST", func(ctx *Context) bool {
		conn, err := ctx.DataConn()
		if err != nil {
			ctx.Logf("connection error occurred: %v", err)
			ctx.Send(ftpcodes.ConnectionTrouble, "Connection trouble")
			return true
		}
		defer conn.Close()

		p := ctx.controlConn.wd
		if len(ctx.Args) == 1 {
			if strings.HasPrefix(ctx.Args[0], "/") {
				p = ctx.Args[0]
			} else {
				p = path.Join(p, ctx.Args[0])
			}
		}

		list, err := ctx.controlConn.fs.LS(p)
		if err != nil {
			ctx.Logf("failed to get list dir: %v", err)
			ctx.Send(ftpcodes.LocalErrorInProcessing, "Requested action aborted. Local error in processing.")
			return true
		}

		ctx.Send(ftpcodes.FileStatusOkay, "File status OK")
		for _, x := range list {
			fmt.Fprint(conn, x, "\r\n")
		}
		ctx.Send(ftpcodes.ClosingDataConnection, "LIST completed")

		return true
	})

	mux.OnFunc("NLST", func(ctx *Context) bool {
		conn, err := ctx.DataConn()
		if err != nil {
			ctx.Logf("connection error occurred: %v", err)
			ctx.Send(ftpcodes.ConnectionTrouble, "Connection trouble")
			return true
		}
		defer conn.Close()

		p := ctx.controlConn.wd
		if len(ctx.Args) == 1 {
			if strings.HasPrefix(ctx.Args[0], "/") {
				p = ctx.Args[0]
			} else {
				p = path.Join(p, ctx.Args[0])
			}
		}

		list, err := ctx.controlConn.fs.List(p)
		if err != nil {
			ctx.Logf("failed to get list dir: %v", err)
			ctx.Send(ftpcodes.LocalErrorInProcessing, "Requested action aborted. Local error in processing.")
			return true
		}

		ctx.Send(ftpcodes.FileStatusOkay, "File status OK")
		for _, x := range list {
			fmt.Fprint(conn, x, "\r\n")
		}
		ctx.Send(ftpcodes.ClosingDataConnection, "NLST completed")

		return true
	})

	mux.OnFunc("EPRT", func(ctx *Context) bool {
		tokens := strings.Split(ctx.Args[0], "|")
		netProto := tokens[1]
		netAddr := tokens[2]
		portStr := tokens[3]

		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			ctx.Logf("invalid port: %v", err)
			CommandSyntaxError(ctx)
			return true
		}

		var network string
		switch netProto {
		case "1": // IPv4
			network = "ip4"
		case "2": // IPv6
			network = "ip6"
		default:
			panic("unknown protocol")
		}

		ip, err := net.ResolveIPAddr(network, netAddr)
		if err != nil {
			ctx.Logf("failed to resolve ip: %v", err)
			CommandSyntaxError(ctx)
			return true
		}
		ctx.controlConn.dataPort = &net.TCPAddr{
			IP:   ip.IP,
			Port: int(port),
		}
		CommandOK(ctx)
		return true
	})

	mux.OnFunc("EPSV", func(ctx *Context) bool {
		if ctx.controlConn.passive != nil {
			ctx.controlConn.passive.Close()
			ctx.controlConn.passive = nil
		}
		lis, err := net.Listen("tcp", ":0")
		if err != nil {
			ctx.Logf("failed to listen: %v", err)
			ctx.Send(ftpcodes.LocalErrorInProcessing, "Failed to listen")
			return true
		}
		ctx.controlConn.passive = lis

		ctx.Send(ftpcodes.ExtenedPassiveModeEntered, fmt.Sprintf("Entering Extended Passive Mode (|||%v|)", lis.Addr().(*net.TCPAddr).Port))
		return true
	})

	mux.OnFunc("CWD", func(ctx *Context) bool {
		ctx.controlConn.ChangeWD(ctx.Args[0])
		ctx.Send(ftpcodes.RequestedFileActionOkey, "OK")
		return true
	})

	return mux
}

type int8Parser struct {
	err error
}

func (p *int8Parser) Parse(s string) byte {
	if p.err != nil {
		return 0
	}
	x, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		p.err = err
		return 0
	}
	return byte(x)
}

func (p *int8Parser) Err() error {
	return p.err
}
