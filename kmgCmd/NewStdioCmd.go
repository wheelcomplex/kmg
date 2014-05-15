package kmgCmd

import (
	"io"
	"os"
	"os/exec"
)

func NewOsStdioCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

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

var OsStdio = osStdio{}

type osStdio struct{}

func (io osStdio) GetStdin() io.ReadCloser {
	return os.Stdin
}

func (io osStdio) GetStdout() io.WriteCloser {
	return os.Stdout
}
func (io osStdio) GetStderr() io.WriteCloser {
	return os.Stderr
}
