package console

import (
	"os/exec"
)

func NewStdioCmd(context *Context, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = context.Stdin
	cmd.Stdout = context.Stdout
	cmd.Stderr = context.Stderr
	return cmd
}
