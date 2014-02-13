package command

import (
	"flag"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgFile"
)

type Json2Yaml struct {
	inputPath  *string
	outputPath *string
}

func (command *Json2Yaml) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "Json2Yaml",
		Short: "convert from json to yaml",
	}
}

func (command *Json2Yaml) ConfigFlagSet(flag *flag.FlagSet) {
	command.inputPath = flag.String("i", "", "input file path")
	command.outputPath = flag.String("o", "", "output file path")
}
func (command *Json2Yaml) Execute(context *console.Context) error {
	inputPath := *command.inputPath
	outputPath := *command.outputPath
	if inputPath == "" || outputPath == "" {
		return kmgYaml.Json2YamlIo(context.Stdin, context.Stdout)
	}
	transform := &kmgFile.DirectoryFileTransform{
		InputExt:  "yml",
		OuputExt:  "json",
		Transform: kmgYaml.Json2YamlIo,
	}
	return transform.Run(inputPath, outputPath)
}
