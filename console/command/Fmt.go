package command

import (
	"github.com/bronze1man/kmg/console"
	"os"
)

type Fmt struct {
}

func (command *Fmt) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "Fmt", Short: `format all golang code in a dir,same as "gofmt -w=true ."`}
}
func (command *Fmt) Execute(context *console.Context) error {
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
