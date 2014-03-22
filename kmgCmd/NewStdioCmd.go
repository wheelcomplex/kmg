package kmgCmd

import (
	"io"
	"os"
	"os/exec"
)

type Stdio interface {
	GetStdin() io.ReadCloser
	GetStdout() io.WriteCloser
	GetStderr() io.WriteCloser
}

func NewStdioCmd(stdio Stdio, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = stdio.GetStdin()
	cmd.Stdout = stdio.GetStdout()
	cmd.Stderr = stdio.GetStderr()
	return cmd
}

var DefaultStdio = defaultStdio{}

type defaultStdio struct{}

func (io defaultStdio) GetStdin() io.ReadCloser {
	return os.Stdin
}

func (io defaultStdio) GetStdout() io.WriteCloser {
	return os.Stdout
}
func (io defaultStdio) GetStderr() io.WriteCloser {
	return os.Stderr
}
