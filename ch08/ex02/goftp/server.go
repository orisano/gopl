package goftp

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/orisano/gopl/ch08/ex02/goftp/ftpcodes"
)

type Server struct {
	FileSystem FileSystem
	Logger     *log.Logger
}

func (s *Server) Serve(l net.Listener) error {
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			s.Logger.Printf("failed to accept: %v", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	writeReply(conn, ftpcodes.ServiceReadyForNewUser)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		s.Logger.Print(">>>> ", line)

		tokens := strings.Split(line, " ")
		cmd, args := tokens[0], tokens[1:]
		s.runCommand(conn, cmd, args)
	}
}

func (s *Server) runCommand(w io.Writer, cmd string, args []string) {
	switch strings.ToLower(cmd) {
	case "user":
		writeReply(w, ftpcodes.UserLoggedOn)
	default:
		writeReply(w, ftpcodes.CommandNotImplemented)
	}
}

func Serve(l net.Listener, fs FileSystem) error {
	w := ioutil.Discard
	if os.Getenv("GOFTP_DEBUG") != "" {
		w = os.Stderr
	}
	s := &Server{
		FileSystem: fs,
		Logger:     log.New(w, "goftp:", log.LstdFlags),
	}
	return s.Serve(l)
}

func ListenAndServe(addr string, fs FileSystem) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return Serve(l, fs)
}

func writeReply(w io.Writer, code int) error {
	write := func(msg string) error {
		_, err := fmt.Fprint(w, code, " ", msg, "\r\n")
		return err
	}
	switch code {
	case ftpcodes.ServiceReadyForNewUser:
		return write("Service ready")
	case ftpcodes.CommandNotImplemented:
		return write("Command not implemented")
	case ftpcodes.UserLoggedOn:
		return write("User logged on, proceed")
	default:
		panic(fmt.Sprintf("unknown code: %v", code))
	}
}
