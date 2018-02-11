package goftp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/orisano/gopl/ch08/ex02/goftp/ftpcodes"
)

type ControlConn struct {
	logger           *log.Logger
	conn             net.Conn
	fs               FileSystem
	workingDirectory string

	dataPort *net.TCPAddr
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
		next := c.runCommand(cmd, args)
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

type byteParser struct {
	err error
}

func (p *byteParser) Parse(s string) byte {
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

func (p *byteParser) Err() error {
	return p.err
}

func (c *ControlConn) runCommand(cmd string, args []string) bool {
	switch strings.ToLower(cmd) {
	case "user":
		c.Send(ftpcodes.UserLoggedOn, "Welcome")
	case "quit":
		c.Send(ftpcodes.ServiceClosingTELNETConnection, "Good bye")
		return false
	case "syst":
		// https://www.iana.org/assignments/operating-system-names/operating-system-names.txt
		c.Send(ftpcodes.SystemStatus, "OSX system type")
	case "pwd":
		c.Send(ftpcodes.Entering, fmt.Sprintf("%q current working directory", c.workingDirectory))
	case "feat":
		c.Send(ftpcodes.SystemStatus, "No features")
	case "port":
		tokens := strings.Split(args[0], ",")
		bp := &byteParser{}
		h1 := bp.Parse(tokens[0])
		h2 := bp.Parse(tokens[1])
		h3 := bp.Parse(tokens[2])
		h4 := bp.Parse(tokens[3])
		p1 := bp.Parse(tokens[4])
		p2 := bp.Parse(tokens[5])
		c.dataPort = &net.TCPAddr{
			IP:   net.IPv4(h1, h2, h3, h4),
			Port: int(p1)*256 + int(p2),
		}
		CommandOK(c)
	case "type":
		if len(args) == 0 {
			CommandSyntaxError(c)
			return true
		}
		switch args[0] {
		case "A":
			if len(args) > 2 {
				CommandSyntaxError(c)
				return true
			}
			if len(args) == 1 {
				switch args[1] {
				case "N":
				case "T", "C":
					CommandNotImplementedForParameter(c)
					return true
				}
			}
		case "E", "I", "L":
			CommandNotImplementedForParameter(c)
			return true
		}
		CommandOK(c)
	case "stru":
		if len(args) != 1 {
			CommandSyntaxError(c)
			return true
		}
		switch args[0] {
		case "F":
		case "R", "P":
			CommandNotImplementedForParameter(c)
			return true
		}
	case "noop":
		CommandOK(c)
	case "mode":
		if len(args) != 1 {
			CommandSyntaxError(c)
			return true
		}
		switch args[0] {
		case "S":
		case "B", "C":
			CommandNotImplementedForParameter(c)
			return true
		}
	case "retr":
		conn, err := c.DataConn()
		if err != nil {
			c.logger.Print("connection error occurred: ", err)
			c.Send(ftpcodes.ConnectionTrouble, "Connection trouble")
			return true
		}
		defer conn.Close()

		f, err := c.fs.Get(args[0])
		if err != nil {
			c.logger.Print("failed to open file: ", err)
			c.Send(ftpcodes.FileUnavailable, "File unavailable")
			return true
		}
		defer f.Close()
		c.Send(ftpcodes.FileStatusOkay, "File status OK")
		io.Copy(conn, f)
		c.Send(ftpcodes.ClosingDataConnection, "Transfer completed")
	default:
		CommandNotImplemented(c)
	}
	return true
}
