package kmgYaml

import (
	"io/ioutil"
	"launchpad.net/goyaml"
)

func ReadYamlFile(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return goyaml.Unmarshal(b, obj)
}
