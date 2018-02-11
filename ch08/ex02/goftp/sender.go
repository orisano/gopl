package goftp

import "github.com/orisano/gopl/ch08/ex02/goftp/ftpcodes"

type Sender interface {
	Send(code int, message string) error
}

func CommandNotImplemented(s Sender) error {
	return s.Send(ftpcodes.CommandNotImplemented, "Command not implemented")
}

func CommandOK(s Sender) error {
	return s.Send(ftpcodes.CommandOkay, "Command OK")
}

func CommandSyntaxError(s Sender) error {
	return s.Send(ftpcodes.CommandSyntaxError, "Command syntax error")
}

func CommandNotImplementedForParameter(s Sender) error {
	return s.Send(ftpcodes.CommandNotImplementedForParameter, "Command not implemented for that parameter")
}
