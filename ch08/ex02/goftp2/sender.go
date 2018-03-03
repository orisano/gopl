package goftp2

type Sender interface {
	Send(code int, message string) error
}

func CommandNotImplemented(s Sender) error {
	return s.Send(502, "Command not implemented")
}

func CommandOK(s Sender) error {
	return s.Send(200, "Command OK")
}

func CommandSyntaxError(s Sender) error {
	return s.Send(500, "Command syntax error")
}

func CommandNotImplementedForParameter(s Sender) error {
	return s.Send(504, "Command not implemented for that parameter")
}
