package command

import (
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/console/kmgContext"
	"github.com/bronze1man/kmg/kmgCmd"
)

type GoRun struct {
}

func (command *GoRun) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoRun", Short: "run some golang code in current project"}
}
func (command *GoRun) Execute(context *console.Context) (err error) {
	args := append([]string{"run"}, context.Args[2:]...)
	cmd := console.NewStdioCmd(context, "go", args...)
	kmgc, err := kmgContext.FindFromWd()
	if err != nil {
		return
	}
	err = kmgCmd.SetCmdEnv(cmd, "GOPATH", kmgc.GOPATHToString())
	if err != nil {
		return err
	}
	return cmd.Run()
}
