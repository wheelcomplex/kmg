package kmgJson

import (
	"github.com/bronze1man/kmg/encoding/kmgYaml"
)

//修正序列化时的类型问题,此处直接把类型搞定,(此处会忽略对系统的json有效的tag)
//TODO finish it
func TypeFixWhenMarshal(obj interface{}) (out interface{}, err error) {
	//a simple work around
	yamlBytes, err := kmgYaml.Marshal(obj)
	if err != nil {
		return
	}
	var yamlData interface{}
	err = kmgYaml.Unmarshal(yamlBytes, &yamlData)
	if err != nil {
		return
	}
	out, err = kmgYaml.Yaml2JsonTransformData(yamlData)
	if err != nil {
		return
	}
	return
}

/*
func marshalTypeFix(in reflect.Value) (out reflect.Value, err error) {
	switch in.Kind() {
	case reflect.Map:

	case reflect.Struct:
	case reflect.Ptr:
		if !in.IsNil(){
			return marshalTypeFix(in.Elem())
		}
	case reflect.Slice:
	case reflect.Interface:
	}
	//其他类型不处理
	return in,err
	return nil,fmt.Errorf("[kmgJson.marshalTypeFix] not support in kind: %s", in.Kind())
}
*/
