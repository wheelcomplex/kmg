package command

import (
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/console/kmgContext"
	"github.com/bronze1man/kmg/kmgFile"
	"go/build"
	"os"
	"path/filepath"
)

type GoTest struct {
	wd      string
	context *console.Context
}

func (command *GoTest) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoTest", Short: `test all go package in a directory in current project`}
}
func (command *GoTest) Execute(context *console.Context) (err error) {
	command.context = context
	kmgc, err := kmgContext.FindFromWd()
	if err != nil {
		return
	}
	command.wd = kmgc.GOPATH[0]
	//TODO handle several GOPATH
	root := filepath.Join(command.wd, "src")
	if context.FlagSet().NArg() == 1 {
		root = filepath.Join(root, context.FlagSet().Arg(0))
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
	cmd := console.NewStdioCmd(command.context, "go", "test")
	cmd.Dir = path
	err := console.SetCmdEnv(cmd, "GOPATH", command.wd)
	if err != nil {
		return err
	}
	return cmd.Run()
}
