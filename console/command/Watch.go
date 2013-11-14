package command

import (
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/fsnotify"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//Warning!! do not use go run to build executes...
// you can kill go run xxx, you can not kill main.exe xxx
//TODO cygwin ctrl+c not exit children processes
//TODO wrap restart cmd stuff
//TODO wrap lastHappendTime watch stuff
type Watch struct {
	context        *console.Context
	wd             string
	watcher        *fsnotify.Watcher
	cmd            *exec.Cmd
	mainFilePath   string
	targetFilePath string
	isDebug        bool //more output
}

func (command *Watch) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "Watch",
		Short: "watch current directory and rebuild and restart app when some file changed",
	}
}

func (command *Watch) Execute(context *console.Context) error {
	command.isDebug = false
	command.context = context
	if len(context.Args) != 3 {
		return errors.Sprintf("usage: %s watch [packages]", context.ExecutionName)
	}
	command.mainFilePath = context.Args[2]
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
	command.watcher.IsIgnorePath = func(path string) bool {
		if fsnotify.DefaultIsIgnorePath(path) {
			return true
		}
		if path == command.targetFilePath {
			return true
		}
		return false
	}
	command.watcher.WatchRecursion(command.wd)
	lastHappendTime := time.Now()
	//start app when command start
	err = command.restart()
	if err != nil {
		fmt.Println("command.restart() error: ", err)
	}
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
func (command *Watch) restart() error {
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
	err = console.SetCmdEnv(command.cmd, "GOPATH", command.wd)
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
	err = console.SetCmdEnv(command.cmd, "GOPATH", command.wd)
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

func (command *Watch) stop() error {
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
func (command *Watch) debugPrintln(o ...interface{}) {
	if command.isDebug {
		fmt.Println(o...)
	}
}

func (command *Watch) setEnv() error {
	env, err := console.NewEnvFromArray(os.Environ())
	if err != nil {
		fmt.Printf("%#v", os.Environ())
		return err
	}
	env.Values["GOPATH"] = command.wd
	command.cmd.Env = env.ToArray()
	return nil
}
