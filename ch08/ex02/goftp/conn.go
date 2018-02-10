package goftp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/orisano/gopl/ch08/ex02/goftp/ftpcodes"
)

type Conn struct {
	logger           *log.Logger
	conn             net.Conn
	fs               FileSystem
	workingDirectory string
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
	default:
		panic(fmt.Sprintf("unknown code: %v", code))
	}
}
