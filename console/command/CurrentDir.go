package command

import (
	"github.com/bronze1man/kmg/console"
	"os"
)

type CurrentDir struct {
}

func (command *CurrentDir) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "CurrentDir",
		Short: "get current dir(usefull in cygwin)",
	}
}

func (command *CurrentDir) Execute(context *console.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = context.Stdout.Write([]byte(wd))
	if err != nil {
		return err
	}
	return nil
}
