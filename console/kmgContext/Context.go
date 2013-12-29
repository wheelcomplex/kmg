package kmgContext

import (
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"os"
	"path/filepath"
	"strings"
)

type Context struct {
	GOPATH             []string
	CrossCompileTarget []CompileTarget
	//should come from environment
	GOROOT string
	//should come from dir of ".kmg.yml"
	ProjectPath string
}

// see http://golang.org/doc/install/source to get all possiable GOOS and GOARCH
// should be something like "windows_amd64","darwin_386",etc..
type CompileTarget string

func (target CompileTarget) GetGOOS() string {
	part := strings.Split(string(target), "_")
	return part[0]
}
func (target CompileTarget) GetGOARCH() string {
	part := strings.Split(string(target), "_")
	return part[1]
}

func (context *Context) GOPATHToString() string {
	return strings.Join(context.GOPATH, ":")
}
func (context *Context) init() {
	for i, p := range context.GOPATH {
		if filepath.IsAbs(p) {
			continue
		}
		context.GOPATH[i] = filepath.Join(context.ProjectPath, p)
	}
	if context.GOROOT == "" {
		context.GOROOT = os.Getenv("GOROOT")
	}
}
func FindFromPath(p string) (context *Context, err error) {
	p, err = filepath.Abs(p)
	if err != nil {
		return
	}
	var kmgFilePath string
	for {
		kmgFilePath = filepath.Join(p, ".kmg.yml")
		_, err = os.Stat(kmgFilePath)
		if err == nil {
			//found it
			break
		}
		if !os.IsNotExist(err) {
			return
		}
		thisP := filepath.Dir(p)
		if p == thisP {
			err = fmt.Errorf("not found .kmg.yml in the project dir")
			return
		}
		p = thisP
	}
	context = &Context{}
	err = kmgYaml.ReadFileGoyaml(kmgFilePath, context)
	if err != nil {
		return
	}
	context.ProjectPath = filepath.Dir(kmgFilePath)
	context.init()
	return
}

func FindFromWd() (context *Context, err error) {
	p, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Println(p)
	return FindFromPath(p)
}
