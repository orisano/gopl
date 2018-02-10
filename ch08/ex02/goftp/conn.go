package goftp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/orisano/gopl/ch08/ex02/goftp/ftpcodes"
)

type Conn struct {
	logger           *log.Logger
	conn             net.Conn
	fs               FileSystem
	workingDirectory string

	dataPort *net.TCPAddr
}

func (c *Conn) Handle() {
	defer c.conn.Close()

	c.writeReply(ftpcodes.ServiceReadyForNewUser)

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

func (c *Conn) runCommand(cmd string, args []string) bool {
	switch strings.ToLower(cmd) {
	case "user":
		c.writeReply(ftpcodes.UserLoggedOn)
	case "quit":
		c.writeReply(ftpcodes.ServiceClosingTELNETConnection)
		return false
	case "syst":
		c.writeReply(ftpcodes.SystemType)
	case "pwd":
		c.writeReply(ftpcodes.Entering)
	case "feat":
		c.writeReply(ftpcodes.SystemStatus)
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
		c.writeReply(ftpcodes.CommandOkay)
	default:
		c.writeReply(ftpcodes.CommandNotImplemented)
	}
	return true
}

func (c *Conn) writeReply(code int) error {
	write := func(msg string) error {
		_, err := fmt.Fprint(c.conn, code, " ", msg, "\r\n")
		return err
	}
	switch code {
	case ftpcodes.ServiceReadyForNewUser:
		return write("Service ready")
	case ftpcodes.CommandNotImplemented:
		return write("Command not implemented")
	case ftpcodes.UserLoggedOn:
		return write("User logged on, proceed")
	case ftpcodes.ServiceClosingTELNETConnection:
		return write("Service closing TELNET connection (logged off if appropriate)")
	case ftpcodes.SystemType:
		// https://www.iana.org/assignments/operating-system-names/operating-system-names.txt
		return write("OSX system type")
	case ftpcodes.Entering:
		return write(fmt.Sprintf("%q currenct working directory", c.workingDirectory))
	case ftpcodes.SystemStatus:
		return write("No features")
	case ftpcodes.CommandOkay:
		return write("Command OK")
	default:
		panic(fmt.Sprintf("unknown code: %v", code))
	}
}
