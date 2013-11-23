package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/kmgFile"
	"io"
	"io/ioutil"
	"launchpad.net/goyaml"
	"strconv"
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
	output, err := yaml2JsonBytes(input)
	if err != nil {
		return err
	}
	_, err = w.Write(output)
	return err
}
func yaml2JsonBytes(input []byte) (output []byte, err error) {
	var data interface{}
	err = goyaml.Unmarshal(input, &data)
	if err != nil {
		return nil, err
	}
	data, err = yaml2JsonTransformData(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}
func yaml2JsonTransformData(in interface{}) (out interface{}, err error) {
	switch in.(type) {
	case map[interface{}]interface{}:
		o := make(map[string]interface{})
		for k, v := range in.(map[interface{}]interface{}) {
			sk := ""
			switch k.(type) {
			case string:
				sk = k.(string)
			case int:
				sk = strconv.Itoa(k.(int))
			default:
				return nil, fmt.Errorf("type not match: expect map key string or int get: %T", k)
			}
			v, err = yaml2JsonTransformData(v)
			if err != nil {
				return nil, err
			}
			o[sk] = v
		}
		return o, nil
	case []interface{}:
		in1 := in.([]interface{})
		len1 := len(in1)
		o := make([]interface{}, len1)
		for i := 0; i < len1; i++ {
			o[i], err = yaml2JsonTransformData(in1[i])
			if err != nil {
				return nil, err
			}
		}
		return o, nil
	default:
		return in, nil
	}
	return in, nil
}
