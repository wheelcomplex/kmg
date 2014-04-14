package kmgJson

import (
	"encoding/json"
	"github.com/bronze1man/kmg/typeTransform"
	"io/ioutil"
	"os"
)

func ReadFile(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

//读取json文件,并修正json的类型问题(map key 必须是string的问题)
func ReadFileTypeFix(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var iobj interface{}
	err = json.Unmarshal(b, &iobj)
	if err != nil {
		return err
	}
	return typeTransform.Transform(iobj, obj)
}

func WriteFile(path string, obj interface{}) (err error) {
	out, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, out, os.FileMode(0777))
}

//写入json文件,并修正json的类型问题(map key 必须是string的问题)
func WriteFileTypeFix(path string, obj interface{}) (err error) {
	//a simple work around
	obj, err = TypeFixWhenMarshal(obj)
	if err != nil {
		return
	}
	outByte, err := json.Marshal(obj)
	if err != nil {
		return
	}
	return ioutil.WriteFile(path, outByte, os.FileMode(0777))
}

func UnmarshalNoType(r []byte) (interface{}, error) {
	var obj interface{}
	err := json.Unmarshal(r, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// for debug to inspect content in obj
func MustMarshalToString(obj interface{}) string {
	output, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(output)
}
