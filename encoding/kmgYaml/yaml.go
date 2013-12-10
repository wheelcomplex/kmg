package kmgYaml

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"strconv"
)

//快速解决 goyaml 的各种奇葩问题...->转成json->转成对象
func ReadFile(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	jsonb, err := Yaml2JsonBytes(b)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonb, obj)
}

func ReadFileGoyaml(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return goyaml.Unmarshal(b, obj)
}

func WriteFileGoyaml(path string, obj interface{}) error {
	out, err := goyaml.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, out, os.FileMode(0777))
}

func Yaml2JsonBytes(input []byte) (output []byte, err error) {
	var data interface{}
	err = goyaml.Unmarshal(input, &data)
	if err != nil {
		return nil, err
	}
	data, err = Yaml2JsonTransformData(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}
func Yaml2JsonTransformData(in interface{}) (out interface{}, err error) {
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
			v, err = Yaml2JsonTransformData(v)
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
			o[i], err = Yaml2JsonTransformData(in1[i])
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
