package command

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/console/kmgContext"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"go/build"
	"os"
	"path/filepath"
)

/*
递归目录的 go test
 支持.kmg.yml目录结构提示文件(该文件必须存在)
 -v 更详细的描述
 -m 一个模块名,从这个模块名开始递归目录测试
 -d 一个目录名,从这个目录开始递归目录测试
*/
type GoTest struct {
	gopath     string
	context    *console.Context
	v          bool
	dir        string
	moduleName string
}

func (command *GoTest) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoTest",
		Short: `递归目录的go test`,
		Detail: `递归目录的go test
 支持.kmg.yml目录结构提示文件(该文件必须存在)
 -v 更详细的描述
 -m 一个模块名,从这个模块名开始递归目录测试
 -d 一个目录名,从这个目录开始递归目录测试`,
	}
}
func (commamd *GoTest) ConfigFlagSet(f *flag.FlagSet) {
	f.BoolVar(&commamd.v, "v", false, "show output of test")
	f.StringVar(&commamd.dir, "d", "", "dir path to test")
	f.StringVar(&commamd.moduleName, "m", "", "module name to test")
}
func (command *GoTest) Execute(context *console.Context) (err error) {
	command.context = context
	kmgc, err := kmgContext.FindFromWd()
	if err != nil {
		return
	}
	command.gopath = kmgc.GOPATH[0]
	//TODO handle several GOPATH
	root := ""
	if context.FlagSet().NArg() == 1 {
		command.moduleName = context.FlagSet().Arg(0)
	}
	if command.dir != "" {
		root = command.dir
		exist, err := kmgFile.FileExist(root)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("[GoTest] dir path:[%s] not found", root)
		}
	}
	if command.moduleName != "" {
		root = filepath.Join(command.gopath, "src", command.moduleName)
		exist, err := kmgFile.FileExist(root)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("[GoTest] module name:[%s] not found", command.moduleName)
		}
	}
	if root == "" {
		root, err = os.Getwd()
		if err != nil {
			return
		}
	}
	c := &build.Context{
		GOPATH:   kmgc.GOPATHToString(),
		Compiler: build.Default.Compiler,
	}
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if kmgFile.IsDotFile(path) {
			return filepath.SkipDir
		}
		pkg, err := c.ImportDir(path, build.ImportMode(0))
		if err != nil {
			return nil
		}
		if len(pkg.TestGoFiles) == 0 {
			return nil
		}
		err = command.gotest(path)
		if err != nil {
			return err
		}
		return nil
	})
}
func (command *GoTest) gotest(path string) error {
	args := []string{"test"}
	if command.v {
		args = append(args, "-v")
	}
	cmd := console.NewStdioCmd(command.context, "go", args...)
	cmd.Dir = path
	err := kmgCmd.SetCmdEnv(cmd, "GOPATH", command.gopath)
	if err != nil {
		return err
	}
	return cmd.Run()
}
