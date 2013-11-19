package command

import (
	"github.com/bronze1man/kmg/console"
	"os"
)

type GoRun struct {
}

func (command *GoRun) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoRun", Short: "run some golang code(auto set current dir as GOPATH)"}
}
func (command *GoRun) Execute(context *console.Context) error {
	args := append([]string{"run"}, context.Args[2:]...)
	cmd := console.NewStdioCmd(context, "go", args...)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = console.SetCmdEnv(cmd, "GOPATH", wd)
	if err != nil {
		return err
	}
	return cmd.Run()
}
