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
 TODO 处理多个GOPATH的问题,从GOPATH里面找到这个模块
 支持.kmg.yml目录结构提示文件(该文件必须存在)
 -v 更详细的描述
 -m 一个模块名,从这个模块名开始递归目录测试
 -d 一个目录名,从这个目录开始递归目录测试
*/
type GoTest struct {
	gopath       string
	context      *console.Context
	v            bool
	dir          string
	moduleName   string
	bench        string
	onePackage   bool
	buildContext *build.Context
}

func (command *GoTest) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoTest",
		Short: `递归目录的go test`,
		Detail: `递归目录的go test
 支持.kmg.yml目录结构提示文件(该文件必须存在)
 -v 更详细的描述
 -m 一个模块名,从这个模块名开始递归目录测试
 -d 一个目录名,从这个目录开始递归目录测试
 -bench benchmarks参数,直接传递到go test
 -onePackage 不递归目录测试,仅测试一个package`,
	}
}
func (commamd *GoTest) ConfigFlagSet(f *flag.FlagSet) {
	f.BoolVar(&commamd.v, "v", false, "show output of test")
	f.StringVar(&commamd.dir, "d", "", "dir path to test")
	f.StringVar(&commamd.moduleName, "m", "", "module name to test")
	f.StringVar(&commamd.bench, "bench", "", "bench parameter pass to go test")
	f.BoolVar(&commamd.onePackage, "onePackage", false, "only test one package")
}
func (command *GoTest) Execute(context *console.Context) (err error) {
	command.context = context
	kmgc, err := kmgContext.FindFromWd()
	if err == nil {
		command.gopath = kmgc.GOPATH[0]
	} else {
		if kmgContext.IsNotFound(err) {
			command.gopath = os.Getenv("GOPATH")
		} else {
			return
		}
	}
	//find root path
	root, err := command.findRootPath(context)
	if err != nil {
		return
	}
	command.buildContext = &build.Context{
		GOPATH:   command.gopath,
		Compiler: build.Default.Compiler,
	}
	if command.onePackage {
		return command.handlePath(root)
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
		return command.handlePath(path)
	})
}

func (command *GoTest) findRootPath(context *console.Context) (root string, err error) {
	if context.FlagSet().NArg() == 1 {
		command.moduleName = context.FlagSet().Arg(0)
	}
	if command.dir != "" {
		root = command.dir
		exist, err := kmgFile.FileExist(root)
		if err != nil {
			return "", err
		}
		if !exist {
			return "", fmt.Errorf("[GoTest] dir path:[%s] not found", root)
		}
		return root, nil
	}
	if command.moduleName != "" {
		//TODO 处理多个GOPATH的问题,从GOPATH里面找到这个模块
		root = filepath.Join(command.gopath, "src", command.moduleName)
		exist, err := kmgFile.FileExist(root)
		if err != nil {
			return "", err
		}
		if !exist {
			return "", fmt.Errorf("[GoTest] module name:[%s] not found", command.moduleName)
		}
		return root, nil
	}
	if root == "" {
		root, err = os.Getwd()
		if err != nil {
			return
		}
	}
	return
}

func (command *GoTest) handlePath(path string) error {
	pkg, err := command.buildContext.ImportDir(path, build.ImportMode(0))
	if err != nil {
		//仅忽略 不是golang的目录的错误
		_, ok := err.(*build.NoGoError)
		if ok {
			return nil
		}
		return err
	}
	if pkg.IsCommand() {
		return nil
	}
	if len(pkg.TestGoFiles) == 0 {
		//如果没有测试文件,还会尝试build一下这个目录
		return command.gobuild(path)
		//return nil
	}
	return command.gotest(path)
}

func (command *GoTest) gotest(path string) error {
	fmt.Printf("[gotest] path[%s]\n", path)
	args := []string{"test"}
	if command.v {
		args = append(args, "-v")
	}
	if command.bench != "" {
		args = append(args, "-bench", command.bench)
	}
	cmd := console.NewStdioCmd(command.context, "go", args...)
	cmd.Dir = path
	err := kmgCmd.SetCmdEnv(cmd, "GOPATH", command.gopath)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func (command *GoTest) gobuild(path string) error {
	/*
		args := []string{"build"}
		if command.v {
			fmt.Println("building ")
		} */
	fmt.Printf("[gobuild] path[%s]\n", path)
	cmd := console.NewStdioCmd(command.context, "go", "build")
	cmd.Dir = path
	err := kmgCmd.SetCmdEnv(cmd, "GOPATH", command.gopath)
	if err != nil {
		return err
	}
	return cmd.Run()
}
