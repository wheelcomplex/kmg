package console

import (
	"os"
	"os/exec"
)

func NewStdioCmd(context *Context, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = context.Stdin
	cmd.Stdout = context.Stdout
	cmd.Stderr = context.Stderr
	return cmd
}

func SetCmdEnv(cmd *exec.Cmd, key string, value string) error {
	if len(cmd.Env) == 0 {
		cmd.Env = os.Environ()
	}
	env, err := NewEnvFromArray(cmd.Env)
	if err != nil {
		return err
	}
	env.Values[key] = value
	cmd.Env = env.ToArray()
	return nil
}
