package buildCommand

import (
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/fsnotify"
	"os"
	"os/exec"
	"time"
)

type WatchCmdCommand struct {
	context *console.Context
	wd      string
	watcher *fsnotify.Watcher
	cmd     *exec.Cmd
	isDebug bool //more output
}

func (command *WatchCmdCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "watchcmd",
		Short: "watch current directory and kill and reexecute a command",
	}
}

func (command *WatchCmdCommand) Execute(context *console.Context) error {
	command.isDebug = false
	command.context = context
	if len(context.Args) <= 2 {
		return errors.Sprintf("usage: %s watch [command]", context.ExecutionName)
	}
	var err error
	command.wd, err = os.Getwd()
	if err != nil {
		return err
	}
	command.watcher, err = fsnotify.NewWatcher(10000)
	if err != nil {
		return err
	}
	command.watcher.ErrorHandler = func(err error) {
		fmt.Println("watcher.Error error: ", err)
	}
	command.watcher.WatchRecursion(command.wd)
	lastHappendTime := time.Now()
	//start app when command start
	command.restart()
	for {
		event := <-command.watcher.Event
		command.debugPrintln("event: ", event)
		if event.Time.Before(lastHappendTime) {
			continue
		}
		//wait 200ms to prevent multiple restart in short time
		time.Sleep(time.Duration(0.2 * float64(time.Second)))
		lastHappendTime = time.Now()
		err := command.restart()
		if err != nil {
			fmt.Println("command.restart() error: ", err)
		}
	}
	// wait forever
	return nil
}
func (command *WatchCmdCommand) restart() error {
	//kill old process
	err := command.stop()
	if err != nil {
		return err
	}

	command.cmd = console.NewStdioCmd(command.context, command.context.Args[2], command.context.Args[3:]...)
	err = command.cmd.Start()
	if err != nil {
		return errors.Sprintf("restart error: %s", err.Error())
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

func (command *WatchCmdCommand) stop() error {
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
func (command *WatchCmdCommand) debugPrintln(o ...interface{}) {
	if command.isDebug {
		fmt.Println(o...)
	}
}
