package command

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/console/kmgContext"
	"github.com/bronze1man/kmg/kmgFile"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type GoAllImport struct {
	gopath      string
	dir         string
	packageName string
}

func (command *GoAllImport) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoAllImport",
		Short: `generate a golang file import all package in a dir`,
		Detail: `generate a golang file import all package in a dir
 需要.kmg.yml目录结构提示文件(该文件必须存在)
 -d 传入的那个目录名`,
	}
}
func (commamd *GoAllImport) ConfigFlagSet(f *flag.FlagSet) {
	f.StringVar(&commamd.dir, "d", "", "dir path")
	f.StringVar(&commamd.packageName, "p", "main", "package name in generate golang file(default to main)")
}
func (command *GoAllImport) Execute(context *console.Context) (err error) {
	kmgc, err := kmgContext.FindFromWd()
	if err != nil {
		return
	}
	command.gopath = kmgc.GOPATH[0]
	root := command.dir
	exist, err := kmgFile.FileExist(root)
	if err != nil {
		return
	}
	if !exist {
		return fmt.Errorf("[GoAllImport] dir path[%s] not exist", root)
	}
	c := &build.Context{
		GOPATH:   kmgc.GOPATHToString(),
		Compiler: build.Default.Compiler,
	}
	root, err = kmgFile.Realpath(root)
	if err != nil {
		return
	}
	moduleList := []string{}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
			//忽略异常文件夹,(不是golang的目录之类的)
			return nil
		}
		if pkg.IsCommand() {
			return nil
		}

		pkgFullName, err := filepath.Rel(command.gopath, path)
		if err != nil {
			return err
		}
		moduleList = append(moduleList, strings.TrimPrefix(filepath.ToSlash(pkgFullName), "src/"))
		return nil
	})
	if err != nil {
		return
	}
	tpl := template.Must(template.New("").Parse(`package {{.PackageName}}
import(
{{range .ModuleList}}	_ "{{.}}"
{{end}}
)
`))
	buf := &bytes.Buffer{}
	tpl.Execute(buf, struct {
		PackageName string
		ModuleList  []string
	}{
		PackageName: command.packageName,
		ModuleList:  moduleList,
	})
	fmt.Println(string(buf.Bytes()))
	return
}
