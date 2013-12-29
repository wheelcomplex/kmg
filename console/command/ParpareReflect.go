package command

import (
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgReflect"
	"os"
	"path/filepath"
)

type ParpareReflect struct {
}

func (command *ParpareReflect) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "ParpareReflect", Short: `parse source code to parpare some data for kmgReflect`}
}
func (command *ParpareReflect) Execute(context *console.Context) error {
	//parse all file in GOPATH
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	root := filepath.Join(wd, "src")
	contextDecl, err := kmgReflect.NewContextDeclFromSrcPath(root)
	if err != nil {
		return err
	}
	out, err := kmgYaml.Marshal(contextDecl)
	if err != nil {
		return err
	}

	_, err = context.Stdout.Write(out)
	if err != nil {
		return err
	}
	return nil
}
