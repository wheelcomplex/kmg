package command

import (
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/console/kmgContext"
	"github.com/bronze1man/kmg/kmgCmd"
	"path/filepath"
	"runtime"
)

type GoCrossCompileInit struct {
}

func (command *GoCrossCompileInit) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoCrossCompileInit", Short: "cross compile init target in current project"}
}
func (command *GoCrossCompileInit) Execute(context *console.Context) (err error) {
	kmgc, err := kmgContext.FindFromWd()
	if err != nil {
		return
	}
	GOROOT := kmgc.GOROOT
	if GOROOT == "" {
		return fmt.Errorf("you must set $GOROOT in environment to use GoCrossComplieInit")
	}
	var makeShellArgs []string
	var makeShellName string
	runCmdPath := filepath.Join(GOROOT, "src")
	if runtime.GOOS == "windows" {
		makeShellName = "cmd"
		makeShellArgs = []string{"/C", filepath.Join(GOROOT, "src", "make.bat"), "--no-clean"}
	} else {
		makeShellName = filepath.Join(GOROOT, "src", "make.bash")
		makeShellArgs = []string{"--no-clean"}
	}
	for _, target := range kmgc.CrossCompileTarget {
		cmd := console.NewStdioCmd(context, makeShellName, makeShellArgs...)
		kmgCmd.SetCmdEnv(cmd, "GOOS", target.GetGOOS())
		kmgCmd.SetCmdEnv(cmd, "GOARCH", target.GetGOARCH())
		cmd.Dir = runCmdPath
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return
}
