package goftp2

import (
	"net"
)

type ConnSource interface {
	Conn() (net.Conn, error)
}

type ActiveMode struct {
	src, dst *net.TCPAddr
}

func (m *ActiveMode) Conn() (net.Conn, error) {
	return net.DialTCP("tcp", m.src, m.dst)
}

func NewActiveMode(src, dst *net.TCPAddr) *ActiveMode {
	return &ActiveMode{src: src, dst: dst}
}

type PassiveMode struct {
	lis net.Listener
}

func (m *PassiveMode) Conn() (net.Conn, error) {
	return m.lis.Accept()
}

func (m *PassiveMode) Addr() net.TCPAddr {
	return *m.lis.Addr().(*net.TCPAddr)
}

func (m *PassiveMode) Close() error {
	return m.lis.Close()
}

func NewPassiveMode(address string) (*PassiveMode, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &PassiveMode{lis: lis}, nil
}
