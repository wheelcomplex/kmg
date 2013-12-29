package command

import (
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/console/kmgContext"
	"github.com/bronze1man/kmg/kmgFile"
	"path/filepath"
)

type GoCrossComplie struct {
}

func (command *GoCrossComplie) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{
		Name:  "GoCrossComplie",
		Short: "cross compile target in current project",
		Detail: `GoCrossComplie [gofile]
the output file will put into $project_root/bin/name_GOOS_GOARCH[.exe]
`,
	}
}
func (command *GoCrossComplie) Execute(context *console.Context) (err error) {
	if len(context.Args) <= 2 {
		return fmt.Errorf("need gofile parameter")
	}
	targetFile, err := filepath.Abs(context.Args[2])
	if err != nil {
		return
	}
	kmgc, err := kmgContext.FindFromWd()
	if err != nil {
		return
	}
	targetName := kmgFile.GetFileBaseWithoutExt(targetFile)
	for _, target := range kmgc.CrossCompileTarget {
		fileName := targetName + "_" + target.GetGOOS() + "_" + target.GetGOARCH()
		if target.GetGOOS() == "windows" {
			fileName = fileName + ".exe"
		}
		outputFilePath := filepath.Join(kmgc.ProjectPath, "bin", fileName)
		cmd := console.NewStdioCmd(context, "go", "build", "-o", outputFilePath, targetFile)
		console.SetCmdEnv(cmd, "GOOS", target.GetGOOS())
		console.SetCmdEnv(cmd, "GOARCH", target.GetGOARCH())
		console.SetCmdEnv(cmd, "GOPATH", kmgc.GOPATHToString())
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return
}
