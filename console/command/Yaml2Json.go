package command

import (
	"flag"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgFile"
	"io"
	"io/ioutil"
)

type Yaml2Json struct {
	inputPath  *string
	outputPath *string
}

func (command *Yaml2Json) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "Yaml2Json",
		Short: "convert from yaml to json",
	}
}

func (command *Yaml2Json) ConfigFlagSet(flag *flag.FlagSet) {
	command.inputPath = flag.String("i", "", "input file path")
	command.outputPath = flag.String("o", "", "output file path")
}
func (command *Yaml2Json) Execute(context *console.Context) error {
	inputPath := *command.inputPath
	outputPath := *command.outputPath
	if inputPath == "" || outputPath == "" {
		return yaml2JsonIo(context.Stdin, context.Stdout)
	}
	transform := &kmgFile.DirectoryFileTransform{
		InputExt:  "yml",
		OuputExt:  "json",
		Transform: yaml2JsonIo,
	}
	return transform.Run(inputPath, outputPath)
}

func yaml2JsonIo(r io.Reader, w io.Writer) error {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	output, err := kmgYaml.Yaml2JsonBytes(input)
	if err != nil {
		return err
	}
	_, err = w.Write(output)
	return err
}
