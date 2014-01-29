package command

import (
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/fsnotify"
	"os"
	"os/exec"
)

type WatchCmd struct {
	context *console.Context
	wd      string
	watcher *fsnotify.Watcher
	cmd     *exec.Cmd
	isDebug bool //more output
}

func (command *WatchCmd) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "WatchCmd",
		Short: "watch current directory and kill and reexecute a command",
	}
}

func (command *WatchCmd) Execute(context *console.Context) error {
	command.isDebug = false
	command.context = context
	if len(context.Args) <= 2 {
		return fmt.Errorf("usage: %s watch [command]", context.ExecutionName)
	}
	var err error
	command.wd, err = os.Getwd()
	if err != nil {
		return err
	}
	runner, err := fsnotify.NewRunner(10000)
	if err != nil {
		return err
	}
	runner.Watcher.ErrorHandler = func(err error) {
		fmt.Println("watcher.Error error: ", err)
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
func (command *WatchCmd) restart() error {
	//kill old process
	err := command.stop()
	if err != nil {
		return err
	}

	command.cmd = console.NewStdioCmd(command.context, command.context.Args[2], command.context.Args[3:]...)
	err = command.cmd.Start()
	if err != nil {
		return fmt.Errorf("restart error: %s", err.Error())
	}
	fmt.Println("[kmg] cmd running pid:", command.cmd.Process.Pid)
	go func() {
		cmd := command.cmd
		err := cmd.Wait()
		fmt.Println("[kmg] cmd exit pid:", cmd.Process.Pid)
		if err != nil {
			command.debugPrintln("wait error: ", err.Error())
		}
	}()
	return nil
}

func (command *WatchCmd) stop() error {
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
func (command *WatchCmd) debugPrintln(o ...interface{}) {
	if command.isDebug {
		fmt.Println(o...)
	}
}
