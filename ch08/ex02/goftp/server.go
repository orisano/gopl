package goftp

import (
	"io/ioutil"
	"log"
	"net"
	"os"
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
		c := &ControlConn{
			logger:           s.Logger,
			conn:             conn,
			fs:               s.FileSystem,
			workingDirectory: "/",
		}
		go c.Handle()
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
