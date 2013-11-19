package command

import (
	"github.com/bronze1man/kmg/console"
	"os"
)

type GoFmt struct {
}

func (command *GoFmt) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoFmt", Short: `format all golang code in a dir,same as "gofmt -w=true ."`}
}
func (command *GoFmt) Execute(context *console.Context) error {
	cmd := console.NewStdioCmd(context, "gofmt", "-w=true", ".")
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
