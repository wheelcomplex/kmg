package buildCommand

import (
	"github.com/bronze1man/kmg/console"
	"os"
	"os/exec"
)

type RunCommand struct {
}

func (command *RunCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "run", Short: "run some golang code(auto set current dir as GOPATH)"}
}
func (command *RunCommand) Execute(context *console.Context) error {
	args := append([]string{"run"}, context.Args[2:]...)
	cmd := console.NewStdioCmd(context, "go", args...)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd.Env = append(cmd.Env, "GOPATH="+wd)
	return cmd.Run()
}
