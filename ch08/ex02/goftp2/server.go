package goftp2

import (
	"io/ioutil"
	"log"
	"net"
	"os"
)

type Server struct {
	FileSystem FileSystem
	Logger     *log.Logger
	Handler    CommandHandler
}

func (s *Server) Serve(l net.Listener) error {
	defer l.Close()
	handler := s.Handler
	if handler == nil {
		handler = DefaultCommandMux()
	}

	logger := s.Logger
	if logger == nil {
		w := ioutil.Discard
		if os.Getenv("GOFTP_DEBUG") != "" {
			w = os.Stderr
		}
		logger = log.New(w, "goftp:", log.LstdFlags)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			s.Logger.Printf("failed to accept: %v", err)
			continue
		}

		src := dataAddr(conn.LocalAddr().(*net.TCPAddr))
		dst := conn.RemoteAddr().(*net.TCPAddr)

		c := &ControlConn{
			logger: logger,
			fs:     s.FileSystem,

			conn:    conn,
			connSrc: NewActiveMode(src, dst),

			workingDir: "/",

			handler: handler,
		}
		go c.handle()
	}
}

func Serve(l net.Listener, fs FileSystem) error {
	s := &Server{
		FileSystem: fs,
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

func dataAddr(connAddr *net.TCPAddr) *net.TCPAddr {
	addr := *connAddr
	addr.Port = addr.Port - 1
	return &addr
}
