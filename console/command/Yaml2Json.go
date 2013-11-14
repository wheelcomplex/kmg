package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bronze1man/kmg/console"
	"io/ioutil"
	"launchpad.net/goyaml"
	"strconv"
)

type Yaml2Json struct {
}

func (command *Yaml2Json) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "Yaml2Json",
		Short: "convert from yaml to json",
	}
}

func (command *Yaml2Json) Execute(context *console.Context) error {
	var data interface{}
	input, err := ioutil.ReadAll(context.Stdin)
	if err != nil {
		return err
	}
	err = goyaml.Unmarshal(input, &data)
	if err != nil {
		return err
	}
	data, err = yaml2JsonTransformData(data)
	if err != nil {
		return err
	}

	output, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = context.Stdout.Write([]byte(output))
	return err
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
				return nil, errors.New(
					fmt.Sprintf("type not match: expect map key string or int get: %T", k))
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
