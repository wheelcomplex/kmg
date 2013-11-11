package buildCommand

import (
	"github.com/bronze1man/kmg/console"
	"os"
	"os/exec"
)

type FmtCommand struct {
}

func (command *FmtCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "fmt", Short: "format all golang code in a dir"}
}
func (command *FmtCommand) Execute(context *console.Context) error {
	cmd := console.NewStdioCmd(context, "gofmt", "-w=true", ".")
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd.Env = append(cmd.Env, "GOPATH="+wd)
	return cmd.Run()
}
