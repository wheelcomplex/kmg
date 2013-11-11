package buildCommand

import (
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/errors"
	"github.com/howeyc/fsnotify"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//Warning!! do not use go run to build executes...
// you can kill go run xxx, you can not kill main.exe xxx
type WatchCommand struct {
	context *console.Context
	wd      string
	watcher *fsnotify.Watcher
	watched map[string]struct{}
	//args []string
	cmd            *exec.Cmd
	mainFilePath   string
	targetFilePath string
	isDebug        bool //more output
}

func (command *WatchCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "watch",
		Short: "watch current directory and rebuild and restart app when some file changed",
	}
}

//TODO trace all watched directory
//TODO wrapper watch directory stuff
func (command *WatchCommand) Execute(context *console.Context) error {
	command.isDebug = false
	command.context = context
	if len(context.Args) != 3 {
		return errors.Sprintf("usage: %s watch [packages]", context.ExecutionName)
	}
	command.mainFilePath = context.Args[2]
	command.watched = make(map[string]struct{})
	var err error
	command.wd, err = os.Getwd()
	if err != nil {
		return err
	}
	command.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	command.addFolder(command.wd)
	changeChan := make(chan time.Time, 10000)
	go func() {
		lastHappendTime := time.Now().Add(-time.Hour)
		//start app when command start
		command.restart()
		for {
			thisTime := <-changeChan
			if thisTime.Before(lastHappendTime) {
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
	}()
	for {
		select {
		case ev := <-command.watcher.Event:
			command.debugPrintln("event: ", ev)
			if command.isIgnorePath(ev.Name) {
				continue
			}
			changeChan <- time.Now()
			// you can not stat a delete file...
			if ev.IsDelete() {
				continue
			}
			fi, err := os.Stat(ev.Name)
			if err != nil {
				//rename send two events,one old file,one new file,here ignore old one
				if os.IsNotExist(err) {
					continue
				}
				fmt.Println("os.Stat error: ", err)
				continue
			}
			if fi.IsDir() {
				if ev.IsCreate() {
					command.addFolder(ev.Name)
				}
			}

		case err := <-command.watcher.Error:
			fmt.Println("watcher.Error error: ", err)
		}
	}
	// wait forever
	return nil
}
func (command *WatchCommand) isIgnorePath(path string) bool {
	base := filepath.Base(path)
	if filepath.HasPrefix(base, ".") {
		return true
	}
	if base == "." {
		return true
	}
	if base == ".." {
		return true
	}
	if path == command.targetFilePath {
		return true
	}
	return false
}
func (command *WatchCommand) addFolder(path string) error {
	folders, err := command.subFolders(path)
	for _, v := range folders {
		err = command.watcher.Watch(v)
		if err != nil {
			return err
		}
	}
	return nil
}
func (command *WatchCommand) subFolders(path string) (paths []string, err error) {
	err = filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}
		// skip folders that begin with a dot
		if command.isIgnorePath(newPath) {
			return filepath.SkipDir
		}
		paths = append(paths, newPath)
		return nil
	})
	return paths, err
}
func (command *WatchCommand) restart() error {
	//kill old process
	err := command.stop()
	if err != nil {
		return err
	}
	//start new one
	fmt.Println("rebuild app...")
	var name string
	if filepath.Ext(command.mainFilePath) == ".go" {
		name = filepath.Base(command.mainFilePath[:len(command.mainFilePath)-len(".go")])
	} else {
		name = filepath.Base(command.mainFilePath)
	}
	command.targetFilePath = filepath.Join(command.wd, "bin", name+".exe")

	command.debugPrintln("target file path: ", command.targetFilePath)

	command.cmd = console.NewStdioCmd(command.context, "go", "build", "-o", command.targetFilePath, command.mainFilePath)
	command.cmd.Env = append(os.Environ(), "GOPATH="+command.wd)

	err = command.cmd.Run()
	if err != nil {
		return errors.Sprintf("rebuild error: %s", err.Error())
	}
	_, err = os.Stat(command.targetFilePath)
	if err != nil {
		return err
	}
	fmt.Println("restart app...")
	command.cmd = console.NewStdioCmd(command.context, command.targetFilePath)
	command.cmd.Env = append(os.Environ(), "GOPATH="+command.wd)
	err = command.cmd.Start()
	if err != nil {
		return errors.Sprintf("restart error: %s", err.Error())
	}
	fmt.Println("app running... pid:", command.cmd.Process.Pid)
	return nil
}

func (command *WatchCommand) stop() error {
	if command.cmd == nil {
		return nil
	}
	if command.cmd.Process == nil {
		return nil
	}
	err := command.cmd.Process.Signal(os.Kill)
	if err != nil {
		command.debugPrintln("Process.Kill() error", err)
	}
	fmt.Println("kill app...")
	err = command.cmd.Wait()
	if err != nil {
		command.debugPrintln("cmd.Wait() error", err)
	}
	return nil
}
func (command *WatchCommand) debugPrintln(o ...interface{}) {
	if command.isDebug {
		fmt.Println(o...)
	}
}
