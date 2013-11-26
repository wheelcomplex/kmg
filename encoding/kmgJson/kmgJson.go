package kmgJson

import (
	"encoding/json"
	"io/ioutil"
)

func ReadFile(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

func UnmarshalNoType(r []byte) (interface{}, error) {
	var obj interface{}
	err := json.Unmarshal(r, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
