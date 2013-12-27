package command

import (
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/console/kmgContext"
	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/fsnotify"
	"os"
	"os/exec"
	"path/filepath"
)

//Warning!! do not use go run to build executes...
// you can kill go run xxx, you can not kill main.exe xxx
//TODO cygwin ctrl+c not exit children processes
//TODO wrap restart cmd stuff
//TODO wrap lastHappendTime watch stuff
type GoWatch struct {
	context        *console.Context
	GOPATH         string
	wd             string
	watcher        *fsnotify.Watcher
	cmd            *exec.Cmd
	mainFilePath   string
	targetFilePath string
	isDebug        bool //more output
}

func (command *GoWatch) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoWatch",
		Short: "watch current project and rebuild and restart app when some file changed",
	}
}

func (command *GoWatch) Execute(context *console.Context) (err error) {
	command.isDebug = false
	command.context = context
	if len(context.Args) != 3 {
		return errors.Sprintf("usage: %s watch [packages]", context.ExecutionName)
	}
	command.mainFilePath = context.Args[2]
	kmgc, err := kmgContext.FindFromWd()
	if err != nil {
		return
	}
	command.GOPATH = kmgc.GOPATHToString()
	command.wd = kmgc.ProjectPath
	runner, err := fsnotify.NewRunner(10000)
	if err != nil {
		return err
	}
	runner.Watcher.ErrorHandler = func(err error) {
		fmt.Println("watcher.Error error: ", err)
	}
	runner.Watcher.IsIgnorePath = func(path string) bool {
		if fsnotify.DefaultIsIgnorePath(path) {
			return true
		}
		if path == command.targetFilePath {
			return true
		}
		return false
	}
	runner.Watcher.WatchRecursion(command.wd)
	// wait forever
	runner.Run(func() {
		err := command.restart()
		if err != nil {
			fmt.Println("command.restart() error: ", err)
		}
	})
	return nil
}
func (command *GoWatch) restart() error {
	//kill old process
	err := command.stop()
	if err != nil {
		return err
	}
	//start new one
	var name string
	if filepath.Ext(command.mainFilePath) == ".go" {
		name = filepath.Base(command.mainFilePath[:len(command.mainFilePath)-len(".go")])
	} else {
		name = filepath.Base(command.mainFilePath)
	}
	command.targetFilePath = filepath.Join(command.wd, "bin", name+".exe")

	command.debugPrintln("target file path: ", command.targetFilePath)

	command.cmd = console.NewStdioCmd(command.context, "go", "build", "-o", command.targetFilePath, command.mainFilePath)
	err = console.SetCmdEnv(command.cmd, "GOPATH", command.GOPATH)
	if err != nil {
		return err
	}

	err = command.cmd.Run()
	if err != nil {
		return errors.Sprintf("rebuild error: %s", err.Error())
	}
	_, err = os.Stat(command.targetFilePath)
	if err != nil {
		return err
	}
	command.cmd = console.NewStdioCmd(command.context, command.targetFilePath)
	err = console.SetCmdEnv(command.cmd, "GOPATH", command.GOPATH)
	if err != nil {
		return err
	}
	err = command.cmd.Start()
	if err != nil {
		return errors.Sprintf("restart error: %s", err.Error())
	}
	fmt.Println("[kmg] app running pid:", command.cmd.Process.Pid)
	go func() {
		cmd := command.cmd
		err := cmd.Wait()
		fmt.Println("[kmg] app exit pid:", cmd.Process.Pid)
		if err != nil {
			command.debugPrintln("wait error: ", err.Error())
		}
	}()
	return nil
}

func (command *GoWatch) stop() error {
	if command.cmd == nil {
		return nil
	}
	if command.cmd.Process == nil {
		return nil
	}
	err := command.cmd.Process.Kill()
	if err != nil {
		command.debugPrintln("Process.Kill() error", err)
	}
	return nil
}
func (command *GoWatch) debugPrintln(o ...interface{}) {
	if command.isDebug {
		fmt.Println(o...)
	}
}
