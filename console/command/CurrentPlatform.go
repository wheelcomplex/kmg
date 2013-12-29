package command

import (
	"github.com/bronze1man/kmg/console"
	"runtime"
)

type CurrentPlatform struct {
}

func (command *CurrentPlatform) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "CurrentPlatform",
		Short: "get current platform(from this binary)",
	}
}

func (command *CurrentPlatform) Execute(context *console.Context) (err error) {
	_, err = context.Stdout.Write([]byte(runtime.GOOS + "_" + runtime.GOARCH))
	return
}
