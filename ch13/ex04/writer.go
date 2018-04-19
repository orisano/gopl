package bzip

import (
	"io"
	"os/exec"
)

type writer struct {
	cmd *exec.Cmd
	wc  io.WriteCloser
}

func NewWriter(out io.Writer) io.WriteCloser {
	cmd := exec.Command("/usr/bin/bzip2")
	cmd.Stdout = out
	wc, _ := cmd.StdinPipe()
	cmd.Start()
	return &writer{cmd, wc}
}

func (w *writer) Write(p []byte) (int, error) {
	return w.wc.Write(p)
}

func (w *writer) Close() error {
	w.wc.Close()
	return w.cmd.Wait()
}
