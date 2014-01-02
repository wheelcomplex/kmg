package kmgContext

import (
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"os"
	"path/filepath"
	"strings"
)

//if you init it like &Context{xxx},please call Init()
type Context struct {
	GOPATH             []string
	CrossCompileTarget []CompileTarget

	//default to $ProjectPath/app
	AppPath string
	//default to $ProjectPath/config
	ConfigPath string
	//default to $AppPath/data
	DataPath string
	//default to $AppPath/tmp
	TmpPath string

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
func (context *Context) Init() {
	for i, p := range context.GOPATH {
		if filepath.IsAbs(p) {
			continue
		}
		context.GOPATH[i] = filepath.Join(context.ProjectPath, p)
	}
	if context.GOROOT == "" {
		context.GOROOT = os.Getenv("GOROOT")
	}
	if context.AppPath == "" {
		context.AppPath = filepath.Join(context.ProjectPath, "app")
	}
	if context.DataPath == "" {
		context.DataPath = filepath.Join(context.AppPath, "data")
	}
	if context.TmpPath == "" {
		context.TmpPath = filepath.Join(context.AppPath, "tmp")
	}
	if context.ConfigPath == "" {
		context.ConfigPath = filepath.Join(context.AppPath, "config")
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
	context.ProjectPath, err = filepath.Abs(filepath.Dir(kmgFilePath))
	if err != nil {
		return
	}
	context.Init()
	return
}

func FindFromWd() (context *Context, err error) {
	p, err := os.Getwd()
	if err != nil {
		return
	}
	return FindFromPath(p)
}
