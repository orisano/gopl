package goftp2

import "fmt"

type Sender interface {
	Send(code int, message string) error
}

func CommandOK(s Sender) error {
	return s.Send(200, "Command Okay")
}

func PathCreated(s Sender, path string) error {
	return s.Send(257, fmt.Sprintf("%q created", path))
}

// This may be a reply to any command if the service knows it must shut down.
func ServiceNotAvailable(s Sender) error {
	return s.Send(421, "Service not available, closing control connection")
}

// This may include errors such as command line too long.
func CommandUnrecognized(s Sender) error {
	return s.Send(500, "Syntax error, command unrecognized and the requested action did not take place")
}

func InvalidParametersOrArguments(s Sender) error {
	return s.Send(501, "Syntax error in parameters or arguments.")
}

func CommandNotImplemented(s Sender) error {
	return s.Send(502, "Command not implemented")
}

func CommandNotImplementedForParameter(s Sender) error {
	return s.Send(504, "Command not implemented for that parameter")
}
