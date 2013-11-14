package command

import (
	"flag"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/encoding/excel"
	"github.com/bronze1man/kmg/errors"
	"launchpad.net/goyaml"
)

type Xlsx2Yaml struct {
	filePath         *string
	format           *string
	isOutputAllSheet *bool
}

func (command *Xlsx2Yaml) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "Xlsx2Yaml",
		Short: "convert from xlsx(Microsoft excel 2007) to yaml",
	}
}
func (command *Xlsx2Yaml) ConfigFlagSet(f *flag.FlagSet) {
	command.filePath = f.String("input", "", "input file path")
	command.format = f.String("format", "grid", "output yaml format(grid,raw)")
	command.isOutputAllSheet = f.Bool("outputAllSheet", false, "is output all sheet(default just out first one)?")
}
func (command *Xlsx2Yaml) Execute(context *console.Context) error {
	if *command.filePath == "" {
		return errors.New("need input file")
	}
	rawArray, err := excel.XlsxFile2Array(*command.filePath)
	if err != nil {
		return err
	}
	var outByte []byte
	switch *command.format {
	case "raw":
		if *command.isOutputAllSheet {
			outByte, err = goyaml.Marshal(rawArray)
			if err != nil {
				return err
			}
		} else {
			outByte, err = goyaml.Marshal(rawArray[0])
			if err != nil {
				return err
			}
		}
	case "grid":
		o := [][]map[string]string{}
		for _, s := range rawArray {
			o1, err := excel.TitleArrayToGrid(s)
			if err != nil {
				return err
			}
			o = append(o, o1)
		}
		if *command.isOutputAllSheet {
			outByte, err = goyaml.Marshal(o)
			if err != nil {
				return err
			}
		} else {
			outByte, err = goyaml.Marshal(o[0])
			if err != nil {
				return err
			}
		}
	default:
		return errors.Sprintf("not support output format: %s", command.format)
	}
	_, err = context.Stdout.Write(outByte)
	if err != nil {
		return err
	}
	return nil
}
