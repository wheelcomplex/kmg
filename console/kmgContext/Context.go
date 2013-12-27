package kmgContext

import (
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"os"
	"path/filepath"
	"strings"
)

type Context struct {
	GOPATH      []string
	ProjectPath string
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
		p = filepath.Dir(p)
		if p == "" {
			err = fmt.Errorf("not found .kmg.yml in the project dir")
			return
		}
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
	return FindFromPath(p)
}
