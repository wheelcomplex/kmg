package command

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/console/kmgContext"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"path/filepath"
)

type GoCrossCompile struct {
	outputPath string
}

func (command *GoCrossCompile) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{
		Name:  "GoCrossCompile",
		Short: "cross compile target in current project",
		Detail: `GoCrossComplie [gofile]
the output file will put into $project_root/bin/name_GOOS_GOARCH[.exe]
`,
	}
}
func (command *GoCrossCompile) ConfigFlagSet(flag *flag.FlagSet) {
	flag.StringVar(&command.outputPath, "o", "", "output file dir(file name come from source file name),default to $project_root/bin")
}
func (command *GoCrossCompile) Execute(context *console.Context) (err error) {
	if len(context.Args) <= 2 {
		return fmt.Errorf("need gofile parameter")
	}
	targetFile := context.FlagSet().Arg(0)
	kmgc, err := kmgContext.FindFromWd()
	if err != nil {
		return
	}
	targetName := kmgFile.GetFileBaseWithoutExt(targetFile)
	if command.outputPath == "" {
		command.outputPath = filepath.Join(kmgc.ProjectPath, "bin")
	}
	for _, target := range kmgc.CrossCompileTarget {
		fileName := targetName + "_" + target.GetGOOS() + "_" + target.GetGOARCH()
		if target.GetGOOS() == "windows" {
			fileName = fileName + ".exe"
		}
		outputFilePath := filepath.Join(command.outputPath, fileName)
		cmd := kmgCmd.NewStdioCmd(context, "go", "build", "-o", outputFilePath, targetFile)
		kmgCmd.SetCmdEnv(cmd, "GOOS", target.GetGOOS())
		kmgCmd.SetCmdEnv(cmd, "GOARCH", target.GetGOARCH())
		kmgCmd.SetCmdEnv(cmd, "GOPATH", kmgc.GOPATHToString())
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return
}
